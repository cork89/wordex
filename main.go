package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//go:embed fantasy.txt
var fantasytxt string

//go:embed fantasynames.txt
var fantasynamestxt string

//go:embed scifi.txt
var scifitxt string

//go:embed mystery.txt
var mysterytxt string

type WordList int

const (
	Fantasy WordList = iota
	Scifi
	Mystery
	Fantasynames
)

func (w WordList) String() string {
	switch w {
	case Fantasy:
		return "fantasy"
	case Scifi:
		return "scifi"
	case Mystery:
		return "mystery"
	case Fantasynames:
		return "fantasynames"
	default:
		return "Unknown"
	}
}

func ParseWordList(s string) (WordList, error) {
	switch strings.ToLower(s) {
	case "fantasy":
		return Fantasy, nil
	case "scifi":
		return Scifi, nil
	case "mystery":
		return Mystery, nil
	case "fantasynames":
		return Fantasynames, nil
	default:
		return Fantasy, errors.New("invalid word list: " + s)
	}
}

var wordGroupByType = map[string]*WordGroup{
	"fantasy":      nil,
	"scifi":        nil,
	"mystery":      nil,
	"fantasynames": nil,
}

var colors = []string{
	"#FFB6C1", // LightPink
	"#ADD8E6", // LightBlue
	"#90EE90", // LightGreen
	"#FFD700", // Gold
	"#FFA07A", // LightSalmon
	"#E6E6FA", // Lavender
	"#F0E68C", // Khaki
	"#D8BFD8", // Thistle
}

var colorsLen = len(colors)

type Word struct {
	Word    string `json:"word"`
	Color   string `json:"color"`
	Subtext string `json:"subtext"`
}

type Words struct {
	Words []Word `json:"words"`
}

func (w Words) getWordsString() (wordString string) {
	words := make([]string, 0, len(w.Words))
	for _, word := range w.Words {
		words = append(words, word.Word)
	}

	return strings.Join(words, ",")
}

type WordGroup struct {
	words    []Word
	index    int
	wordType WordList
}

func (w *WordGroup) shuffleWords() {
	rand.Shuffle(len(w.words), func(i, j int) {
		w.words[i], w.words[j] = w.words[j], w.words[i]
	})
	for i := range w.words {
		w.words[i].Color = colors[i%colorsLen]
	}
}

func (w *WordGroup) getRandomWords(n int) Words {
	words := make([]Word, 0, n)
	if w.index+n >= len(w.words) {
		w.index = 0
		w.shuffleWords()
	}

	for i := 0; i < n; i += 1 {
		words = append(words, w.words[w.index+i])
		//Word{Word: w.words[w.index+i], Color: colors[colorIdx]})
		colorIdx = (colorIdx + 1) % len(colors)
	}
	w.index += n

	return Words{Words: words}
}

func (w *WordGroup) simpleString() string {
	return fmt.Sprintf("type: %v, index: %d", w.wordType, w.index)
}

func (w *WordGroup) getSavedWords(savedWords string) Words {
	splitwords := strings.Split(savedWords, ",")

	words := make([]Word, 0, len(splitwords))
	for i, v := range splitwords {
		word := html.UnescapeString(v)
		var subtext = ""
		for _, group := range w.words {
			if group.Word == word {
				subtext = group.Subtext
			}
		}
		words = append(words, Word{Word: word, Color: colors[colorIdx], Subtext: subtext})
		colorIdx = (colorIdx + 1) % len(colors)
		if i > 2 {
			break
		}
	}

	return Words{Words: words}
}

var colorIdx = 0

type Colors struct {
	Colors []string `json:"Colors"`
}

// func typeHandler(r *http.Request) *WordGroup {
// 	wordsType := r.URL.Query().Get("type")

// 	wordGroup := wordGroupByType[wordsType]

// 	if wordGroup == nil {
// 		wordGroup = wordGroupByType["fantasy"]
// 	}

// 	return wordGroup
// }

func typeHandler(r *http.Request) (WordList, error) {
	wordsType := r.URL.Query().Get("type")

	if wordsType == "" {
		return Fantasy, nil
	}
	return ParseWordList(wordsType)
}

func handler(w http.ResponseWriter, r *http.Request) {
	savedWords := r.URL.Query().Get("words")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var words Words

	wordCategory, err := typeHandler(r)

	if err != nil {
		log.Println("defaulting to fantasy")
	}

	if savedWords != "" {
		words, err = getSavedWords(savedWords, wordCategory)
	} else {
		words, err = getRandomWords(2+rand.Intn(2), wordCategory)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := Index(words)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template render error: %v", err)
	}
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	wordCategory, err := typeHandler(r)

	if err != nil {
		log.Println("defaulting to fantasy")
	}

	words, err := getRandomWords(2+rand.Intn(2), wordCategory)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("words", words.getWordsString())

	component := WordsDiv(words)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Failed to render words", http.StatusInternalServerError)
	}
}

var validWordIdx = regexp.MustCompile("^[0-2]{1}$")

func extractGameId(r *http.Request) (string, error) {
	m := validWordIdx.FindStringSubmatch(r.PathValue("word"))
	if len(m) != 1 {
		return "", errors.New("invalid word idx")
	}
	return m[0], nil
}

func wordHandler(w http.ResponseWriter, r *http.Request) {
	wordIdx, err := extractGameId(r)
	if err != nil {
		http.Error(w, "Failed to retrieve word index", http.StatusInternalServerError)
		return
	}
	savedWords := r.URL.Query().Get("words")
	savedWordsSplit := strings.Split(savedWords, ",")
	idx, err := strconv.Atoi(wordIdx)
	if err != nil {
		http.Error(w, "Failed to convert word index", http.StatusInternalServerError)
		return
	}

	wordCategory, err := typeHandler(r)

	if err != nil {
		log.Println("defaulting to fantasy")
	}

	words, err := getRandomWords(1, wordCategory)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if savedWordsSplit[0] == "" {
		savedWordsSplit = []string{"", "", ""}
	}
	savedWordsSplit[idx] = words.Words[0].Word

	w.Header().Set("words", strings.Join(savedWordsSplit, ","))

	component := WordDiv(words.Words[0], wordIdx)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Failed to render word", http.StatusInternalServerError)
	}
}

var re = regexp.MustCompile((`,\r?\n`))

func createWordGroup(wordsraw [2]string) WordGroup {
	var wordGroup WordGroup
	words := re.Split(wordsraw[1], -1)
	wordGroup.words = make([]Word, 0, len(words))
	for i, v := range words {
		wordparts := strings.Split(v, " : ")

		word := Word{Word: wordparts[0], Color: colors[i%colorsLen]}
		if len(wordparts) > 1 {
			word.Subtext = wordparts[1]
		}
		wordGroup.words = append(wordGroup.words, word)
	}
	return wordGroup
}

func main() {
	err := initDataaccess()

	if err != nil {
		log.Panicf("failed to init dataaccess, err: %v\n", err)
	}

	var initialWords = [][2]string{
		{"fantasy", fantasytxt},
		{"scifi", scifitxt},
		{"mystery", mysterytxt},
		{"fantasynames", fantasynamestxt},
	}

	for i := range initialWords {
		tempWord := createWordGroup(initialWords[i])
		tempWord.shuffleWords()
		wordGroupByType[initialWords[i][0]] = &tempWord
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/words/", wordsHandler)
	http.HandleFunc("/word/{word}/", wordHandler)

	log.Println("Starting server at http://localhost:8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

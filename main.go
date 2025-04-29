package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//go:embed fantasy.txt
var fantasytxt string

//go:embed scifi.txt
var scifitxt string

//go:embed mystery.txt
var mysterytxt string

var fantasywords WordGroup
var scifiwords WordGroup
var mysterywords WordGroup

type WordList int

const (
	Fantasy WordList = iota
	Scifi
	Mystery
)

var wordGroupByType = map[string]*WordGroup{
	"fantasy": nil,
	"scifi":   nil,
	"mystery": nil,
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

type Word struct {
	Word  string `json:"word"`
	Color string `json:"color"`
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
	words    []string
	index    int
	wordType WordList
}

func (w *WordGroup) shuffleWords() {
	rand.Shuffle(len(w.words), func(i, j int) {
		w.words[i], w.words[j] = w.words[j], w.words[i]
	})
}

func (w *WordGroup) getRandomWords(n int) Words {
	words := make([]Word, 0, n)
	if w.index+n >= len(w.words) {
		w.index = 0
		w.shuffleWords()
	}

	for i := 0; i < n; i += 1 {
		words = append(words, Word{Word: w.words[w.index+i], Color: colors[colorIdx]})
		colorIdx = (colorIdx + 1) % len(colors)
	}
	w.index += n

	return Words{Words: words}
}

func (w *WordGroup) simpleString() string {
	return fmt.Sprintf("type: %v, index: %d", w.wordType, w.index)
}

var colorIdx = 0

func getSavedWords(savedWords string) Words {
	splitwords := strings.Split(savedWords, ",")

	words := make([]Word, 0, len(splitwords))
	for i, word := range splitwords {
		words = append(words, Word{Word: word, Color: colors[colorIdx]})
		colorIdx = (colorIdx + 1) % len(colors)
		if i > 2 {
			break
		}
	}

	return Words{Words: words}
}

type Colors struct {
	Colors []string `json:"Colors"`
}

func typeHandler(r *http.Request) *WordGroup {
	wordsType := r.URL.Query().Get("type")

	wordGroup := wordGroupByType[wordsType]

	if wordGroup == nil {
		wordGroup = wordGroupByType["fantasy"]
	}

	return wordGroup
}

func handler(w http.ResponseWriter, r *http.Request) {
	savedWords := r.URL.Query().Get("words")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var words Words

	wordGroup := typeHandler(r)

	if savedWords != "" {
		words = getSavedWords(savedWords)

	} else {
		words = wordGroup.getRandomWords(2 + rand.Intn(2))
	}

	component := Index(words)

	err := component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template render error: %v", err)
	}
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	wordGroup := typeHandler(r)

	words := wordGroup.getRandomWords(2 + rand.Intn(2))
	w.Header().Set("words", words.getWordsString())

	component := WordsDiv(words)

	err := component.Render(context.Background(), w)

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

	wordGroup := typeHandler(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	words := wordGroup.getRandomWords(1)
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

func main() {
	re := regexp.MustCompile((`,\r?\n`))
	fantasywords = WordGroup{words: re.Split(fantasytxt, -1)}
	scifiwords = WordGroup{words: re.Split(scifitxt, -1)}
	mysterywords = WordGroup{words: re.Split(mysterytxt, -1)}

	fantasywords.shuffleWords()
	scifiwords.shuffleWords()
	mysterywords.shuffleWords()

	wordGroupByType["fantasy"] = &fantasywords
	wordGroupByType["scifi"] = &scifiwords
	wordGroupByType["mystery"] = &mysterywords

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/words/", wordsHandler)
	http.HandleFunc("/word/{word}/", wordHandler)

	log.Println("Starting server at http://localhost:8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

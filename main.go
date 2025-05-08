package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type WordList int

const (
	Fantasy WordList = iota
	Scifi
	Mystery
	Fantasynames
	Fantasypics
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
	case Fantasypics:
		return "fantasypics"
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
	case "fantasypics":
		return Fantasypics, nil
	default:
		return Fantasy, errors.New("invalid word list: " + s)
	}
}

//go:embed images.txt
var imagestxt string
var images []string

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
	Id      int    `json:"id"`
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

	return url.QueryEscape(strings.Join(words, ","))
}

func (w Words) getWordsHypenated() (wordString string) {
	words := make([]string, 0, len(w.Words))
	for _, word := range w.Words {
		words = append(words, word.Word)
	}

	sort.Strings(words)

	return strings.Join(words, "-")
}

func typeHandler(r *http.Request) WordList {
	wordsType := r.URL.Query().Get("type")

	wordList, err := ParseWordList(wordsType)

	if err != nil {
		log.Println("defaulting to fantasy")
	}

	return wordList
}

func handler(w http.ResponseWriter, r *http.Request) {
	savedWords := r.URL.Query().Get("words")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var words Words
	var err error

	wordCategory := typeHandler(r)

	if savedWords != "" {
		words, err = getSavedWords(savedWords, wordCategory)
	} else {
		words, err = getRandomWords(2+rand.Intn(2), wordCategory)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := Index(words, wordCategory)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template render error: %v", err)
	}
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {

	wordCategory := typeHandler(r)

	var words Words
	var err error

	if wordCategory == Fantasypics {
		savedWord := images[rand.Intn(len(images))]
		savedWord = strings.ReplaceAll(savedWord, "-", ",")
		words, err = getSavedWords(savedWord, wordCategory)
	} else {
		words, err = getRandomWords(2+rand.Intn(2), wordCategory)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		for i := range words.Words {
			words.Words[i].Id = i
		}
		json, err := json.Marshal(words)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(json)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	wordCategory := typeHandler(r)

	words, err := getRandomWords(1, wordCategory)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		json, err := json.Marshal(words.Words[0])

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(json)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if savedWordsSplit[0] == "" {
		savedWordsSplit = []string{"", "", ""}
	}
	savedWordsSplit[idx] = words.Words[0].Word

	var newWords Words
	newWords.Words = make([]Word, 0, len(savedWordsSplit))
	for i := range savedWordsSplit {
		if i == idx {
			newWords.Words = append(newWords.Words, words.Words[0])
		} else {
			newWords.Words = append(newWords.Words, Word{Word: savedWordsSplit[i]})
		}
	}

	w.Header().Set("words", newWords.getWordsString())

	component := WordDiv(words.Words[0], wordIdx)

	err = component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, "Failed to render word", http.StatusInternalServerError)
	}
}

func main() {
	err := initDataaccess()

	if err != nil {
		log.Panicf("failed to init dataaccess, err: %v\n", err)
	}
	imagestxt = strings.ReplaceAll(imagestxt, "\r\n", "\n")
	images = strings.Split(imagestxt, "\n")

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("GET /images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", handler)
	http.HandleFunc("/words/", wordsHandler)
	http.HandleFunc("/word/{word}/", wordHandler)

	log.Println("Starting server at http://localhost:8001")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

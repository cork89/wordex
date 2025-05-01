package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	dataaccess "com.github.cork89/goodmorning/db"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

var queries *dataaccess.Queries

func initDataaccess() error {
	ctx := context.Background()
	dbFile := "file:wordex.db"
	// conn, err := sql.Open("sqlite3", ":memory:")
	conn, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		return err
	}

	// create tables
	if _, err := conn.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries = dataaccess.New(conn)

	return nil
}

func insertWords(words Words, category WordList) error {
	var err error
	for _, v := range words.Words {
		_, err = queries.CreateWord(context.Background(), dataaccess.CreateWordParams{Word: v.Word, Category: category.String(), Subtext: v.Subtext})

		if err != nil {
			fmt.Println("failed to insert word, err: ", err)
		}
	}
	return err
}

func getRandomWords(count int, category WordList) (Words, error) {
	var dbwords []dataaccess.Word
	var err error
	if count > 0 && count < 4 {
		dbwords, err = queries.GetRandomWords(context.Background(), dataaccess.GetRandomWordsParams{Category: category.String(), Limit: int64(count)})
	} else {
		return Words{}, errors.New("count must be 1, 2, or 3")
	}

	if err != nil {
		fmt.Println("failed to retrieve random words, count:", count, " err:", err)
		return Words{}, err
	}

	words := make([]Word, 0, len(dbwords))
	colorsIdx := rand.Intn(colorsLen)
	for i, v := range dbwords {
		words = append(words, Word{Word: v.Word, Subtext: v.Subtext, Color: colors[(colorsIdx+i)%colorsLen]})
	}
	return Words{words}, nil
}

func wordExists(words []Word, savedWord string) bool {
	for _, v := range words {
		if v.Word == savedWord {
			return true
		}
	}
	return false
}

func getSavedWords(savedWordsStr string, category WordList) (Words, error) {
	var dbwords []dataaccess.Word
	var err error
	savedWords := strings.Split(savedWordsStr, ",")

	if len(savedWords) == 2 {
		dbwords, err = queries.GetTwoSavedWords(context.Background(), dataaccess.GetTwoSavedWordsParams{Word: savedWords[0], Word_2: savedWords[1], Category: category.String()})
	} else if len(savedWords) == 3 {
		dbwords, err = queries.GetThreeSavedWords(context.Background(), dataaccess.GetThreeSavedWordsParams{Word: savedWords[0], Word_2: savedWords[1], Word_3: savedWords[2], Category: category.String()})
	} else {
		return Words{}, errors.New("invalid number of saved words")
	}

	if err != nil {
		fmt.Println("failed to retrieve saved words, err:", err)
		return Words{}, err
	}

	words := make([]Word, 0, len(dbwords))
	colorsIdx := rand.Intn(colorsLen)
	for i, v := range dbwords {
		words = append(words, Word{Word: v.Word, Subtext: v.Subtext, Color: colors[(colorsIdx+i)%colorsLen]})
	}
	if len(dbwords) != len(savedWords) {
		for _, v := range savedWords {
			if !wordExists(words, v) {
				words = append(words, Word{Word: v, Color: colors[rand.Intn(colorsLen)]})
			}
		}
	}
	return Words{words}, nil
}

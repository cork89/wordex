package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"

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

func insertWords(wordGroup WordGroup, category WordList) error {
	var err error
	for _, v := range wordGroup.words {
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
	for _, v := range dbwords {
		words = append(words, Word{Word: v.Word, Subtext: v.Subtext, Color: colors[rand.Intn(colorsLen)]})
	}
	return Words{words}, nil
}

func getSavedWords(savedWords string, category WordList) (Words, error) {
	dbwords, err := queries.GetSavedWords(context.Background(), dataaccess.GetSavedWordsParams{Word: savedWords, Category: category.String()})

	if err != nil {
		fmt.Println("failed to retrieve saved words")
		return Words{}, err
	}

	words := make([]Word, 0, len(dbwords))
	for _, v := range dbwords {
		words = append(words, Word{Word: v.Word, Subtext: v.Subtext, Color: colors[rand.Intn(colorsLen)]})
	}
	return Words{words}, nil
}

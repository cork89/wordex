// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
)

const createWord = `-- name: CreateWord :one
INSERT INTO words (
  word, category, subtext
) VALUES (
  ?, ?, ?
)
RETURNING id, word, category, subtext
`

type CreateWordParams struct {
	Word     string
	Category string
	Subtext  string
}

func (q *Queries) CreateWord(ctx context.Context, arg CreateWordParams) (Word, error) {
	row := q.db.QueryRowContext(ctx, createWord, arg.Word, arg.Category, arg.Subtext)
	var i Word
	err := row.Scan(
		&i.ID,
		&i.Word,
		&i.Category,
		&i.Subtext,
	)
	return i, err
}

const getRandomWords = `-- name: GetRandomWords :many
SELECT id, word, category, subtext FROM words
WHERE category = ?
ORDER BY RANDOM()
LIMIT ?
`

type GetRandomWordsParams struct {
	Category string
	Limit    int64
}

func (q *Queries) GetRandomWords(ctx context.Context, arg GetRandomWordsParams) ([]Word, error) {
	rows, err := q.db.QueryContext(ctx, getRandomWords, arg.Category, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Word
	for rows.Next() {
		var i Word
		if err := rows.Scan(
			&i.ID,
			&i.Word,
			&i.Category,
			&i.Subtext,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThreeSavedWords = `-- name: GetThreeSavedWords :many
SELECT id, word, category, subtext FROM words
WHERE word in (?,?,?) and category = ?
LIMIT 3
`

type GetThreeSavedWordsParams struct {
	Word     string
	Word_2   string
	Word_3   string
	Category string
}

func (q *Queries) GetThreeSavedWords(ctx context.Context, arg GetThreeSavedWordsParams) ([]Word, error) {
	rows, err := q.db.QueryContext(ctx, getThreeSavedWords,
		arg.Word,
		arg.Word_2,
		arg.Word_3,
		arg.Category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Word
	for rows.Next() {
		var i Word
		if err := rows.Scan(
			&i.ID,
			&i.Word,
			&i.Category,
			&i.Subtext,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTwoSavedWords = `-- name: GetTwoSavedWords :many
SELECT id, word, category, subtext FROM words
WHERE word in (?,?) and category = ?
LIMIT 2
`

type GetTwoSavedWordsParams struct {
	Word     string
	Word_2   string
	Category string
}

func (q *Queries) GetTwoSavedWords(ctx context.Context, arg GetTwoSavedWordsParams) ([]Word, error) {
	rows, err := q.db.QueryContext(ctx, getTwoSavedWords, arg.Word, arg.Word_2, arg.Category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Word
	for rows.Next() {
		var i Word
		if err := rows.Scan(
			&i.ID,
			&i.Word,
			&i.Category,
			&i.Subtext,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

-- name: GetRandomWords :many
SELECT * FROM words
WHERE category = ?
ORDER BY RANDOM()
LIMIT ?;

-- name: GetTwoSavedWords :many
SELECT * FROM words
WHERE word in (?,?) and category = ?
LIMIT 2;

-- name: GetThreeSavedWords :many
SELECT * FROM words
WHERE word in (?,?,?) and category = ?
LIMIT 3;

-- name: CreateWord :one
INSERT INTO words (
  word, category, subtext
) VALUES (
  ?, ?, ?
)
RETURNING *;

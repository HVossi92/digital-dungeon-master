// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
)

type ChatMessage struct {
	ID         int64
	Role       string
	Content    string
	SaveGameID int64
	Created    sql.NullString
}

type SaveGame struct {
	ID      int64
	Name    string
	Created sql.NullString
}

type Setting struct {
	ID  int64
	Url string
	Llm string
}

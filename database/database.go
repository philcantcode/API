package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func init() {
	database, _ = sql.Open("sqlite3", "E:\\Google Drive\\elements.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS watchFolders (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE)")
	statement.Exec()
}

type TrackFolders struct {
	ID   int
	Path string
}

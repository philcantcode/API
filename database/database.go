package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philcantcode/goApi/utils"
)

var database *sql.DB

func init() {
	database, _ = sql.Open("sqlite3", utils.DBLoc)

	statement, _ := database.Prepare(
		"CREATE TABLE IF NOT EXISTS watchFolders" +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE)")
	statement.Exec()

	statement, _ = database.Prepare(
		"CREATE TABLE IF NOT EXISTS playHistory " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" name TEXT NOT NULL, " +
			" hash TEXT NOT NULL, " +
			" path TEXT NOT NULL, " +
			" playTime TEXT, " +
			" date INTEGER NOT NULL)")
	statement.Exec()

	statement, _ = database.Prepare(
		"CREATE TABLE IF NOT EXISTS folderMeta " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" type INTEGER)")
	statement.Exec()

	statement, _ = database.Prepare(
		"CREATE TABLE IF NOT EXISTS fileTrack " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" dateAdded TEXT NOT NULL)")
	statement.Exec()

	statement, _ = database.Prepare(
		"CREATE TABLE IF NOT EXISTS settings " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" key TEXT UNIQUE, " +
			" value TEXT)")
	statement.Exec()
}

// FindOrCreateMedia searches for or creates a media
func FindOrCreateMedia(path string) MediaInfo {
	media := SelectMedia(path)

	// Doesn't exist in DB
	if media.ID == 0 {
		InsertMedia(path)
		media = SelectMedia(path)
	}

	return media
}

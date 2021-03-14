package database

import (
	"database/sql"
	"fmt"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

var con *sql.DB

func init() {
	var dbLoc string
	os := runtime.GOOS

	switch os {
	case "windows":
		dbLoc = "E:/Google Drive/elements.db"
	case "darwin":
		dbLoc = "/Users/Phil/Google Drive/elements.db"
	case "linux":
		fmt.Println("OS Not Supported")
	default:
		fmt.Printf("%s.\n", os)
	}

	fmt.Printf("Database loaded: %s\n", dbLoc)

	con, _ = sql.Open("sqlite3", dbLoc)

	statement, _ := con.Prepare(
		"CREATE TABLE IF NOT EXISTS watchFolders" +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE)")
	statement.Exec()

	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS playHistory " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" name TEXT NOT NULL, " +
			" hash TEXT NOT NULL, " +
			" path TEXT NOT NULL, " +
			" playTime TEXT, " +
			" date INTEGER NOT NULL)")
	statement.Exec()

	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS folderMeta " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" type INTEGER)")
	statement.Exec()

	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS fileTrack " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" dateAdded TEXT NOT NULL)")
	statement.Exec()

	statement, _ = con.Prepare(
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

package database

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/philcantcode/goApi/utils"
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
			" hash TEXT DEFAULT '', " +
			" altHash TEXT DEFAULT '', " +
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

	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS ffmpeg " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" archivePath TEXT, " +
			" mp4Path TEXT, " +
			" codecs TEXT, " +
			" conversions TEXT, " +
			" duration TEXT, " +
			" date INTEGER NOT NULL)")
	statement.Exec()

	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS ffmpegPriority " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT UNIQUE)")
	statement.Exec()

	// Convert dates created in Java (System.currentTimeMillis()) which are
	// 1000x larger than golang times to be compatible
	statement, _ = con.Prepare(
		"UPDATE playHistory SET date = (date / 1000) WHERE date > 1000000000000")
	statement.Exec()
}

// FindOrCreateMedia searches for or creates a media
func FindOrCreateMedia(path string) MediaInfo {

	// Try find by path first
	mediaInfo, err := SelectMediaByPath(path)

	if err == nil {
		return mediaInfo
	}

	// Try find by computing hash
	fmt.Println("Computing File Hash, Please Wait")
	hash, _ := utils.Hash(path)
	mediaInfo, err = SelectMediaByHash(hash)

	// Hash Found
	if err == nil {
		UpdateMediaPathByHash(mediaInfo.File.AbsPath, hash)
		return mediaInfo
	}

	// Doesn't exist in DB
	id := InsertMedia(utils.ProcessFile(path))
	mediaInfo, err = SelectMediaByID(id)

	if err != nil {
		log.Fatal("Couldn't find or create media in playHistory")
	}

	return mediaInfo
}

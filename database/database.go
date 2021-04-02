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
		dbLoc = "E:/Google Drive/MediaPlayer.db"
	case "darwin":
		dbLoc = "/Users/Phil/Google Drive/MediaPlayer.db"
	case "linux":
		fmt.Println("OS Not Supported")
	default:
		fmt.Printf("%s.\n", os)
	}

	fmt.Printf("Database loaded: %s\n", dbLoc)

	con, _ = sql.Open("sqlite3", dbLoc)

	// Top level directories to keep track of
	statement, _ := con.Prepare(
		"CREATE TABLE IF NOT EXISTS RootDirectories" +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT UNIQUE)")
	statement.Exec()

	// The playhistory for media files
	// The hash is the fist computed hash for the file, the alt hash is used if a
	// conversion takes place resulting in an old hash and a new hash for the two
	// versions of the file
	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS MediaPlayback " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" playtime TEXT, " +
			" modified INTEGER)")
	statement.Exec()

	// Depricated database only used by Java
	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS folderMeta " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" type INTEGER)")
	statement.Exec()

	// Depricated database only used by Java
	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS fileTrack " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT NOT NULL, " +
			" dateAdded TEXT NOT NULL)")
	statement.Exec()

	// Depricated database only used by Java
	statement, _ = con.Prepare(
		"CREATE TABLE IF NOT EXISTS settings " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" key TEXT UNIQUE, " +
			" value TEXT)")
	statement.Exec()

	// Keeps track of file conversions using FFMPEG
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

	// Folders the user has manually marked has high priority for conversion
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

// FindOrCreateMedia searches for or creates a media by a given path
func FindOrCreateMedia(path string) MediaInfo {
	// Try find by path first
	mediaInfo, err := SelectMediaByPath(path)

	if err == nil {
		return mediaInfo
	}

	// Try find by computing hash
	fmt.Println("[FindOrCreateMedia] Computing File Hash, Please Wait")
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

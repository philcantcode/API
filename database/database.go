package database

import (
	"database/sql"
	"fmt"
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
	stmt, err := con.Prepare(
		"CREATE TABLE IF NOT EXISTS RootDirectories" +
			"(path TEXT PRIMARY KEY UNIQUE)")
	stmt.Exec()
	utils.Error("Couldn't create Database: RootDirectories", err)

	// Where each media is up to in its payback and date it was modified
	stmt, err = con.Prepare(
		"CREATE TABLE IF NOT EXISTS MediaPlayback " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" playtime INTEGER DEFAULT 0, " +
			" modified INTEGER DEFAULT -1)")
	stmt.Exec()
	utils.Error("Couldn't create Database: MediaPlayback", err)

	// Hashes of files point to a media playback so multiple files can have
	// be in the same playback time, e.g., the .mp4 and .avi files can still
	// playback at the same time
	stmt, err = con.Prepare(
		"CREATE TABLE IF NOT EXISTS MediaHashes " +
			"(id INTEGER PRIMARY KEY UNIQUE, " +
			" hash TEXT NOT NULL, " +
			" path TEXT NOT NULL, " +
			" mediaPlaybackID INTEGER NOT NULL)")
	stmt.Exec()
	utils.Error("Couldn't create Database: MediaHashes", err)

	// Keeps track of file conversions using FFMPEG
	stmt, err = con.Prepare(
		"CREATE TABLE IF NOT EXISTS FfmpegConversions " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" path TEXT, " + // The original location of the file
			" archivePath TEXT, " + // Where the file was moved to
			" originalCodecs TEXT, " + // Original Audio / Video format
			" convertedCodecs TEXT, " + // Conversion Audio / Video Format
			" duration TEXT, " + // Duration of how long it took
			" date INTEGER NOT NULL)") // When the conversion occured
	stmt.Exec()
	utils.Error("Couldn't create Database: FfmpegConversions", err)

	// Folders the user has manually marked has high priority for conversion
	stmt, err = con.Prepare(
		"CREATE TABLE IF NOT EXISTS FfmpegPriority " +
			"(path TEXT PRIMARY KEY UNIQUE)")
	stmt.Exec()
	utils.Error("Couldn't create Database: FfmpegPriority", err)

	stmt.Close()
}

// FindOrCreatePlayback takes in a path and either finds or creates
// the playback in the database
// Step 1: Hash the file
// Step 2: Check if the file has been previously hashed and find playback
// Step 3: If not, create it
func FindOrCreatePlayback(path string) Playback {
	// Fast lookup to see if the path exists in DB
	mediaPlaybackID, _ := SelectPlaybackID_ByPath(path)

	if mediaPlaybackID != -1 {
		playback := SelectMediaPlayback_ByID(mediaPlaybackID)
		playback.PrefLoc, _ = GetPreferredLocation(playback)
		return playback
	}

	// Slower hash check to see if the hash exists
	fmt.Printf("FindOrCreatePlayback hashing %s, please wait...\n", path)
	hash, err := utils.MD5Hash(path)
	fmt.Printf("FindOrCreatePlayback done: %s\n", hash)
	utils.Error("Couldn't MD5Hash: "+path, err)
	mediaPlaybackID, err = SelectPlaybackID_ByHash(hash)

	if mediaPlaybackID != -1 {
		return SelectMediaPlayback_ByID(mediaPlaybackID)
	}

	// Couldn't find a hash or path entry in the DB so create it
	mediaPlaybackID = InsertMediaPlayback()
	InsertMediaHash(hash, path, mediaPlaybackID)

	playback := SelectMediaPlayback_ByID(mediaPlaybackID)
	playback.PrefLoc, _ = GetPreferredLocation(playback)

	return playback
}

// Returns the preferred (loaded) media location where there are
// multiple files on disk
func GetPreferredLocation(playback Playback) (int, error) {
	for i := 0; i < len(playback.Locations); i++ {
		if playback.Locations[i].Exists {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no matching file exists on disk")
}

func DatabaseStats() {
	fmt.Printf("Num Connections: %d\n", con.Stats().OpenConnections)
	fmt.Printf("Num InUse Connections: %d\n", con.Stats().InUse)
}

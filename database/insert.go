package database

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

func InsertMediaHash(hash string, path string, mediaPlaybackID int) {
	stmt, err := con.Prepare("INSERT INTO `MediaHashes`" +
		"(`hash`, `path`, `mediaPlaybackID`) VALUES (?, ?, ?);")

	utils.Error("Couldn't Insert Into MediaHashes", err)

	_, err = stmt.Exec(hash, path, mediaPlaybackID)
	utils.Error("Results error from InsertMediaHash", err)

	stmt.Close()
}

// InsertMediaPlayback generates a new playback tracker and returns the ID
func InsertMediaPlayback() int {
	stmt, err := con.Prepare("INSERT INTO `MediaPlayback`" +
		"(`playtime`, `modified`) VALUES (?, ?);")

	utils.Error("Couldn't Insert Into MediaPlayback", err)

	res, err := stmt.Exec(0, time.Now().Unix())
	utils.Error("Results error from InsertMediaPlayback", err)

	insertID, _ := res.LastInsertId()
	utils.Error("LastInsertId error from InsertMediaPlayback", err)

	stmt.Close()
	return int(insertID)
}

////////////////////////////////////////////////////////////////////////////////

// Adds a root directory to track (G:\\, F:\\, etc)
func InsertRootDirectory(folder string) {
	stmt, err := con.Prepare("INSERT INTO `RootDirectories` (path) VALUES (?);")
	stmt.Exec(folder)
	utils.Error("Couldn't insert into RootDirectories", err)

	stmt.Close()
}

func InsertMedia(file utils.File) int {
	stmt, err := con.Prepare("INSERT INTO `playHistory`" +
		"(path, hash, altHash, playTime, date) VALUES (?, ?, ?, ?, ?)")

	res, _ := stmt.Exec(file.AbsPath, "", "", 0, time.Now().Unix())
	utils.Error("Couldn't insert into PlayHistory", err)

	insertID, _ := res.LastInsertId()

	stmt.Close()
	return int(insertID)
}

func InsertFfmpeg(archivePath string, mp4Path string, codecs string, conversions string, duration string) {
	stmt, err := con.Prepare("INSERT INTO `ffmpeg`" +
		"(archivePath, mp4Path, codecs, conversions, duration, date) VALUES (?, ?, ?, ?, ?, ?)")

	stmt.Exec(archivePath, mp4Path, codecs, conversions, duration, time.Now().Unix())
	utils.Error("Couldn't insert into Ffmpeg", err)

	stmt.Close()
}

func InsertFfmpegPriority(path string) {
	stmt, err := con.Prepare("INSERT INTO `ffmpegPriority` (path) VALUES (?)")
	stmt.Exec(path)
	stmt.Close()

	utils.Error("Couldn't insert into FfmpegPriority", err)
}

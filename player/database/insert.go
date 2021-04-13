package database

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

func InsertMediaHash(hash string, path string, mediaPlaybackID int) {
	stmt, err := playerCon.Prepare("INSERT INTO `MediaHashes`" +
		"(`hash`, `path`, `mediaPlaybackID`) VALUES (?, ?, ?);")

	utils.Error("Couldn't Insert Into MediaHashes", err)

	_, err = stmt.Exec(hash, path, mediaPlaybackID)
	utils.Error("Results error from InsertMediaHash", err)

	stmt.Close()
}

// InsertMediaPlayback generates a new playback tracker and returns the ID
func InsertMediaPlayback() int {
	stmt, err := playerCon.Prepare("INSERT INTO `MediaPlayback`" +
		"(`playtime`, `modified`) VALUES (?, ?);")

	utils.Error("Couldn't Insert Into MediaPlayback", err)

	res, err := stmt.Exec(0, 0)
	utils.Error("Results error from InsertMediaPlayback", err)

	insertID, _ := res.LastInsertId()
	utils.Error("LastInsertId error from InsertMediaPlayback", err)

	stmt.Close()
	return int(insertID)
}

// Adds a root directory to track (G:\\, F:\\, etc)
func InsertRootDirectory(folder string) {
	stmt, err := playerCon.Prepare("INSERT INTO `RootDirectories` (path) VALUES (?);")
	stmt.Exec(folder)
	utils.Error("Couldn't insert into RootDirectories", err)

	stmt.Close()
}

func InsertMedia(file utils.File) int {
	stmt, err := playerCon.Prepare("INSERT INTO `playHistory`" +
		"(path, hash, altHash, playTime, date) VALUES (?, ?, ?, ?, ?)")

	res, _ := stmt.Exec(file.AbsPath, "", "", 0, time.Now().Unix())
	utils.Error("Couldn't insert into PlayHistory", err)

	insertID, _ := res.LastInsertId()

	defer stmt.Close()
	return int(insertID)
}

func InsertFfmpeg(mp4Path string, archivePath string, codecs string, convCodecs string, duration string) {
	stmt, err := playerCon.Prepare("INSERT INTO `FfmpegConversions`" +
		"(`path`, `archivePath`, `originalCodecs`, `convertedCodecs`, `duration`, `date`) VALUES (?, ?, ?, ?, ?, ?)")

	stmt.Exec(mp4Path, archivePath, codecs, convCodecs, duration, time.Now().Unix())
	utils.Error("Couldn't insert into Ffmpeg", err)

	defer stmt.Close()
}

func InsertFfmpegPriority(path string) {
	stmt, err := playerCon.Prepare("INSERT INTO `ffmpegPriority` (path) VALUES (?)")
	stmt.Exec(path)
	defer stmt.Close()

	utils.Error("Couldn't insert into FfmpegPriority", err)
}

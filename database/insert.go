package database

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

func InsertFolder(folder string) {
	statement, _ := con.Prepare("INSERT INTO watchFolders (path)" +
		"VALUES (?);")

	statement.Exec(folder)
}

func InsertMedia(file utils.File) int {
	statement, _ := con.Prepare("INSERT INTO `playHistory`" +
		"(path, hash, altHash, playTime, date) VALUES (?, ?, ?, ?, ?)")

	res, _ := statement.Exec(file.AbsPath, "", "", 0, time.Now().Unix())
	insertID, _ := res.LastInsertId()

	return int(insertID)
}

func InsertFfmpeg(archivePath string, mp4Path string, codecs string, conversions string, duration string) {
	statement, _ := con.Prepare("INSERT INTO `ffmpeg`" +
		"(archivePath, mp4Path, codecs, conversions, duration, date) VALUES (?, ?, ?, ?, ?, ?)")

	statement.Exec(archivePath, mp4Path, codecs, conversions, duration, time.Now().Unix())
}

func InsertFfmpegPriority(path string) {
	statement, _ := con.Prepare("INSERT INTO `ffmpegPriority`" +
		"(path) VALUES (?)")

	statement.Exec(path)
}

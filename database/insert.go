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

func InsertMedia(path string) {
	statement, _ := con.Prepare("INSERT INTO `playHistory`" +
		"(name, hash, path, playTime, date) VALUES (?, ?, ?, ?, ?)")

	f := utils.ProcessFile(path)
	name := f.Name + f.Ext
	path = f.Path + name

	statement.Exec(name, "", path, 0, time.Now().Unix())
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

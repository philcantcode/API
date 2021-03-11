package database

import (
	"github.com/philcantcode/goApi/utils"
)

func InsertFolder(folder string) {
	statement, _ := database.Prepare("INSERT INTO watchFolders (path)" +
		"VALUES (?);")

	statement.Exec(folder)
}

func InsertMedia(path string) {
	statement, _ := database.Prepare("INSERT INTO `playHistory`" +
		"(name, hash, path, playTime, date) VALUES (?, ?, ?, ?, ?)")

	statement.Exec(utils.ExtractFileName(path), "", path, 0, utils.GetTime())
}

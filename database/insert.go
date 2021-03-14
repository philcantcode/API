package database

import "github.com/philcantcode/goApi/utils"

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

	statement.Exec(name, "", path, 0, utils.GetTime())
}

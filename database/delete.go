package database

import (
	"github.com/philcantcode/goApi/utils"
)

func DeleteRootDirectory(path string) {
	stmt, err := con.Prepare("DELETE FROM `RootDirectories` WHERE `path` = ?;")
	stmt.Exec(path)
	stmt.Close()

	utils.Error("Couldn't Delete From RootDirectory", err)
}

func DeleteFfmpegEntry(path string) {
	stmt, err := con.Prepare("DELETE FROM `FfmpegConversions` WHERE `originalPath` = ? OR `archivePath` = ?;")
	stmt.Exec(path, path)
	stmt.Close()

	utils.Error("Couldn't Delete From FfmpegHistory", err)
}

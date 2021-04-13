package database

import (
	"github.com/philcantcode/goApi/utils"
)

func DeleteRootDirectory(path string) {
	stmt, err := playerCon.Prepare("DELETE FROM `RootDirectories` WHERE `path` = ?;")
	stmt.Exec(path)
	stmt.Close()

	utils.Error("Couldn't Delete From RootDirectory", err)
}

func DeleteFfmpegEntry(id int) {
	stmt, err := playerCon.Prepare("DELETE FROM `FfmpegConversions` WHERE `id` = ?;")
	stmt.Exec(id)
	stmt.Close()

	utils.Error("Couldn't Delete From FfmpegHistory", err)
}

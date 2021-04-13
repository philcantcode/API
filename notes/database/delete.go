package notes

import "github.com/philcantcode/goApi/utils"

func DeleteNote(id int) {
	stmt, err := wikiCon.Prepare("DELETE FROM `Notes` WHERE `id` = ?;")
	stmt.Exec(id)
	defer stmt.Close()

	utils.Error("Couldn't Delete From FfmpegHistory", err)
}

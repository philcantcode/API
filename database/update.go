package database

import (
	"github.com/philcantcode/goApi/utils"
)

// UpdatePlaytime updates the playtime for a media ID
func UpdatePlaytime(id int, playtime int) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET playTime = ?," +
			"date = ? WHERE id = ?;")

	stmt.Exec(playtime, utils.GetTime(), id)
}

func UpdateMediaName(oldName string, newName string) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET name = ? WHERE name = ?;")

	stmt.Exec(newName, oldName)
}

func UpdateMediaPath(oldPath string, newPath string) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET path = ? WHERE path = ?;")

	stmt.Exec(newPath, oldPath)
}

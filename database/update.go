package database

import (
	"github.com/philcantcode/goApi/utils"
)

// UpdatePlaytime updates the playtime for a media ID
func UpdatePlaytime(id int, playtime int) {
	stmt, _ := database.Prepare(
		"UPDATE `playHistory` SET playTime = ?," +
			"date = ? WHERE id = ?;")

	stmt.Exec(playtime, utils.GetTime(), id)
}

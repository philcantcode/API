package database

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

// UpdatePlaytime updates the playtime for a media ID
func UpdatePlaytime(id int, playtime int) {
	stmt, err := playerCon.Prepare("UPDATE `MediaPlayback` SET `playtime` = ?, `modified` = ? WHERE `id` = ?;")
	defer stmt.Close()
	utils.Error("Couldn't update UpdatePlaytime", err)

	_, err2 := stmt.Exec(playtime, time.Now().Unix(), id)

	utils.Error("Results error from UpdatePlaytime", err2)
	stmt.Close()
}

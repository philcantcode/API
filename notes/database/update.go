package notes

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

// UpdateNote updates the playtime for a media ID
func UpdateNote(id int, keyword string, desc string, text string) {
	stmt, err := wikiCon.Prepare("UPDATE `Notes` SET `keyword` = ?, `desc` = ?, `text` = ?, `modified` = ? WHERE `id` = ?;")

	utils.Error("Couldn't update UpdatePlaytime", err)
	defer stmt.Close()

	_, err2 := stmt.Exec(keyword, desc, text, time.Now().Unix(), id)

	utils.Error("Results error from UpdatePlaytime", err2)
	stmt.Close()
}

package notes

import (
	"time"

	"github.com/philcantcode/goApi/utils"
)

// InsertNote creates a new note with the current timestamp
func InsertNote(keyword string, desc string, text string) (int, error) {
	stmt, err := wikiCon.Prepare("INSERT INTO `Notes` " +
		"(`keyword`, `desc`, `text`, `modified`) VALUES (?, ?, ?, ?);")

	utils.Error("Couldn't Prepare Insert Into Notes", err)
	defer stmt.Close()

	res, err := stmt.Exec(keyword, desc, text, time.Now().Unix())
	utils.Error("Results error from InsertNote", err)

	insertID, err := res.LastInsertId()
	utils.Error("LastInsertId error from InsertNote", err)

	return int(insertID), nil
}

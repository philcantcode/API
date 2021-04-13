package notes

import (
	"database/sql"
	"fmt"
	"runtime"

	_ "github.com/mattn/go-sqlite3"

	"github.com/philcantcode/goApi/utils"
)

var wikiCon *sql.DB

func init() {
	var dbLoc string
	os := runtime.GOOS

	switch os {
	case "windows":
		dbLoc = "E:/Google Drive/WikiNotes.db"
	case "darwin":
		dbLoc = "/Users/Phil/Google Drive/WikiNotes.db"
	case "linux":
		fmt.Println("OS Not Supported")
	default:
		fmt.Printf("%s.\n", os)
	}

	wikiCon, _ = sql.Open("sqlite3", dbLoc)
	fmt.Printf("Database loaded: %s\n", dbLoc)

	// Where each media is up to in its payback and date it was modified
	stmt, err := wikiCon.Prepare(
		"CREATE TABLE IF NOT EXISTS Notes " +
			"(id INTEGER PRIMARY KEY AUTOINCREMENT," +
			" keyword TEXT DEFAULT '', " +
			" desc TEXT DEFAULT '', " +
			" text TEXT DEFAULT '', " +
			" modified INTEGER DEFAULT -1)")
	stmt.Exec()

	defer stmt.Close()
	utils.Error("Couldn't create Database: Notes", err)
}

package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func init() {
	database, _ = sql.Open("sqlite3", "./database/mediaPlayer.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS trackedFolders (id INTEGER PRIMARY KEY AUTOINCREMENT, folder TEXT UNIQUE)")
	statement.Exec()
}

func TrackFolder(folder string) {
	statement, _ := database.Prepare("INSERT INTO trackedFolders (folder) VALUES (?);")
	statement.Exec(folder)
}

func UnTrackFolder(folder string) {
	statement, _ := database.Prepare("DELETE FROM trackedFolders WHERE folder = ?;")
	statement.Exec(folder)
}

type TrackFolders struct {
	ID     int
	Folder string
}

func GetTrackedFolders() []TrackFolders {
	rows, _ := database.Query("SELECT * FROM trackedFolders;")

	var res []TrackFolders

	for rows.Next() {
		var id int
		var folder string

		rows.Scan(&id, &folder)
		res = append(res, TrackFolders{ID: id, Folder: folder})
	}

	return res
}

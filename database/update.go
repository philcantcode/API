package database

import (
	"log"
	"time"

	"github.com/philcantcode/goApi/utils"
)

// UpdatePlaytime updates the playtime for a media ID
func UpdatePlaytime(id int, playtime int) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET playTime = ?, date = ? WHERE id = ?;")

	stmt.Exec(playtime, time.Now().Unix(), id)
}

// UpdateMediaPath changes the name of the path in the playHistory DB
// Mainly used during conversion when the file extensions change
func UpdateMediaPath(oldPath utils.File, newPath utils.File) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET path = ? WHERE path = ?;")

	res, err := stmt.Exec(newPath.AbsPath, oldPath.AbsPath)
	rowsUpdated, _ := res.RowsAffected()

	if err != nil || rowsUpdated == 0 {
		log.Fatalf("UpdateMediaPath Failed, Couldn't Change: %s to, %s", oldPath.AbsPath, newPath.AbsPath)
	}
}

func UpdateMediaPathByHash(newPath string, hash string) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET path = ? WHERE hash = ? OR altHash = ?;")

	res, err := stmt.Exec(newPath, hash, hash)
	rowsUpdated, _ := res.RowsAffected()

	if err != nil || rowsUpdated == 0 {
		log.Fatalf("UpdateMediaPath Failed, Couldn't Change: %s to, %s", hash, newPath)
	}
}

func UpdateMediaHash(path string, hash string) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET hash = ? WHERE path = ?;")

	stmt.Exec(hash, path)
}

func UpdateMediaAltHash(path string, hash string) {
	stmt, _ := con.Prepare(
		"UPDATE `playHistory` SET `altHash` = ? WHERE `path` = ?;")

	stmt.Exec(hash, path)
}

package database

import (
	"github.com/philcantcode/goApi/utils"
)

// Directory is a Top level directory
type Directory struct {
	ID   int
	Path string
}

// SelectDirectories returns all the locations monitored on disk
func SelectDirectories() []Directory {
	rows, _ := con.Query("SELECT * FROM watchFolders ORDER BY id DESC;")
	var res []Directory

	for rows.Next() {
		var id int
		var folder string

		rows.Scan(&id, &folder)
		res = append(res, Directory{ID: id, Path: folder})
	}

	return res
}

// MediaInfo is the default struct for a database item
type MediaInfo struct {
	ID       int
	Hash     string
	PlayTime int
	Date     int
	File     utils.File
}

// SelectMedia finds the playback in the database
func SelectMedia(path string) MediaInfo {
	stmt, _ := con.Prepare(
		"SELECT `id`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `path` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(path)

	for rows.Next() {
		var dbPath string

		rows.Scan(&media.ID, &media.Hash,
			&dbPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(dbPath)
	}

	return media
}

// SelectMediaByID finds the playback in the database
func SelectMediaByID(id int) MediaInfo {
	stmt, _ := con.Prepare(
		"SELECT `id`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `id` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(id)

	for rows.Next() {
		var dbPath string

		rows.Scan(&media.ID, &media.Hash,
			&dbPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(dbPath)
	}

	return media
}

// SelectMediaByTime finds the playback in the from a given timepoint
func SelectMediaByTime(unixTime int64) []MediaInfo {

	var mediaList []MediaInfo

	stmt, _ := con.Prepare(
		"SELECT `id`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `date` >= ?;")

	rows, _ := stmt.Query(unixTime)

	for rows.Next() {
		media := MediaInfo{}
		var dbPath string

		rows.Scan(&media.ID, &media.Hash,
			&dbPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(dbPath)
		mediaList = append(mediaList, media)
	}

	return mediaList
}

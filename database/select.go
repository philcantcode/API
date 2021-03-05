package database

import (
	"github.com/philcantcode/goApi/utils"
)

// Directory is a Top level directory
type Directory struct {
	ID   int
	Path string
}

// GetDirectories returns all the locations monitored on disk
func GetDirectories() []Directory {
	rows, _ := database.Query("SELECT * FROM watchFolders;")
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
	Title    string
	Hash     string
	Path     string
	Folder   string
	PlayTime int
	Date     int
}

// SelectMedia finds the playback in the database
func SelectMedia(path string) MediaInfo {
	stmt, _ := database.Prepare(
		"SELECT `id`, `name`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `path` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(path)

	for rows.Next() {
		rows.Scan(&media.ID, &media.Title, &media.Hash,
			&media.Path, &media.PlayTime, &media.Date)
	}

	media.Folder = utils.ExtractFolderName(media.Path)

	return media
}

// SelectMediaByID finds the playback in the database
func SelectMediaByID(id int) MediaInfo {
	stmt, _ := database.Prepare(
		"SELECT `id`, `name`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `id` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(id)

	for rows.Next() {
		rows.Scan(&media.ID, &media.Title, &media.Hash,
			&media.Path, &media.PlayTime, &media.Date)
	}

	media.Folder = utils.ExtractFolderName(media.Path)

	return media
}

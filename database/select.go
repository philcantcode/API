package database

type TrackFolders struct {
	ID   int
	Path string
}

// GetTrackedFolders returns all the locations monitored on disk
func GetTrackedFolders() []TrackFolders {
	rows, _ := database.Query("SELECT * FROM watchFolders;")
	var res []TrackFolders

	for rows.Next() {
		var id int
		var folder string

		rows.Scan(&id, &folder)
		res = append(res, TrackFolders{ID: id, Path: folder})
	}

	return res
}

// Media is the default struct for a database item
type Media struct {
	ID       int
	Title    string
	Hash     string
	Path     string
	PlayTime int
	Date     string
}

// SelectMedia finds the playback in the database
func SelectMedia(path string) Media {
	stmt, _ := database.Prepare(
		"SELECT `id`, `name`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `path` = ? LIMIT 1;")

	media := Media{}
	rows, _ := stmt.Query(path)

	for rows.Next() {
		rows.Scan(&media.ID, &media.Title, &media.Hash,
			&media.Path, &media.PlayTime, &media.Date)
	}

	return media
}

// SelectMediaByID finds the playback in the database
func SelectMediaByID(id int64) Media {
	stmt, _ := database.Prepare(
		"SELECT `id`, `name`, `hash`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `id` = ? LIMIT 1;")

	media := Media{}
	rows, _ := stmt.Query(id)

	for rows.Next() {
		rows.Scan(&media.ID, &media.Title, &media.Hash,
			&media.Path, &media.PlayTime, &media.Date)
	}

	return media
}

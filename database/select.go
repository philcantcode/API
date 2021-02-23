package database

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

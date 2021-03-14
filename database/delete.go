package database

func UnTrackFolder(folder string) {
	statement, _ := con.Prepare("DELETE FROM watchFolders WHERE path = ?;")
	statement.Exec(folder)
}

package database

func UnTrackFolder(folder string) {
	statement, _ := database.Prepare("DELETE FROM watchFolders WHERE path = ?;")
	statement.Exec(folder)
}

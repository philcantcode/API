package database

func TrackFolder(folder string) {
	statement, _ := database.Prepare("INSERT INTO watchFolders (path) VALUES (?);")
	statement.Exec(folder)
}

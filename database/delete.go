package database

func UnTrackFolder(folder string) {
	statement, _ := con.Prepare("DELETE FROM `watchFolders` WHERE `path` = ?;")
	statement.Exec(folder)
}

func DeleteFfmpegEntry(path string) {
	statement, _ := con.Prepare("DELETE FROM `ffmpeg` WHERE `archivePath` = ? OR `mp4Path` = ?;")
	statement.Exec(path, path)
}

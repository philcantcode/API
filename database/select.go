package database

import (
	"errors"

	"github.com/philcantcode/goApi/utils"
)

// SelectDirectories returns all the locations monitored on disk
func SelectDirectories() []utils.File {
	rows, _ := con.Query("SELECT * FROM watchFolders ORDER BY id DESC;")
	var res []utils.File

	for rows.Next() {
		var id int
		var folder string

		rows.Scan(&id, &folder)
		res = append(res, utils.ProcessFile(folder))
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

// SelectMediaByPath finds the playback in the database
func SelectMediaByPath(path string) (MediaInfo, error) {
	stmt, _ := con.Prepare(
		"SELECT `id`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `path` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(path)
	err := errors.New("Media not found by path in playHistory")

	for rows.Next() {
		var mediaPath string

		rows.Scan(&media.ID, &mediaPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(mediaPath)
		err = nil
	}

	return media, err
}

// SelectAllMedia finds the playback in the database
func SelectAllMedia() []MediaInfo {
	stmt, _ := con.Prepare(
		"SELECT `id`, `hash`, `path`, `playTime`, `date` FROM `playHistory` ORDER BY `id` ASC;")

	var mediaList []MediaInfo
	rows, _ := stmt.Query()

	for rows.Next() {
		media := MediaInfo{}
		var mediaPath string

		rows.Scan(&media.ID, &media.Hash, &mediaPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(mediaPath)
		mediaList = append(mediaList, media)
	}

	return mediaList
}

// SelectMediaHash finds the hash for a playHistory path
func SelectMediaHash(path string) string {
	stmt, _ := con.Prepare(
		"SELECT `hash` FROM `playHistory` WHERE `path` = ? LIMIT 1;")

	rows, _ := stmt.Query(path)
	var hash string

	for rows.Next() {
		rows.Scan(&hash)
	}

	return hash
}

// SelectMediaByID finds the playback in the database
func SelectMediaByID(id int) (MediaInfo, error) {
	stmt, _ := con.Prepare(
		"SELECT `id`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `id` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(id)
	err := errors.New("Media not found by ID in playHistory")

	for rows.Next() {
		var mediaPath string

		rows.Scan(&media.ID, &mediaPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(mediaPath)
		err = nil
	}

	return media, err
}

// SelectMediaByHash finds the playback in the database given a hash or alt hash
func SelectMediaByHash(hash string) (MediaInfo, error) {
	stmt, _ := con.Prepare(
		"SELECT `id`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `hash` = ? OR `altHash` = ? LIMIT 1;")

	media := MediaInfo{}
	rows, _ := stmt.Query(hash, hash)
	err := errors.New("Media not found by hash in playHistory")

	for rows.Next() {
		var mediaPath string

		rows.Scan(&media.ID, &mediaPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(mediaPath)
		err = nil
	}

	return media, err
}

// SelectMediaByTime finds the playback in the from a given timepoint
func SelectMediaByTime(unixTime int64) []MediaInfo {

	var mediaList []MediaInfo

	stmt, _ := con.Prepare(
		"SELECT `id`, `path`, `playTime`, `date`" +
			"FROM `playHistory` WHERE `date` >= ?;")

	rows, _ := stmt.Query(unixTime)

	for rows.Next() {
		media := MediaInfo{}
		var mediaPath string

		rows.Scan(&media.ID, &mediaPath, &media.PlayTime, &media.Date)

		media.File = utils.ProcessFile(mediaPath)
		mediaList = append(mediaList, media)
	}

	return mediaList
}

type FfmpegHistory struct {
	ArchivePath utils.File
	Mp4Path     utils.File
	codecs      string
	conversions string
	Duration    string
	Date        int
}

// FindFfmpegHistory gets the history for an ffmpeg conversion
func FindFfmpegHistory(anyPath string) FfmpegHistory {
	stmt, _ := con.Prepare(
		"SELECT `archivePath`, `mp4Path`, `codecs`, `conversions`, `duration`, `date`" +
			"FROM `ffmpeg` WHERE `archivePath` = ? OR `mp4Path` = ? LIMIT 1;")

	rows, _ := stmt.Query(anyPath, anyPath)
	var h FfmpegHistory

	for rows.Next() {
		var arcPath string
		var mp4Path string

		rows.Scan(&arcPath, &mp4Path, &h.codecs,
			&h.conversions, &h.Duration, &h.Date)

		h.ArchivePath = utils.ProcessFile(arcPath)
		h.Mp4Path = utils.ProcessFile(mp4Path)
	}

	return h
}

// SelectAllFfmpeg gets all ffmpeg histories
func SelectAllFfmpeg() []FfmpegHistory {
	var histories []FfmpegHistory

	stmt, _ := con.Prepare(
		"SELECT `archivePath`, `mp4Path`, `codecs`, `conversions`, `duration`, `date`" +
			"FROM `ffmpeg` ORDER BY `id` DESC;")

	rows, _ := stmt.Query()

	for rows.Next() {
		h := FfmpegHistory{}
		var arcPath string
		var mp4Path string

		rows.Scan(&arcPath, &mp4Path, &h.codecs,
			&h.conversions, &h.Duration, &h.Date)

		h.ArchivePath = utils.ProcessFile(arcPath)
		h.Mp4Path = utils.ProcessFile(mp4Path)

		histories = append(histories, h)
	}

	return histories
}

// SelectFfmpegPriority gets all ffmpeg proritiy folders
func SelectFfmpegPriority() []utils.File {
	var priorityFolders []utils.File

	stmt, _ := con.Prepare(
		"SELECT `path` FROM `ffmpegPriority` ORDER BY `id` DESC;")

	rows, _ := stmt.Query()

	for rows.Next() {
		var path string

		rows.Scan(&path)

		priorityFolders = append(priorityFolders, utils.ProcessFile(path))
	}

	return priorityFolders
}

package database

import (
	"fmt"

	"github.com/philcantcode/goApi/utils"
)

// Playback is the default struct for a database item
type Playback struct {
	ID        int
	Playtime  int
	Modified  int
	Locations []utils.File // All of the matching files on the disk
	PrefLoc   int          // The preferred [pos] for Locations
}

// SelectRootDirectories returns all the locations monitored on disk
func SelectRootDirectories() []utils.File {
	stmt, err := playerCon.Query("SELECT `path` FROM `RootDirectories`;")
	defer stmt.Close()
	utils.Error("Couldn't select from RootDirectories", err)

	var directories []utils.File

	for stmt.Next() {
		var path string

		stmt.Scan(&path)
		directories = append(directories, utils.ProcessFile(path))
	}

	return directories
}

// SelectPlaybackID_ByHash finds a single mediaPlaybackID which correlates to
// a playback in the MediaPlayback table
// Returns -1 for the ID if not found
func SelectPlaybackID_ByHash(hash string) (int, error) {
	stmt, err := playerCon.Prepare("SELECT `mediaPlaybackID` FROM `MediaHashes` WHERE `hash` = ? LIMIT 1;")
	defer stmt.Close()
	utils.Error("Couldn't SelectPlaybackID_ByHash from MediaHashes", err)

	rows, err := stmt.Query(hash)
	defer rows.Close()
	utils.Error("Results error from SelectPlaybackID_ByHash", err)

	for rows.Next() {
		var id int
		rows.Scan(&id)

		return id, nil
	}

	return -1, fmt.Errorf("Couldn't find entry in SelectPlaybackID_ByHash")
}

// SelectPlaybackID_ByPath finds the mediaPlaybackID by the path
// for quick lookup after the file has been hashed
// Returns -1 if not found
func SelectPlaybackID_ByPath(path string) (int, error) {
	stmt, err := playerCon.Prepare("SELECT `mediaPlaybackID` FROM `MediaHashes` WHERE `path` = ? LIMIT 1;")
	defer stmt.Close()
	utils.Error("Couldn't SelectPlaybackID_ByPath from MediaHashes", err)

	rows, err := stmt.Query(path)
	defer rows.Close()
	utils.Error("Results error from SelectPlaybackID_ByHash", err)

	for rows.Next() {
		var id int
		rows.Scan(&id)

		return id, nil
	}

	return -1, fmt.Errorf("Couldn't find entry in SelectPlaybackID_ByPath")
}

// SelectMediaPlayback_ByID returns the playback info given an ID
// The ID is guarenteed to exist at this point
func SelectMediaPlayback_ByID(id int) Playback {
	stmt, err := playerCon.Prepare("SELECT `playtime`, `modified` FROM `MediaPlayback` WHERE `id` = ?;")
	defer stmt.Close()
	utils.Error("Couldn't SelectMediaPlayback_ByID from MediaPlayback", err)

	rows, err := stmt.Query(id)
	defer rows.Close()
	utils.Error("Results error from SelectMediaPlayback_ByID", err)

	paths, err := selectMediaPaths_ByID(id)
	utils.Error("Couldn't retrieve path locations from selectMediaPaths_ByID", err)

	playback := Playback{ID: id, Locations: paths}

	for rows.Next() {
		rows.Scan(&playback.Playtime, &playback.Modified)
	}

	return playback
}

// selectMediaPaths_ByID returns an array of files that match the ID from
// the MediaHashes database
func selectMediaPaths_ByID(id int) ([]utils.File, error) {
	stmt, err := playerCon.Prepare("SELECT `path` FROM `MediaHashes` WHERE `mediaPlaybackID` = ?;")
	defer stmt.Close()
	utils.Error("Couldn't SelectMediaPath_ByID from MediaHashes", err)

	rows, err := stmt.Query(id)
	defer rows.Close()
	utils.Error("Results error from SelectMediaPath_ByID", err)

	paths := []utils.File{}

	for rows.Next() {
		var path string
		rows.Scan(&path)
		err = nil

		paths = append(paths, utils.ProcessFile(path))
	}

	if len(paths) == 0 {
		err = fmt.Errorf("Couldn't find any paths by ID in SelectMediaPaths_ByID")
	}

	return paths, err
}

// SelectPlaybacks_ByTime finds the playback in the from a given timepoint
func SelectPlaybacks_ByTime(unixTime int64) []Playback {
	var recentPlayList []Playback

	stmt, err := playerCon.Prepare("SELECT `id` FROM `MediaPlayback` WHERE `modified` >= ?;")
	defer stmt.Close()
	utils.Error("Couldn't SelectPlaybacks_ByTime from MediaPlayback", err)

	rows, err2 := stmt.Query(unixTime)
	defer rows.Close()
	utils.Error("Results error SelectPlaybacks_ByTime", err2)

	for rows.Next() {
		var id int
		rows.Scan(&id)

		recentPlayList = append(recentPlayList, SelectMediaPlayback_ByID(id))
	}

	return recentPlayList
}

// SelectAllPlaybacks finds the playback in the database
func SelectAllPlaybacks() []Playback {
	stmt, err := playerCon.Prepare("SELECT `id`, `playtime`, `modified` FROM `MediaPlayback`;")
	defer stmt.Close()
	utils.Error("Couldn't SelectAllPlaybacks from MediaPlayback", err)

	rows, err := stmt.Query()
	defer rows.Close()
	utils.Error("Results error from SelectAllPlaybacks", err)

	var playbackList []Playback

	for rows.Next() {
		playback := Playback{}
		rows.Scan(&playback.ID, &playback.Playtime, &playback.Modified)

		paths, err := selectMediaPaths_ByID(playback.ID)
		utils.Error("Couldn't retrieve path locations from SelectAllPlaybacks", err)
		playback.Locations = paths

		playbackList = append(playbackList, playback)
	}

	return playbackList
}

type FfmpegHistory struct {
	ID          int
	Path        utils.File // The new .mp4 path
	ArchivePath utils.File // The old .avi path in .ffmpeg
	OrigCodecs  string
	ConvCodecs  string
	Duration    string
	Date        int
}

// FindFfmpegHistory gets the history for an ffmpeg conversion
func FindFfmpegHistory(anyPath string) (FfmpegHistory, error) {
	stmt, err := playerCon.Prepare(
		"SELECT `id`, `path`, `archivePath`, `originalCodecs`, `convertedCodecs`, `duration`, `date`" +
			"FROM `FfmpegConversions` WHERE `path` = ? LIMIT 1;")

	defer stmt.Close()
	utils.Error("Couldn't select from FfmpegConversions", err)

	rows, _ := stmt.Query(anyPath)
	defer rows.Close()

	var h FfmpegHistory
	retError := fmt.Errorf("Conversion not found")

	for rows.Next() {
		var path string       // New .mp4
		var arcivePath string // Old .avi

		rows.Scan(&h.ID, &path, &arcivePath, &h.OrigCodecs,
			&h.ConvCodecs, &h.Duration, &h.Date)

		h.Path = utils.ProcessFile(path)
		h.ArchivePath = utils.ProcessFile(arcivePath)
		retError = nil
	}

	return h, retError
}

// SelectAllFfmpeg gets all ffmpeg histories
func SelectAllFfmpeg() []FfmpegHistory {
	var histories []FfmpegHistory

	stmt, err := playerCon.Prepare("SELECT `path`, `archivePath`, " +
		"`originalCodecs`, `convertedCodecs`, `duration`, `date` " +
		"FROM `FfmpegConversions` ORDER BY `id` DESC;")

	utils.Error("Couldn't select all from FfmpegConversions", err)
	defer stmt.Close()

	rows, err := stmt.Query()

	utils.Error("Results error from FfmpegConversions", err)
	defer rows.Close()

	for rows.Next() {
		h := FfmpegHistory{}
		var arcPath string
		var mp4Path string

		rows.Scan(&mp4Path, &arcPath, &h.OrigCodecs,
			&h.ConvCodecs, &h.Duration, &h.Date)

		h.Path = utils.ProcessFile(mp4Path)
		h.ArchivePath = utils.ProcessFile(arcPath)

		histories = append(histories, h)
	}

	return histories
}

// SelectFfmpegPriority gets all ffmpeg proritiy folders
func SelectFfmpegPriority() []utils.File {
	var priorityFolders []utils.File

	stmt, err := playerCon.Prepare("SELECT `path` FROM `FfmpegPriority` ORDER BY `ROWID` DESC;")
	utils.Error("Couldn't select from FfmpegPriority", err)
	defer stmt.Close()

	rows, _ := stmt.Query()
	defer rows.Close()

	for rows.Next() {
		var path string

		rows.Scan(&path)

		priorityFolders = append(priorityFolders, utils.ProcessFile(path))
	}

	return priorityFolders
}

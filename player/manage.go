package player

import (
	"fmt"
	"net/http"
	"os"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var managePage = page{
	PageTitle:       "Manage Local Files",
	PageDescription: "Manage Local Files",
	PreviousPath:    "Player",
	PreviousPathURL: "player",
	CurrentPath:     "player/manage",
}

// ManagePage handles the /player/manage request
func ManagePage(w http.ResponseWriter, r *http.Request) {
	reload()

	pathParam := r.FormValue("path")
	trackParam := r.FormValue("track")
	untrackParam := r.FormValue("untrack")

	data := struct {
		Selected       string
		Drives         []utils.File
		SubFolders     []utils.File
		TrackedFolders []utils.File
		FfmpegMetrics  []FfmpegMetrics
		FfmpegHistory  []database.FfmpegHistory
	}{
		Selected:       pathParam,
		TrackedFolders: database.SelectDirectories(),
		Drives:         utils.GetDefaultSystemDrives(),
		FfmpegMetrics:  FfmpegStat,
		FfmpegHistory:  database.SelectAllFfmpeg(),
	}

	if pathParam != "" {
		data.SubFolders = utils.GetFolderLayer(pathParam)
	}

	if trackParam != "" {
		database.InsertFolder(trackParam)
	}

	if untrackParam != "" {
		database.UnTrackFolder(untrackParam)
	}

	managePage.Contents = data
	err := templates.ExecuteTemplate(w, "playerManage", managePage)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

// RestoreFfmpeg reverts the file conversion by looking up the archive file
// in the database, moving the archive file back to the folder
// and deleting the original file
func RestoreFfmpeg(w http.ResponseWriter, r *http.Request) {
	// Original path, e.g., G:/folder/movie.avi
	origPath := r.FormValue("path")

	fmt.Println("[Restoring] " + origPath)

	f := utils.ProcessFile(origPath)
	archivedPath := database.FindFfmpegHistory(f.AbsPath)

	// Not found in database
	// Probably because not finished processing yet
	if archivedPath.ArchivePath.AbsPath == "" {
		w.Write([]byte("Not Found"))
		return
	}

	// Moves archived file back & delete mp4 file
	restorationPath := f.Path + archivedPath.ArchivePath.Name + archivedPath.ArchivePath.Ext
	fmt.Printf("Restoring (moving) %s to %s\n", archivedPath.ArchivePath.AbsPath, restorationPath)
	err := os.Rename(archivedPath.ArchivePath.AbsPath, restorationPath)

	if err != nil {
		fmt.Println("Could Not Restore File (Move Err)")
		w.Write([]byte("Move Error"))
		return
	}

	// Remove the old .mp4 file
	err = os.Remove(f.Path + f.Name + ".mp4")

	if err != nil {
		w.Write([]byte("Delete Error"))
		return
	}

	// Remove entry from the DB
	database.DeleteFfmpegEntry(archivedPath.ArchivePath.AbsPath)
}

func PlayFfmpeg(w http.ResponseWriter, r *http.Request) {
	pathParam := r.FormValue("path")

	fmt.Println("[Playing] " + pathParam)
}

// ControlFfmpeg handles web requests that stop the file
// conversion or puts it into fast / slow mode (slow mode
// is single threaded )
func ControlFfmpeg(w http.ResponseWriter, r *http.Request) {
	controlType := r.FormValue("type")
	prioritise := r.FormValue("prioritise")

	switch controlType {
	case "fast":
		DisableFfmpeg = false
		NumFfmpegThreads = 0
		fmt.Println("Switching FFMPEG to fast convert")
	case "slow":
		DisableFfmpeg = false
		NumFfmpegThreads = 1
		fmt.Println("Switching FFMPEG to slow convert")
	case "disable":
		DisableFfmpeg = true
		NumFfmpegThreads = 0
		fmt.Println("Disabling FFMPEG conversion")
	case "prioritise":
		fmt.Println("Prioritising " + prioritise)
		database.InsertFfmpegPriority(prioritise)
	}
}

package player

import (
	"fmt"
	"net/http"
	"os"

	"github.com/philcantcode/goApi/index"
	"github.com/philcantcode/goApi/player/database"
	"github.com/philcantcode/goApi/utils"
)

var pageContents = index.HTMLContents{
	PageTitle:       "Manage Local Files",
	PageDescription: "Manage Local Player Files",
}

// ManagePage handles the /player/manage request
func ManagePage(w http.ResponseWriter, r *http.Request) {
	index.Reload()

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
		TrackedFolders: database.SelectRootDirectories(),
		Drives:         utils.GetDefaultSystemDrives(),
		FfmpegMetrics:  FfmpegStats,
		FfmpegHistory:  database.SelectAllFfmpeg(),
	}

	if pathParam != "" {
		data.SubFolders = utils.GetFolderLayer(pathParam)
	}

	if trackParam != "" {
		database.InsertRootDirectory(trackParam)
	}

	if untrackParam != "" {
		database.DeleteRootDirectory(untrackParam)
	}

	pageContents.Contents = data
	err := index.TemplateLoader.ExecuteTemplate(w, "playerManage", pageContents)

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
	mp4Path := r.FormValue("path")
	fmt.Println("[Restoring] " + mp4Path)

	history, err := database.FindFfmpegHistory(utils.ProcessFile(mp4Path).AbsPath)

	// Not found in database
	// Probably because not finished processing yet
	if err != nil {
		w.Write([]byte("Not Found"))
		return
	}

	if !history.ArchivePath.Exists {
		w.Write([]byte("Archive File Not Found"))
		return
	}

	// Moves archived file back & delete mp4 file
	srcPath := history.ArchivePath.AbsPath
	dstPath := history.Path.Path + history.Path.FileName + history.ArchivePath.Ext
	fmt.Printf("Restoring (moving) %s to %s\n", srcPath, dstPath)
	err = os.Rename(srcPath, dstPath)

	if err != nil {
		fmt.Println("Could Not Restore File (Move Err)")
		w.Write([]byte("Move Error"))
		return
	}

	// Remove archive and old file
	err = os.Remove(mp4Path)

	if err != nil {
		w.Write([]byte("Delete Error"))
		return
	}

	// Remove entry from the DB
	database.DeleteFfmpegEntry(history.ID)
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

		if prioritise == "" {
			break
		}

		database.InsertFfmpegPriority(prioritise)
	}
}

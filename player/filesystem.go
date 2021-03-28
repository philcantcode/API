package player

import (
	"fmt"
	"net/http"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

// FileTrack handles the /os request
func FileTrack(w http.ResponseWriter, r *http.Request) {
	reload()

	pathParam := r.FormValue("path")
	trackParam := r.FormValue("track")
	untrackParam := r.FormValue("untrack")

	data := struct {
		Selected       string
		Drives         []string
		SubFolders     []string
		TrackedFolders []database.Directory
		FfmpegStat     []ConversionHistory
	}{
		Selected:       pathParam,
		TrackedFolders: database.SelectDirectories(),
		Drives:         utils.GetDrives(),
		FfmpegStat:     FfmpegStat,
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

	localTrackPage.Contents = data
	err := templates.ExecuteTemplate(w, "localTracker", localTrackPage)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

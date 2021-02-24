package server

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
		TrackedFolders []database.TrackFolders
	}{
		Selected:       pathParam,
		TrackedFolders: database.GetTrackedFolders(),
		Drives:         utils.GetDrives(),
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

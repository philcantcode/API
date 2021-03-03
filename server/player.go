package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

// PlayerPage handles the /player request
func PlayerPage(w http.ResponseWriter, r *http.Request) {
	reload()

	openParam := r.FormValue("open")
	playParam := r.FormValue("play")

	data := struct {
		IP   string
		Port string

		OpenParam string

		// File selector menu containing Directories > Folders > Files
		Directories []string
		SubFolders  []string
		Files       []utils.File

		// Media that is being played
		MediaInfo     database.MediaInfo
		NextMediaItem int
	}{
		IP:        utils.Host,
		Port:      utils.Port,
		MediaInfo: database.FindOrCreateMedia(playParam),
		OpenParam: openParam,
	}

	// Find top level directories
	for _, s := range database.GetDirectories() {
		data.Directories = append(data.Directories, s.Path)
	}

	if openParam != "" {
		// Find sub folders
		for _, s := range utils.GetFolderLayer(openParam) {
			data.SubFolders = append(data.SubFolders, s)
		}

		// Find files in folder
		for _, s := range utils.GetFilesLayer(openParam) {
			data.Files = append(data.Files, s)
		}
	}

	// Play the media in the playParam
	if playParam != "" {

		// Find the next media ID
		nextMedia := utils.GetNextMatchingOrderedFile(data.MediaInfo.Folder, data.MediaInfo.Path)
		data.NextMediaItem = database.FindOrCreateMedia(nextMedia).ID
	}

	playerPage.Contents = data
	templates.ExecuteTemplate(w, "player", playerPage)
}

// LoadMedia takes a file or ID GET param, then loads the media
func LoadMedia(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("file")
	id, _ := strconv.Atoi(r.FormValue("id"))
	var mediaInfo database.MediaInfo

	if file != "" {
		mediaInfo = database.SelectMedia(file)
		http.ServeFile(w, r, file)
	}

	if id != 0 {
		mediaInfo = database.SelectMediaByID(id)
		http.ServeFile(w, r, mediaInfo.Path)
	}

	fmt.Printf("Media Loaded (%d) %s\n", mediaInfo.ID, mediaInfo.Title)
	OpenChannel(mediaInfo)
}

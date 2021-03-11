package player

import (
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
		MediaInfo database.MediaInfo
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

	//OpenChannel(mediaInfo)
}

// When a playback update comes in
func playbackUpdate(playTimeStr string, mediaID int) {
	playTime, _ := strconv.ParseFloat(playTimeStr, 64)
	database.UpdatePlaytime(mediaID, int(playTime))
}

func findNextMedia(mediaID int) int {
	prevMedia := database.SelectMediaByID(mediaID)
	nextMedia := utils.GetNextMatchingOrderedFile(prevMedia.Folder, prevMedia.Path)
	nextID := database.FindOrCreateMedia(nextMedia).ID

	return nextID
}

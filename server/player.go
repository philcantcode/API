package server

import (
	"fmt"
	"net/http"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

// PlayerControls handles the media player controls
func PlayerControls(w http.ResponseWriter, r *http.Request) {
	reload()

	play := r.FormValue("play")
	pause := r.FormValue("pause")
	forward := r.FormValue("forward")
	back := r.FormValue("back")
	skip := r.FormValue("skip")
	update := r.FormValue("update")

	if play != "" {

	}

	if pause != "" {

	}

	if forward != "" {

	}

	if back != "" {

	}

	if skip != "" {

	}

	if update != "" {

	}

}

// Player handles the /player request
func Player(w http.ResponseWriter, r *http.Request) {
	reload()

	openParam := r.FormValue("open")
	playParam := r.FormValue("play")

	media := database.FindOrCreateMedia(playParam)

	//https://blog.addpipe.com/10-advanced-features-in-html5-video-player/
	data := struct {
		IP   string
		Port string

		Folders    []string
		SubFolders []string
		Files      []utils.File

		Media database.Media
	}{
		IP:    utils.Host,
		Port:  utils.Port,
		Media: media,
	}

	for _, s := range database.GetTrackedFolders() {
		data.Folders = append(data.Folders, s.Path)
	}

	if openParam != "" {
		for _, s := range utils.GetFolderLayer(openParam) {
			data.SubFolders = append(data.SubFolders, s)
		}

		for _, s := range utils.GetFilesLayer(openParam) {
			data.Files = append(data.Files, s)
		}
	}

	playerPage.Contents = data
	err := templates.ExecuteTemplate(w, "player", playerPage)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

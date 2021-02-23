package api

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var templates *template.Template

func init() {
	reload()
}

func reload() { // When done, remove calls to reload
	var err error
	templates, err = template.ParseFiles(utils.FilePath+"index.html", utils.FilePath+"player.html",
		utils.FilePath+"os.html", utils.FilePath+"header.html", utils.FilePath+"footer.html")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

// IndexHandler handles the / (Root) request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", indexPage)
}

// MediaHandler to handle file requests
func MediaHandler(w http.ResponseWriter, r *http.Request) {
	loadParam := r.FormValue("load")
	http.ServeFile(w, r, loadParam)
}

// PlayerHandler handles the /player request
func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	reload()

	openParam := r.FormValue("open")
	playParam := r.FormValue("play")

	//https://blog.addpipe.com/10-advanced-features-in-html5-video-player/
	data := struct {
		Folders    []string
		SubFolders []string
		Files      []utils.File

		Media         string
		MediaTitle    string
		MediaPlayback string
	}{
		Media:         playParam,
		MediaTitle:    utils.ExtractFileName(playParam),
		MediaPlayback: "0",
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

// LocalTrackHandler handles the /os request
func LocalTrackHandler(w http.ResponseWriter, r *http.Request) {
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
		database.TrackFolder(trackParam)
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

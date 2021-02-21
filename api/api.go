package api

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var indexTemplate *template.Template
var playerTemplate *template.Template
var locTrackTemplate *template.Template

func init() {
	reload()
}

func reload() { // When done, remove calls to reload
	var err error
	indexTemplate, err = template.ParseFiles(utils.FilePath+"index.html", utils.FilePath+"header.html", utils.FilePath+"footer.html")
	playerTemplate, err = template.ParseFiles(utils.FilePath+"player.html", utils.FilePath+"header.html", utils.FilePath+"footer.html")
	locTrackTemplate, err = template.ParseFiles(utils.FilePath+"os.html", utils.FilePath+"header.html", utils.FilePath+"footer.html")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

// IndexHandler handles the / (Root) request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.ExecuteTemplate(w, "index", indexPage)
}

// PlayerHandler handles the /player request
func PlayerHandler(w http.ResponseWriter, r *http.Request) {
	reload()
	playerTemplate.ExecuteTemplate(w, "player", playerPage)
}

// LocalTrackHandler handles the /os request
func LocalTrackHandler(w http.ResponseWriter, r *http.Request) {
	reload()

	data := struct {
		Selected       string
		Drives         []string
		SubFolders     []string
		TrackedFolders []database.TrackFolders
	}{}

	pathParam := r.FormValue("path")
	trackParam := r.FormValue("track")
	untrackParam := r.FormValue("untrack")

	data.Selected = pathParam

	if pathParam != "" {
		data.SubFolders = utils.GetFolderLayer(pathParam)
	}

	if trackParam != "" {
		database.TrackFolder(trackParam)
	}

	if untrackParam != "" {
		database.UnTrackFolder(untrackParam)
	}

	data.TrackedFolders = database.GetTrackedFolders()
	data.Drives = utils.GetDrives()

	localTrackPage.Contents = data
	err := locTrackTemplate.ExecuteTemplate(w, "localTracker", localTrackPage)

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

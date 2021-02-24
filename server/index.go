package server

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/philcantcode/goApi/utils"
)

type page struct {
	PageTitle       string
	PageDescription string
	PreviousPath    string
	PreviousPathURL string
	CurrentPath     string

	Contents interface{}
}

var (
	indexPage = page{
		PageTitle:       "HomePage",
		PageDescription: "OwO Player Homepage",
		PreviousPath:    "Home",
		PreviousPathURL: "/",
		CurrentPath:     "Home",
	}

	playerPage = page{
		PageTitle:       "Player",
		PageDescription: "Local Player",
		PreviousPath:    "Home",
		PreviousPathURL: "/",
		CurrentPath:     "Player",
	}

	localTrackPage = page{
		PageTitle:       "Local Files",
		PageDescription: "Track Local Files",
		PreviousPath:    "Home",
		PreviousPathURL: "/",
		CurrentPath:     "OS",
	}
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

// Index handles the / (Root) request
func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", indexPage)
}

// Media to handle file requests
func Media(w http.ResponseWriter, r *http.Request) {
	loadParam := r.FormValue("load")
	http.ServeFile(w, r, loadParam)
}

package player

import (
	"fmt"
	"math/rand"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
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

	remotePage = page{
		PageTitle:       "Remote Control",
		PageDescription: "Control the playback on other screens",
		PreviousPath:    "Player",
		PreviousPathURL: "/player",
		CurrentPath:     "Remote",
	}
)

var templates *template.Template

var store = sessions.NewCookieStore([]byte("temp"))

func init() {
	token := make([]byte, 32)
	rand.Read(token)
	store = sessions.NewCookieStore(token)

	reload()
}

func reload() { // When done, remove calls to reload
	var err error
	templates, err = template.ParseFiles(
		utils.FilePath+"index.html", utils.FilePath+"player.html",
		utils.FilePath+"os.html", utils.FilePath+"header.html",
		utils.FilePath+"footer.html", utils.FilePath+"remote.html")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

// IndexPage handles the / (Root) request
func IndexPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", indexPage)
}

package player

import (
	"fmt"
	"net/http"
	"text/template"
)

type page struct {
	PageTitle       string
	PageDescription string
	PreviousPath    string
	PreviousPathURL string
	CurrentPath     string

	Contents interface{}
}

var indexPage = page{
	PageTitle:       "HomePage",
	PageDescription: "OwO Player Homepage",
	PreviousPath:    "Home",
	PreviousPathURL: "/",
	CurrentPath:     "Home",
}

var templates *template.Template

func init() {
	reload()
}

func reload() { // When done, remove calls to reload
	var err error

	// Parse all .gohtml template files
	templates, err = template.ParseGlob("web/*.gohtml")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

// IndexPage handles the / (Root) request
func IndexPage(w http.ResponseWriter, r *http.Request) {
	// Loads the template, 2nd param is the name in the .gohtml file at top
	templates.ExecuteTemplate(w, "index", indexPage)
}

// FavIcon sends back icon image
func FavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/web/res/icons/favicon.ico")
}

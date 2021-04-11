package notes

import (
	"net/http"

	"github.com/philcantcode/goApi/index"
)

var notesPage = index.HTMLContents{
	PageTitle:       "Notes",
	PageDescription: "Wiki Notes",
}

// NotesPage handles the /notes request
func NotesPage(w http.ResponseWriter, r *http.Request) {

	index.TemplateLoader.ExecuteTemplate(w, "notes", notesPage)
}

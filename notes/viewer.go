package notes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/philcantcode/goApi/index"
	notes "github.com/philcantcode/goApi/notes/database"
	"github.com/philcantcode/goApi/utils"
)

var notesViewer = index.HTMLContents{
	PageTitle:       "WikiNotes",
	PageDescription: "Wiki Notes",
}

// NotesPage handles the /notes request
func ViewerPage(w http.ResponseWriter, r *http.Request) {
	index.Reload()

	urlParams := mux.Vars(r)
	noteID, err := strconv.Atoi(urlParams["id"])

	keywords, err := json.Marshal(notes.SelectKeywords())
	utils.Error("Couldn't marshall keywords", err)

	data := struct {
		RecentNotes []notes.NoteContents
		Note        notes.NoteContents
		NoteJson    string
		Keywords    string
	}{
		RecentNotes: notes.SelectRecentNotes(3),
		Keywords:    string(keywords),
	}

	if urlParams["id"] == "new" {
		data.Note = notes.NoteContents{}
		notesViewer.PageDescription = "Create New Note"
	}

	if err == nil {
		data.Note = notes.SelectNote(noteID)
		notesViewer.PageDescription = data.Note.Keyword
	}

	dat, _ := json.Marshal(data.Note)
	data.NoteJson = string(dat)

	notesViewer.Contents = data
	index.TemplateLoader.ExecuteTemplate(w, "viewer", notesViewer)
}

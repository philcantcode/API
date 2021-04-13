package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/philcantcode/goApi/index"
	notes "github.com/philcantcode/goApi/notes/database"
	"github.com/philcantcode/goApi/utils"
)

var notesPage = index.HTMLContents{
	PageTitle:       "WikiNotes",
	PageDescription: "Wiki Notes",
}

// NotesPage handles the /notes request
func NotesPage(w http.ResponseWriter, r *http.Request) {
	index.Reload()

	data := struct {
		RecentNotes []notes.NoteContents
	}{
		RecentNotes: notes.SelectRecentNotes(50),
	}

	notesPage.Contents = data
	index.TemplateLoader.ExecuteTemplate(w, "notes", notesPage)
}

type Response struct {
	Type    string
	Message string
	Value   string
}

// CreateNote handles the POST request with new notes
func CreateNote(w http.ResponseWriter, r *http.Request) {
	contents := r.FormValue("contents")

	note := notes.NoteContents{}
	err := json.Unmarshal([]byte(contents), &note)
	utils.Error("CreateNote JSON error", err)

	note.Keyword = strings.ToLower(note.Keyword)

	if len(note.Keyword) == 0 {
		response := jsonResponse(
			Response{
				Type:    "Error",
				Message: "Keyword Required",
			})

		w.Write([]byte(response))
		return
	}

	// Detect duplicate keywords
	for _, v := range notes.SelectKeywords() {
		if v.Keyword == note.Keyword {
			response := jsonResponse(
				Response{
					Type:    "Error",
					Message: "Keyword Exists Already",
				})

			w.Write([]byte(response))
			return
		}
	}

	id, err := notes.InsertNote(note.Keyword, note.Desc, contents)

	if err == nil {
		response := jsonResponse(
			Response{
				Type:    "Success",
				Message: "Note Inserted",
				Value:   fmt.Sprintf("%d", id),
			})

		w.Write([]byte(response))
	}
}

// CreateNote handles the POST request with new notes
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	contents := r.FormValue("contents")

	note := notes.NoteContents{}
	err := json.Unmarshal([]byte(contents), &note)
	utils.Error("CreateNote JSON error", err)
	note.Keyword = strings.ToLower(note.Keyword)

	notes.UpdateNote(note.ID, note.Keyword, note.Desc, contents)

	response := jsonResponse(
		Response{
			Type:    "Success",
			Message: "Note Updated Successfully",
		})

	w.Write([]byte(response))
}

// DeleteNote handles the POST delete request with new notes
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	utils.Error("Couldn't delete note", err)

	notes.DeleteNote(id)
	fmt.Printf("Note Deleted\n")

	response := jsonResponse(
		Response{
			Type:    "Success",
			Message: "Note Successfulyl Deleted",
		})

	w.Write([]byte(response))
}

func jsonResponse(r Response) string {
	jsonStruct, err := json.Marshal(r)
	utils.Error("Couldn't convert response to Json", err)

	return string(jsonStruct)
}

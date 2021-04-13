package notes

import (
	"encoding/json"

	"github.com/philcantcode/goApi/utils"
)

type Element struct {
	Key   string
	Value string
}

type NoteContents struct {
	ID       int
	Modified int

	Keyword  string
	Desc     string
	Elements []Element
}

type Keyword struct {
	ID      int
	Keyword string
}

// SelectKeywords returns all keywords
func SelectKeywords() []Keyword {
	stmt, err := wikiCon.Prepare("SELECT `id`, `keyword` FROM `Notes`;")

	utils.Error("Couldn't select from SelectKeywords", err)
	defer stmt.Close()

	rows, err := stmt.Query()
	utils.Error("Results error from SelectKeywords", err)
	defer rows.Close()

	var keywords []Keyword

	for rows.Next() {
		var keyword Keyword

		rows.Scan(&keyword.ID, &keyword.Keyword)
		keywords = append(keywords, keyword)
	}

	return keywords
}

// SelectRecentNotes returns all recent notes
func SelectRecentNotes(limit int) []NoteContents {
	stmt, err := wikiCon.Prepare("SELECT `id`, `keyword`, `desc`, `text`, `modified` FROM `Notes` ORDER BY `id` DESC LIMIT ?;")

	utils.Error("Couldn't select from SelectRecentNotes", err)
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	utils.Error("Results error from SelectRecentNotes", err)
	defer rows.Close()

	var notesList []NoteContents

	for rows.Next() {
		note := NoteContents{}
		var text string
		var id int

		rows.Scan(&id, &note.Keyword, &note.Desc, &text, &note.Modified)
		json.Unmarshal([]byte(text), &note)
		note.ID = id
		notesList = append(notesList, note)
	}

	return notesList
}

// SelectNote returns 1 notes
func SelectNote(id int) NoteContents {
	stmt, err := wikiCon.Prepare("SELECT `id`, `keyword`, `desc`, `text`, `modified` FROM `Notes` WHERE `id` = ?;")

	utils.Error("Couldn't select from SelectNote", err)
	defer stmt.Close()

	rows, err := stmt.Query(id)
	utils.Error("Results error from SelectNote", err)
	defer rows.Close()

	note := NoteContents{}

	for rows.Next() {
		var text string
		var id int

		rows.Scan(&id, &note.Keyword, &note.Desc, &text, &note.Modified)
		json.Unmarshal([]byte(text), &note)
		note.ID = id
	}

	return note
}

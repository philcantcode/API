package index

import (
	"fmt"
	"text/template"
)

type HTMLContents struct {
	PageTitle       string
	PageDescription string

	Contents interface{}
}

var TemplateLoader *template.Template

func Reload() { // When done, remove calls to reload
	var err error

	// Parse all .gohtml template files
	TemplateLoader, err = template.ParseGlob("web/*/*.gohtml")

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

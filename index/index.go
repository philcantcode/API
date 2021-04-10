package index

import (
	"net/http"

	"github.com/philcantcode/goApi/utils"
)

func init() {
	Reload()
}

// IndexPage handles the / (Root) request
func IndexPage(w http.ResponseWriter, r *http.Request) {
	// Loads the template, 2nd param is the name in the .gohtml file at top
	err := TemplateLoader.ExecuteTemplate(w, "index", HTMLContents{PageTitle: "HomePage", PageDescription: "OwO Homepage"})
	utils.Error("Index template load error", err)
}

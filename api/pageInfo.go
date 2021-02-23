package api

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

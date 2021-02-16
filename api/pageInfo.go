package api

type page struct {
	PageTitle       string
	PageDescription string
	PreviousPath    string
	PreviousPathURL string
	CurrentPath     string

	Year string
}

var indexPage = page{Year: "2021", PageTitle: "HomePage", PageDescription: "OwO Player Homepage", PreviousPath: "Home", PreviousPathURL: "/", CurrentPath: "Home"}
var playerPage = page{Year: "2021", PageTitle: "Player", PageDescription: "Local Player", PreviousPath: "Home", PreviousPathURL: "/", CurrentPath: "Player"}
var localTrackPage = page{Year: "2021", PageTitle: "Local Files", PageDescription: "Track Local Files", PreviousPath: "Home", PreviousPathURL: "/", CurrentPath: "OS"}

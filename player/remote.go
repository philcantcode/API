package player

import (
	"net/http"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var remotePage = page{
	PageTitle:       "Remote Control",
	PageDescription: "Control the playback on other screens",
	PreviousPath:    "Player",
	PreviousPathURL: "/player",
	CurrentPath:     "Remote",
}

type LoadedMedia struct {
	RemoteID string
	Playback database.Playback
}

// RemotePage handles the remote controller
func RemotePage(w http.ResponseWriter, r *http.Request) {
	reload()

	idParam := r.FormValue("id")

	data := struct {
		IP   string
		Port string

		// List of media info with active channels
		LoadedMedia      []LoadedMedia
		ControllingMedia LoadedMedia
		RemoteID         string
	}{
		IP:       utils.Host,
		Port:     utils.Port,
		RemoteID: idParam,
	}

	for remoteID, channel := range players {
		data.LoadedMedia = append(data.LoadedMedia, LoadedMedia{RemoteID: remoteID, Playback: channel.playback})
	}

	if idParam != "" {
		data.ControllingMedia = LoadedMedia{RemoteID: idParam, Playback: players[idParam].playback}
	}

	remotePage.Contents = data
	templates.ExecuteTemplate(w, "playerRemote", remotePage)
}

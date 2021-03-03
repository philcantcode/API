package server

import (
	"net/http"
	"strconv"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

// RemotePage handles the remote controller
func RemotePage(w http.ResponseWriter, r *http.Request) {
	reload()

	remoteID, _ := strconv.Atoi(r.FormValue("controller"))

	data := struct {
		IP   string
		Port string

		// List of media info with active channels
		OpenMediaInfoList []database.MediaInfo

		// If a remote is selected, media is selected
		RemoteID        int
		RemoteMediaInfo database.MediaInfo
		NextMediaItem   int
	}{
		IP:              utils.Host,
		Port:            utils.Port,
		RemoteID:        remoteID,
		RemoteMediaInfo: database.SelectMediaByID(remoteID),
	}

	nextMedia := utils.GetNextMatchingOrderedFile(data.RemoteMediaInfo.Folder, data.RemoteMediaInfo.Path)
	data.NextMediaItem = database.FindOrCreateMedia(nextMedia).ID

	for _, channel := range channels {
		data.OpenMediaInfoList = append(data.OpenMediaInfoList, channel.mediaInfo)
	}

	remotePage.Contents = data
	templates.ExecuteTemplate(w, "remote", remotePage)
}

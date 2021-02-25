package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var loadedMedia []database.Media

// PlayerRemote handles the remote controller
func PlayerRemote(w http.ResponseWriter, r *http.Request) {
	reload()

	controller := r.FormValue("controller")
	controllerInt, _ := strconv.ParseInt(controller, 10, 64)

	data := struct {
		IP              string
		Port            string
		Loaded          []database.Media
		Controller      string
		ControllerMedia database.Media
	}{
		IP:              utils.Host,
		Port:            utils.Port,
		Loaded:          loadedMedia,
		Controller:      controller,
		ControllerMedia: database.SelectMediaByID(controllerInt),
	}

	remotePage.Contents = data
	templates.ExecuteTemplate(w, "remote", remotePage)
}

// PlayerControlsSocket fuq off
func PlayerControlsSocket(ws *websocket.Conn) {
	log.Printf("Socket thing")
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// PlayerControls handles the media player controls
func PlayerControls(w http.ResponseWriter, r *http.Request) {
	reload()

	/*
		if r.Method == http.MethodGet {
			load := r.FormValue("load")

			if load != "" {
				http.ServeFile(w, r, load)
			}

			return
		}*/

	ws, err := upgrader.Upgrade(w, r, nil)
	utils.Err("Couldn't upgrade", err)
	// register client
	go PlayerControlsSocket(ws)

}

// PlayerStatus handles post pushes from the client > server
func PlayerStatus(w http.ResponseWriter, r *http.Request) {
	reload()

	id := r.FormValue("id")
	playtime := r.FormValue("playtime")

	if id == "" {
		os.Exit(0)
	}

	if playtime != "" {
		int_id, _ := strconv.ParseInt(id, 10, 64)
		int_pt, _ := strconv.ParseFloat(playtime, 64)

		database.UpdatePlaytime(int(int_id), int(int_pt))
	}

	w.Write([]byte("aaaaaaaaaaa"))

}

// Player handles the /player request
func Player(w http.ResponseWriter, r *http.Request) {
	reload()

	openParam := r.FormValue("open")
	playParam := r.FormValue("play")

	// Keep track of what mediais loaded on the client side
	if playParam != "" {
		media := database.SelectMedia(playParam)
		alreadyLoaded := false

		for i := 0; i < len(loadedMedia); i++ {
			if loadedMedia[i].ID == media.ID {
				alreadyLoaded = true
			}
		}

		if !alreadyLoaded {
			loadedMedia = append(loadedMedia, media)

			fmt.Printf("Loaded Media %d, currently %d loaded\n",
				loadedMedia[len(loadedMedia)-1].ID, len(loadedMedia))
		}
	}

	media := database.FindOrCreateMedia(playParam)

	data := struct {
		IP   string
		Port string

		Folders    []string
		SubFolders []string
		Files      []utils.File

		Media database.Media
	}{
		IP:    utils.Host,
		Port:  utils.Port,
		Media: media,
	}

	for _, s := range database.GetTrackedFolders() {
		data.Folders = append(data.Folders, s.Path)
	}

	if openParam != "" {
		for _, s := range utils.GetFolderLayer(openParam) {
			data.SubFolders = append(data.SubFolders, s)
		}

		for _, s := range utils.GetFilesLayer(openParam) {
			data.Files = append(data.Files, s)
		}
	}

	playerPage.Contents = data
	err := templates.ExecuteTemplate(w, "player", playerPage)

	utils.Err("Player Handler", err)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

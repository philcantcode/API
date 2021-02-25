package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/goApi/server"
	"github.com/philcantcode/goApi/utils"
)

// Handle incoming web requests and direct them to the folder
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", server.Index)
	router.HandleFunc("/player", server.Player)
	router.HandleFunc("/player/controls", server.PlayerControls)
	router.HandleFunc("/player/status", server.PlayerStatus)
	router.HandleFunc("/player/remote", server.PlayerRemote)
	router.HandleFunc("/os", server.FileTrack)

	fileServer := http.FileServer(http.Dir(utils.FilePath))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	err := http.ListenAndServe(utils.Host+":"+utils.Port, router)

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
	}
}

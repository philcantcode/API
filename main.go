package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/goApi/player"
	"github.com/philcantcode/goApi/utils"
)

// Handle incoming web requests and direct them to the folder
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", player.IndexPage)
	router.HandleFunc("/player", player.PlayerPage)
	router.HandleFunc("/player/remote", player.RemotePage)

	router.HandleFunc("/player/ws-setup", player.SocketSetup)
	router.HandleFunc("/player/load", player.LoadMedia)

	router.HandleFunc("/os", player.FileTrack)

	fileServer := http.FileServer(http.Dir(utils.FilePath))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	err := http.ListenAndServe(utils.Host+":"+utils.Port, router)

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
	}
}

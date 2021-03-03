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

	router.HandleFunc("/", server.IndexPage)
	router.HandleFunc("/player", server.PlayerPage)
	router.HandleFunc("/player/remote", server.RemotePage)

	router.HandleFunc("/player/ws-setup", server.WebSocketSetup)
	router.HandleFunc("/player/load", server.LoadMedia)

	router.HandleFunc("/os", server.FileTrack)

	fileServer := http.FileServer(http.Dir(utils.FilePath))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	err := http.ListenAndServe(utils.Host+":"+utils.Port, router)

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
	}
}

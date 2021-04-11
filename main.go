package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/goApi/index"
	"github.com/philcantcode/goApi/notes"
	"github.com/philcantcode/goApi/player"
	"github.com/philcantcode/goApi/utils"
)

// Handle incoming web requests and direct them to the folder
func main() {
	flag.IntVar(&player.NumFfmpegThreads, "ffthreads", 1, "Number of threads for FFMPEG Conversions")
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/", index.IndexPage)

	router.HandleFunc("/player", player.PlayerPage)
	router.HandleFunc("/player/remote", player.RemotePage)
	router.HandleFunc("/player/manage", player.ManagePage)
	router.HandleFunc("/player/ffmpeg/revert", player.RestoreFfmpeg)
	router.HandleFunc("/player/ffmpeg/play", player.PlayFfmpeg)
	router.HandleFunc("/player/ffmpeg/control", player.ControlFfmpeg)
	router.HandleFunc("/player/ws-setup/{pageType}/{devID}", player.SocketSetup)
	router.HandleFunc("/player/load", player.LoadMedia)

	router.HandleFunc("/notes", notes.NotesPage)

	fileServer := http.FileServer(http.Dir(utils.FilePath))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	err := http.ListenAndServe(":"+utils.Port, router) //

	if err != nil {
		log.Fatal("Error Starting the HTTP Server :", err)
	}
}

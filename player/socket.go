package player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var players = make(map[string]player)

type player struct {
	playback database.Playback
	players  []*websocket.Conn
	remotes  []*websocket.Conn

	playCH   chan command
	remoteCH chan command
}

type command struct {
	Type  string
	Key   string
	Value string
}

type Response struct {
	Type     string
	Key      string
	Value    string
	Playback database.Playback
}

func init() {
	go processor()
}

// SocketSetup creates web sockets, optional id parameter
func SocketSetup(w http.ResponseWriter, r *http.Request) {
	// Media ID as GET param
	urlParams := mux.Vars(r)
	pageType := urlParams["pageType"]
	devID := urlParams["devID"]

	// uncomment for websocket object
	ws, err := upgrader.Upgrade(w, r, nil)
	utils.Error("Couldn't upgrade websocket in SocketSetup", err)

	_, exists := players[devID]
	var ch player

	if !exists {
		ch = player{
			playCH:   make(chan command, 10),
			remoteCH: make(chan command, 10),
		}

		if pageType == "remote" {
			ch.remotes = append(ch.remotes, ws)
			go remoteSocket(ws, devID)
		}

		if pageType == "player" {
			ch.players = append(ch.players, ws)
			go playerSocket(ws, devID)
		}

		players[devID] = ch
		fmt.Printf("SocketSetup (new): type: %s, device ID: %s\n", pageType, devID)
	} else {
		var ch player = players[devID]

		if pageType == "remote" {
			ch.remotes = append(players[devID].remotes, ws)
			go remoteSocket(ws, devID)
		}

		if pageType == "player" {
			ch.players = append(players[devID].players, ws)
			go playerSocket(ws, devID)
		}

		players[devID] = ch
		fmt.Printf("SocketSetup (existing): type: %s, device ID: %s\n", pageType, devID)
	}
}

func processor() {
	for range time.Tick(300 * time.Millisecond) {
		for id, val := range players {
			// Remote channels
			for i := 0; i < len(val.remoteCH); i++ {
				cmd := <-val.remoteCH

				switch cmd.Type {
				case "control":
					controls(cmd, id)
				default:
					fmt.Println("Socket Processor doesn't recognise remote command")
				}
			}

			// Player channels
			for i := 0; i < len(val.playCH); i++ {
				cmd := <-val.playCH

				switch cmd.Type {
				case "status":
					status(cmd, id)
				case "control":
					controls(cmd, id)
				default:
					utils.ErrorC("Socket Processor doesn't recognise player command")
				}
			}
		}
	}
}

// The messages come in from the player here
func playerSocket(ws *websocket.Conn, devID string) {
	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			if err.Error() == "websocket: close 1001 (going away)" {
				fmt.Printf("Device (%s) closed websocket\n", devID)
				delete(players, devID)
				return
			}

			continue
		}

		cmd := command{}
		err = json.Unmarshal(msg, &cmd)
		utils.Error("Error opening JSON", err)

		players[devID].playCH <- cmd
		fmt.Printf("Device %s (player) -> %+v\n", devID, cmd)
	}
}

// The messages come in from the remote here
func remoteSocket(ws *websocket.Conn, devID string) {
	for {
		_, msg, err := ws.ReadMessage()

		// Fix error where if you refresh remote, it kills the whole socket
		if err != nil {
			if err.Error() == "websocket: close 1001 (going away)" {
				fmt.Printf("Device (%s) disconnected from websocket\n", devID)
				ws.Close()
				return
			}

			continue
		}

		cmd := command{}
		err = json.Unmarshal(msg, &cmd)
		utils.Error("Error opening JSON", err)

		players[devID].playCH <- cmd
		fmt.Printf("Device %s (remote) -> %+v\n", devID, cmd)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func controls(cmd command, devID string) {
	switch cmd.Key {
	case "change-media":
		// Update the current playback ID in the players struct
		playback := database.FindOrCreatePlayback(cmd.Value)
		player, exists := players[devID]
		player.playback = playback
		players[devID] = player

		if !exists || players[devID].playback.ID == 0 {
			utils.ErrorC("controls error, playback doesn't exist")
		}

		response := jsonResponse(
			Response{
				Type:     "update",
				Key:      "change-media",
				Value:    "",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
		fmt.Printf("Device %s (MediaPlayback) change-media -> %d\n", devID, players[devID].playback.ID)
	case "play":
		response := jsonResponse(
			Response{
				Type:     "command",
				Key:      "play",
				Value:    "",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
	case "pause":
		response := jsonResponse(
			Response{
				Type:     "command",
				Key:      "pause",
				Value:    "",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
	case "toggle":
		response := jsonResponse(
			Response{
				Type:     "command",
				Key:      "toggle",
				Value:    "",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
	case "rewind":
		response := jsonResponse(
			Response{
				Type:     "command",
				Key:      "rewind",
				Value:    "10",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
	case "fastforward":
		response := jsonResponse(
			Response{
				Type:     "command",
				Key:      "fastforward",
				Value:    "10",
				Playback: players[devID].playback})

		SendToPlayers(response, devID)
	case "skip": // Find next ID, send to remotes + players, change channel ID, update details
		prefLoc := players[devID].playback.PrefLoc
		currentMedia := players[devID].playback.Locations[prefLoc]
		nextMedia, err := utils.GetNextMatchingOrderedFile(currentMedia)

		response := jsonResponse(
			Response{
				Type:     "error",
				Key:      "change-media-fail",
				Value:    "",
				Playback: database.Playback{}})

		if err == nil {
			fmt.Printf("findNextMedia found: %s\n", nextMedia.AbsPath)
			nextPlayback := database.FindOrCreatePlayback(nextMedia.AbsPath)

			player, exists := players[devID]
			player.playback = nextPlayback
			players[devID] = player

			if !exists || players[devID].playback.ID == 0 {
				utils.ErrorC("controls error, playback doesn't exist")
			}

			response = jsonResponse(
				Response{
					Type:     "update",
					Key:      "change-media",
					Value:    fmt.Sprintf("%d", nextPlayback.ID),
					Playback: nextPlayback})
		}

		SendToPlayers(response, devID)
		SendToRemotes(response, devID)
	default:
		fmt.Printf("Command: %+v", cmd)
		utils.ErrorC("Command Unknown")
	}
}

func status(cmd command, devID string) {
	switch cmd.Key {
	case "playback": // Update media playback
		playTime, _ := strconv.ParseFloat(cmd.Value, 64)
		database.UpdatePlaytime(players[devID].playback.ID, int(playTime))
	case "ping":
		response := jsonResponse(
			Response{
				Type:     "status",
				Key:      "pong",
				Value:    "",
				Playback: database.Playback{}})

		SendToPlayers(response, devID)
		SendToRemotes(response, devID)
	}
}

func jsonResponse(r Response) string {
	jsonStruct, err := json.Marshal(r)
	utils.Error("Couldn't convert response to Json", err)

	return string(jsonStruct)
}

func SendToPlayers(command string, devID string) {
	for i := 0; i < len(players[devID].players); i++ {
		err := players[devID].players[i].WriteMessage(1, []byte(command))
		utils.Error("Web socket closed", err)
	}
}

func SendToRemotes(command string, devID string) {
	for i := 0; i < len(players[devID].remotes); i++ {
		err := players[devID].remotes[i].WriteMessage(1, []byte(command))

		// Clean up remote stocket if the user refreshes
		if err != nil {
			newPlayers := players[devID]
			newPlayers.remotes = utils.RemoveSocketIndex(players[devID].remotes, i)
			players[devID] = newPlayers
		}
	}
}

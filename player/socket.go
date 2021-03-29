package player

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var channels = make(map[int]channel)

type channel struct {
	mediaInfo database.MediaInfo
	players   []*websocket.Conn
	remotes   []*websocket.Conn

	playCH   chan command
	remoteCH chan command
}

type command struct {
	ctype string
	key   string
	value string
}

func init() {
	go processor()
}

// SocketSetup creates web sockets, optional id parameter
func SocketSetup(w http.ResponseWriter, r *http.Request) {
	// Media ID as GET param
	mID, _ := strconv.Atoi(r.FormValue("id"))

	ws, err := upgrader.Upgrade(w, r, nil)
	utils.Err("Couldn't upgrade websocket", err)

	_, exists := channels[mID]
	var ch channel

	if !exists {
		ch = channel{
			mediaInfo: database.SelectMediaByID(mID),
			playCH:    make(chan command, 10),
			remoteCH:  make(chan command, 10),
		}

		if mux.Vars(r)["page"] == "remote" {
			ch.remotes = append(ch.remotes, ws)
			go remoteSocket(ws, mID)
		}

		if mux.Vars(r)["page"] == "player" {
			ch.players = append(ch.players, ws)
			go playerSocket(ws, mID)
		}

		channels[mID] = ch
	} else {
		var ch channel = channels[mID]

		if mux.Vars(r)["page"] == "remote" {
			ch.remotes = append(channels[mID].remotes, ws)
			go remoteSocket(ws, mID)
		}

		if mux.Vars(r)["page"] == "player" {
			ch.players = append(channels[mID].players, ws)
			go playerSocket(ws, mID)
		}

		channels[mID] = ch
	}

	fmt.Printf("Socket (%s) opened for id: %d (#chan: %d)\n", mux.Vars(r)["page"], mID, len(channels))
}

func processor() {
	for range time.Tick(300 * time.Millisecond) {
		for id, val := range channels {
			// Remote channels
			for i := 0; i < len(val.remoteCH); i++ {
				cmd := <-val.remoteCH

				switch cmd.ctype {
				case "control":
					controls(cmd, id)
				default:
					fmt.Println("Socket Processor doesn't recognise remote command")
				}
			}

			// Player channels
			for i := 0; i < len(val.playCH); i++ {
				cmd := <-val.playCH

				switch cmd.ctype {
				case "status":
					status(cmd, id)
				case "control":
					controls(cmd, id)
				default:
					fmt.Println("Socket Processor doesn't recognise player command")
				}
			}
		}
	}
}

func playerSocket(ws *websocket.Conn, id int) {
	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			return
		}

		cmd := strings.Split(string(msg), ":")

		if len(cmd) < 3 {
			cmd = append(cmd, "")
		}

		if cmd[0] == "change-id" {
			newID, _ := strconv.Atoi(cmd[1])
			fmt.Printf("Switching ID from %d -> %d\n", id, newID)
			id = newID

			continue
		}

		channels[id].playCH <- command{ctype: cmd[0], key: cmd[1], value: cmd[2]}

		fmt.Printf("Player %d -> %v\n", id, cmd)
	}
}

func remoteSocket(ws *websocket.Conn, id int) {
	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			return
		}

		cmd := strings.Split(string(msg), ":")

		if len(cmd) < 3 {
			cmd = append(cmd, "")
		}

		if cmd[0] == "change-id" {
			newID, _ := strconv.Atoi(cmd[1])
			fmt.Printf("Switching ID from %d -> %d\n", id, newID)
			id = newID

			continue
		}

		channels[id].remoteCH <- command{ctype: cmd[0], key: cmd[1], value: cmd[2]}

		fmt.Printf("Remote %d -> %v\n", id, cmd)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func controls(cmd command, id int) {
	switch cmd.key {
	case "play":
		sendToPlayers("command:play", id)
	case "pause":
		sendToPlayers("command:pause", id)
	case "rewind":
		sendToPlayers("command:rewind:10", id)
	case "fastforward":
		sendToPlayers("command:fastforward:10", id)
	case "skip": // Find next ID, send to remotes + players, change channel ID, update details
		nextID := findNextMedia(id)
		skipCMD := fmt.Sprintf("update:skip-to:%d", nextID)

		sendToPlayers(skipCMD, id)
		sendToRemotes(skipCMD, id)

		channels[nextID] = channels[id]
		delete(channels, id)

		media := database.SelectMediaByID(nextID)
		ret := fmt.Sprintf("update:media-info:%d;%s;%s;%s;%s;%d;%d", id, media.File.Name, media.Hash, media.File.AbsPath, media.File.Path, media.PlayTime, media.Date)

		sendToPlayers(ret, nextID)
		sendToRemotes(ret, nextID)
	}
}

func status(cmd command, id int) {
	switch cmd.key {
	case "playback":
		playbackUpdate(cmd.value, id)
	}
}

func sendToPlayers(command string, id int) {
	for i := 0; i < len(channels[id].players); i++ {
		channels[id].players[i].WriteMessage(1, []byte(fmt.Sprintf(command)))
	}
}

func sendToRemotes(command string, id int) {
	for i := 0; i < len(channels[id].remotes); i++ {
		channels[id].remotes[i].WriteMessage(1, []byte(fmt.Sprintf(command)))
	}
}

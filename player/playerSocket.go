package player

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var channels = make(map[int]channel)

type channel struct {
	mediaInfo database.MediaInfo
	c         chan command
}

type command struct {
	key   string
	value string
}

func commandSocket(ws *websocket.Conn, id int) {

	for {
		mt, msg, err := ws.ReadMessage()

		if err != nil {
			CloseChannel(id)
			return
		}

		// Incoming requests from client = key:value
		cmd := strings.Split(string(msg), ":")

		fmt.Printf("Chan %d/%d has %d commands: %v\n", id, len(channels), len(channels[id].c), cmd)

		// Handle incoming queries
		switch cmd[0] {
		case "close-channel":
			CloseChannel(id)
		case "playback":
			playbackUpdate(cmd[1], id)
		case "query":
			switch cmd[1] {
			case "nextID":
				if len(cmd) == 3 {
					queryID, _ := strconv.Atoi(cmd[2])
					ws.WriteMessage(mt, []byte(fmt.Sprintf("nextID:%d", findNextMedia(queryID))))
				} else {
					ws.WriteMessage(mt, []byte(fmt.Sprintf("nextID:%d", findNextMedia(id))))
				}
			case "prevID":
				// do something
			case "mediaInfo":
				media := database.SelectMediaByID(id)
				ret := fmt.Sprintf("id:%d:title:%s:hash:%s:path:%s:folder:%s:playtime:%d:date:%d", id, media.Title, media.Hash, media.Path, media.Folder, media.PlayTime, media.Date)
				ws.WriteMessage(mt, []byte(ret))
			}
		default:
			var com command

			if len(cmd) == 2 {
				com = command{key: cmd[1], value: ""}
			} else if len(cmd) == 3 {
				com = command{key: cmd[1], value: cmd[2]}
			}

			chID, _ := strconv.Atoi(cmd[0])
			channels[chID].c <- com
		}
	}
}

// SocketSetup creates web sockets, optional id parameter
func SocketSetup(w http.ResponseWriter, r *http.Request) {

	idParam := 0
	if r.FormValue("id") != "" {
		idParam, _ = strconv.Atoi(r.FormValue("id"))
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	utils.Err("Couldn't upgrade websocket", err)

	// Open a channel and socket
	OpenChannel(database.SelectMediaByID(idParam))

	go commandSocket(ws, idParam)

	if idParam != 0 {
		go receiverSocket(ws, idParam)
	}

	fmt.Printf("Socket opened for id: %d\n", idParam)
}

func receiverSocket(ws *websocket.Conn, id int) {
	for range time.Tick(300 * time.Millisecond) {
		var cmdBuff string

		_, exists := channels[id]

		if !exists {
			fmt.Printf("CHAN %d doesn't exist, breaking\n", id)
			break
		}

		for j := 0; j < len(channels[id].c); j++ {
			var cmd = <-channels[id].c
			cmdBuff += fmt.Sprintf("%s:%s:%s", cmdBuff, cmd.key, cmd.value)
		}

		cmdBuff = strings.TrimPrefix(cmdBuff, ":")
		cmdBuff = strings.TrimSuffix(cmdBuff, ":")

		ws.WriteMessage(1, []byte(cmdBuff))
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// OpenChannel opens up a new channel given a mediaInfo object
func OpenChannel(mediaInfo database.MediaInfo) {
	_, channelExists := channels[mediaInfo.ID]

	if !channelExists {
		fmt.Printf("Channel opened for id: %d\n", mediaInfo.ID)
		c := channel{mediaInfo: mediaInfo, c: make(chan command, 10)}
		channels[mediaInfo.ID] = c
	}
}

// CloseChannel closes a channel given an id
func CloseChannel(mediaID int) {
	_, channelExists := channels[mediaID]

	if channelExists {
		fmt.Printf("Closing channel %d\n", mediaID)
		delete(channels, mediaID)
	}
}

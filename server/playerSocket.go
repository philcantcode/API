package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func controlsSocket(ws *websocket.Conn) {
	for {
		mt, msg, err := ws.ReadMessage()

		if err != nil {
			log.Println("Error Reading Websocket:", err)
			break
		}

		msgArr := strings.Split(string(msg), ":")
		mediaID, _ := strconv.Atoi(msgArr[0])
		cmdKey := msgArr[1]
		cmdVal := msgArr[2]

		fmt.Printf("Chan %d/%d has %d cmds %s -> %s\n", mediaID, len(channels), len(channels[mediaID].c), cmdKey, cmdVal)

		// Processing incoming commands
		if cmdKey == "status" {
			playTime, _ := strconv.ParseFloat(cmdVal, 64)
			database.UpdatePlaytime(mediaID, int(playTime))
		} else if cmdKey == "query" {
			if cmdVal == "nextID" { // Client Player asks for the nextID
				prevMedia := database.SelectMediaByID(mediaID)
				nextMedia := utils.GetNextMatchingOrderedFile(prevMedia.Folder, prevMedia.Path)
				nextID := database.FindOrCreateMedia(nextMedia).ID

				ws.WriteMessage(mt, []byte(utils.JoinStr("nextID", fmt.Sprintf("%d", nextID))))
			} else if cmdVal == "mediaInfo" { // Client Player asks for all media info
				media := database.SelectMediaByID(mediaID)

				ws.WriteMessage(mt, []byte(utils.JoinStr(
					"id", fmt.Sprintf("%d", media.ID),
					"title", media.Title,
					"hash", media.Hash,
					"path", media.Path,
					"folder", media.Folder,
					"playtime", fmt.Sprintf("%d", media.PlayTime),
					"date", media.Date,
				)))
			}
		} else if cmdKey == "close" {
			CloseChannel(mediaID)
		} else if cmdKey == "open" {
			OpenChannel(database.SelectMediaByID(mediaID))
		} else {
			// Find the correct channel, new command to it
			cmd := command{key: cmdKey, value: cmdVal}
			channels[mediaID].c <- cmd
		}

		// Sending back commands to correct client
		if cmdKey == "status" {
			var returnCMD = ""
			channel, exists := channels[mediaID]

			if exists {
				for j := 0; j < len(channel.c); j++ {
					var cmd = <-channel.c
					returnCMD = utils.JoinStr(returnCMD, cmd.key, cmd.value)
				}

				ws.WriteMessage(mt, []byte(returnCMD))
			}
		}
	}
}

// WebSocketSetup creates web sockets
func WebSocketSetup(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	utils.Err("Couldn't upgrade websocket", err)
	go controlsSocket(ws)
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
		fmt.Printf("Opening new channel %d\n", mediaInfo.ID)
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

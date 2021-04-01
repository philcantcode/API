package player

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

var playerPage = page{
	PageTitle:       "Player",
	PageDescription: "Local Player",
	PreviousPath:    "Home",
	PreviousPathURL: "/",
	CurrentPath:     "player",
}

type RecentlyPlayed struct {
	Title string
	File  utils.File
}

// PlayerPage handles the /player request
func PlayerPage(w http.ResponseWriter, r *http.Request) {
	reload()

	openParam := r.FormValue("open")
	playParam := r.FormValue("play")

	data := struct {
		IP   string
		Port string

		OpenParam     string
		SafeOpenParam string

		// File selector menu containing Directories > Folders > Files
		Directories    []utils.File
		SubFolders     []utils.File
		Files          []utils.File
		RecentlyPlayed []RecentlyPlayed

		// Media that is being played
		MediaInfo database.MediaInfo
	}{
		IP:             utils.Host,
		Port:           utils.Port,
		Directories:    database.SelectDirectories(),
		OpenParam:      openParam,
		SafeOpenParam:  strings.ReplaceAll(openParam, "\\", "\\\\"),
		RecentlyPlayed: getRecentlyWatched(),
	}

	if playParam != "" {
		data.MediaInfo = database.FindOrCreateMedia(playParam)
	}

	if openParam != "" {
		// Find sub folders
		data.SubFolders = utils.GetFolderLayer(openParam)

		// Find files in folder for file display menu
		for _, s := range utils.GetFilesLayer(openParam) {
			// Hide subtitles from file display menu
			if s.Ext != ".srt" {
				data.Files = append(data.Files, s)
			}
			ConversionPriorityFolder = openParam
		}
	}

	playerPage.Contents = data
	templates.ExecuteTemplate(w, "player", playerPage)
}

func getRecentlyWatched() []RecentlyPlayed {
	// Get a past time period -n days
	timeRange := time.Now().AddDate(0, 0, -10).Unix()
	mediaList := database.SelectMediaByTime(timeRange)
	recent := make(map[string]database.MediaInfo)

	// Loop over all the returned media & group by folder titles
	for i := 0; i < len(mediaList); i++ {
		hasCategory := false

		// Loop over each token in the path
		for j := 0; j < len(mediaList[i].File.PathTokens); j++ {
			nthToken := mediaList[i].File.PathTokens[j]

			// Handles media ordered by category
			if strings.Contains(nthToken, "Category - ") {
				nextToken := mediaList[i].File.PathTokens[j+1]
				media, found := recent[nextToken]

				if !found {
					recent[nextToken] = mediaList[i]
				} else {
					if media.Date >= recent[nextToken].Date {
						recent[nextToken] = mediaList[i]
					}
				}

				hasCategory = true
			}
		}

		// If the media isn't ordered by a category
		if !hasCategory {
			folderName := mediaList[i].File.PathTokens[1]
			_, titleFound := recent[folderName]

			if !titleFound {
				recent[folderName] = mediaList[i]
			} else {
				if mediaList[i].Date > recent[folderName].Date {
					recent[folderName] = mediaList[i]
				}
			}
		}
	}

	var recentFiles []RecentlyPlayed
	var processed []string

	// Reorder by date
	for i := 0; i < len(recent); i++ {
		highest := 0
		highestStr := ""

		for title, value := range recent {
			if value.Date >= highest && !utils.Contains(title, processed) {
				highest = value.Date
				highestStr = title
			}
		}

		processed = append(processed, highestStr)
	}

	// Return
	for i := 0; i < len(processed); i++ {
		recentFiles = append(recentFiles, RecentlyPlayed{Title: processed[i], File: recent[processed[i]].File})
	}

	return recentFiles
}

// LoadMedia takes a file or ID GET param, then loads the media
func LoadMedia(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("file")
	id, _ := strconv.Atoi(r.FormValue("id"))

	// Find by ID
	if id != 0 {
		mediaInfo, err := database.SelectMediaByID(id)

		if err != nil {
			log.Fatalf("Couldn't retrieve media by ID: %d\n", id)
		}

		http.ServeFile(w, r, mediaInfo.File.AbsPath)
	}

	// find by file path
	if file != "" {
		mediaInfo := database.FindOrCreateMedia(file)
		http.ServeFile(w, r, mediaInfo.File.AbsPath)
	}
}

// When a playback update comes in
func playbackUpdate(playTimeStr string, mediaID int) {
	playTime, _ := strconv.ParseFloat(playTimeStr, 64)
	database.UpdatePlaytime(mediaID, int(playTime))
}

func findNextMedia(mediaID int) int {
	prevMedia, _ := database.SelectMediaByID(mediaID)
	nextMedia := utils.GetNextMatchingOrderedFile(utils.ProcessFile(prevMedia.File.AbsPath))
	nextID := database.FindOrCreateMedia(nextMedia).ID

	return nextID
}

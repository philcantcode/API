package player

import (
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
		IP       string
		Port     string
		DeviceID string

		OpenParam utils.File
		PlayParam utils.File

		SafeOpenParam string
		SafePlayParam string

		// File selector menu containing Directories > Folders > Files
		Directories    []utils.File
		SubFolders     []utils.File
		Files          []utils.File
		RecentlyPlayed []RecentlyPlayed

		// Media that is being played
		Playback database.Playback
	}{
		IP:             utils.Host,
		Port:           utils.Port,
		Directories:    database.SelectRootDirectories(),
		RecentlyPlayed: getRecentlyWatched(),
		DeviceID:       utils.RandomString(8),
	}

	if openParam != "" {
		data.OpenParam = utils.ProcessFile(openParam)
		data.SafeOpenParam = strings.ReplaceAll(openParam, "\\", "\\\\")
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

	if playParam != "" {
		data.PlayParam = utils.ProcessFile(playParam)
		data.SafePlayParam = strings.ReplaceAll(playParam, "\\", "\\\\")
		data.Playback = database.FindOrCreatePlayback(playParam)
	}

	playerPage.Contents = data
	templates.ExecuteTemplate(w, "player", playerPage)
}

// Returns the media last changed by the database (e.g., played)
func getRecentlyWatched() []RecentlyPlayed {
	// Get a past time period -n days
	timeRange := time.Now().AddDate(0, 0, -10).Unix()
	recentPlaybackList := database.SelectPlaybacks_ByTime(timeRange)
	recent := make(map[string]database.Playback)
	locations := make(map[string]utils.File)

	for _, playback := range recentPlaybackList {
		isIndexedByCategory := false

		// For each utils.File location
		for _, location := range playback.Locations {
			if location.Exists {
				// Loop over each token in the path
				for i := 0; i < len(location.PathTokens); i++ {
					nthToken := location.PathTokens[i]

					// Handles media ordered by category
					if strings.Contains(nthToken, "Category - ") && len(location.PathTokens) > i+1 {
						seriesName := location.PathTokens[i+1] // Token after category
						_, found := recent[seriesName]         // Already catalogued

						if !found { // Add to recent list
							recent[seriesName] = playback
							locations[seriesName] = location
							isIndexedByCategory = true
						} else { // Overwrite if bigger
							if playback.Modified >= recent[seriesName].Modified {
								recent[seriesName] = playback
								locations[seriesName] = location
								isIndexedByCategory = true
							}
						}
					}
				}

				// If the media isn't ordered by a category
				if !isIndexedByCategory {
					folderName := location.PathTokens[len(location.PathTokens)-2]

					for _, dir := range database.SelectRootDirectories() {

						if dir.Path == location.Path {
							folderName = location.PathTokens[len(location.PathTokens)-1] // File name
							break
						}
					}
					_, titleFound := recent[folderName]

					if !titleFound {
						recent[folderName] = playback
						locations[folderName] = location
					} else {
						if playback.Modified > recent[folderName].Modified {
							recent[folderName] = playback
							locations[folderName] = location
						}
					}
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
			if value.Modified >= highest && !utils.Contains(title, processed) {
				highest = value.Modified
				highestStr = title
			}
		}

		processed = append(processed, highestStr)
	}

	// Return
	for i := 0; i < len(processed); i++ {
		recentFiles = append(recentFiles, RecentlyPlayed{Title: processed[i], File: locations[processed[i]]})
	}

	return recentFiles
}

// LoadMedia takes a file or ID GET param, then loads the media
func LoadMedia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	utils.Error("Couldn't convert LoadMedia ID to integer", err)

	// Find by ID - the ID is guarenteed to already exist
	playback := database.SelectMediaPlayback_ByID(id)
	playback.PrefLoc = database.GetPreferredLocation(playback)
	http.ServeFile(w, r, playback.Locations[playback.PrefLoc].AbsPath)
}

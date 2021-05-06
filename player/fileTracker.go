package player

import (
	"os"
	"path/filepath"
	"time"

	"github.com/philcantcode/goApi/player/database"
)

type RecentlyModified struct {
	Path    string
	TimeAgo time.Duration
}

var RecentlyModifiedMedia []RecentlyModified

var nowTime time.Time

func FindRecentlyAdded() {
	for {
		nowTime = time.Now()

		drives := database.SelectFfmpegPriority()
		drives = append(drives, database.SelectRootDirectories()...)

		for i := 0; i < len(drives); i++ {
			filepath.Walk(drives[i].AbsPath, doFindRecentlyAdded)
		}

		// Recheck for new drives every n-seconds
		time.Sleep(20 * time.Second)
	}
}

func doFindRecentlyAdded(path string, info os.FileInfo, err error) error {

	if err == nil {
		modTime := info.ModTime()
		diffTime := nowTime.Sub(modTime)

		// Days * hours
		if diffTime < time.Duration(7*24*time.Hour) {
			RecentlyModifiedMedia = append(RecentlyModifiedMedia, RecentlyModified{Path: path, TimeAgo: diffTime})
		}
	}

	return nil
}

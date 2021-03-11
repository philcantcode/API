package utils

import (
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

// GetTime returns the current server time
func GetTime() int32 {
	return int32(time.Now().Unix())
}

// Err prints error messages
func Err(msg string, err error) {
	if err != nil {
		log.Fatalf("[%s] %s\n", msg, err)
	}
}

// ConvertMediaFile uses FFMPEG to convert to MP4
func ConvertMediaFile(file File) {
	exec.Command(FfmpegExe, "-i", file.PathWithName, file.Path+string(filepath.Separator)+file.NameWithoutExt+".mp4")
}

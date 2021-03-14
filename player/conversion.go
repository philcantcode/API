package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

const videoCodec = "libx264"
const audioCodec = "libmp3lame"

var ffmpegExe string

func init() {
	os := runtime.GOOS

	switch os {
	case "windows":
		ffmpegExe = "./res/ffmpeg-win7/bin/ffmpeg.exe"
	case "darwin":
		ffmpegExe = "/res/ffmpeg-osx/ffmpeg"
	case "linux":
		fmt.Println("OS Not Supported For File Conversion")
	default:
		fmt.Printf("%s.\n", os)
	}
}

// ConvertTrackedMediaDrives should be run on a new thread
func ConvertTrackedMediaDrives() {
	drives := database.SelectDirectories()

	for i := 0; i < len(drives); i++ {
		filepath.Walk(drives[i].Path, convertWalkFunc)
	}
}

func convertWalkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		ConvertToMP4(utils.ProcessFile(path), false, false)
	}

	return nil
}

// ConvertToMP4 uses FFMPEG to convert to MP4
func ConvertToMP4(file utils.File, stdout bool, remove bool) {

	if file.Ext != ".avi" { // && file.Ext != ".mkv" {
		return
	}

	// Full file path
	absFile := file.Path + file.Name + file.Ext

	fmt.Printf("Converting to MP4 [%s] %s\n", file.Ext, absFile)

	exec := exec.Command(ffmpegExe, "-hwaccel", "cuda", "-y", "-i", absFile, "-c:v", videoCodec, "-c:a", audioCodec, file.Path+file.Name+".mp4")

	if stdout {
		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
	}

	err := exec.Run()

	if err != nil {
		log.Fatal(err)
		return
	}

	if remove {
		os.Remove(absFile)
		return
	}

	sep := string(filepath.Separator)
	root := strings.Split(file.Path, sep)[0]
	convPath := root + sep + "conversion"

	os.Mkdir(convPath, 0755)
	os.Rename(absFile, convPath+sep+file.Name+file.Ext)
}

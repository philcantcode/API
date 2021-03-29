package player

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

const videoCodec = "libx264"
const audioCodec = "libmp3lame"

var ffmpegPath string
var ffmpegZip string

var FfmpegStat []ConversionHistory
var ConversionPriorityFolder = ""

type ConversionHistory struct {
	File      string
	Status    string
	StartTime time.Time
	EndTime   time.Time
}

func (p *ConversionHistory) NowTime() string {
	if p.EndTime.IsZero() {
		return fmt.Sprintf("%s", time.Since(p.StartTime).String())
	}

	return fmt.Sprintf("%s", p.EndTime.Sub(p.StartTime).String())
}

func init() {
	os := runtime.GOOS

	switch os {
	case "windows":
		ffmpegPath = "./res/ffmpeg/ffmpeg.exe"
		ffmpegZip = "res/ffmpeg/ffmpeg.zip"
		unzipFFMPEG()
		go ConvertTrackedMediaDrives()
	case "darwin":
		ffmpegPath = "/res/ffmpeg/ffmpeg-osx"
		ffmpegZip = "res/ffmpeg/ffmpeg-osx.zip"
		unzipFFMPEG()
		go ConvertTrackedMediaDrives()
	case "linux":
		fmt.Println("OS Not Supported For File Conversion")
	default:
		fmt.Printf("%s.\n", os)
	}
}

// ConvertTrackedMediaDrives should be run on a new thread
func ConvertTrackedMediaDrives() {
	for {
		drives := database.SelectDirectories()

		if ConversionPriorityFolder != "" {
			folder := ConversionPriorityFolder
			ConversionPriorityFolder = ""

			filepath.Walk(folder, convertWalkFunc)
		} else {
			for i := 0; i < len(drives); i++ {
				filepath.Walk(drives[i].Path, convertWalkFunc)

				if ConversionPriorityFolder != "" {
					break
				}
			}
		}
	}
}

func convertWalkFunc(path string, info os.FileInfo, err error) error {
	_, fErr := os.Stat(path)

	if fErr == nil && !info.IsDir() {
		ConvertToMP4(utils.ProcessFile(path), false, false)

		if ConversionPriorityFolder != "" {
			return io.EOF
		}
	}

	return nil
}

// ConvertToMP4 uses FFMPEG to convert to MP4
func ConvertToMP4(file utils.File, stdout bool, remove bool) {

	// Setup commands for file types
	if file.Ext != ".avi" && file.Ext != ".mkv" {
		return
	}

	var info = ConversionHistory{StartTime: time.Now()}
	info.File = file.Path + file.Name + file.Ext

	if !utils.IsLegalPath(info.File) {
		return
	}

	info.Status = "In Progress"
	FfmpegStat = append(FfmpegStat, info)
	pos := len(FfmpegStat) - 1

	newName := file.Name + ".mp4"
	newPath := file.Path + newName

	oldName := file.Name + file.Ext
	oldPath := file.Path + oldName

	fmt.Printf("Converting to MP4 [%s] %s\n", file.Ext, info.File)
	exec := exec.Command(ffmpegPath, "-hide_banner", "-loglevel", "error", "-hwaccel", "cuda", "-y", "-i", oldPath, "-c:v", videoCodec, "-c:a", audioCodec, newPath)

	if stdout {
		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
	}

	var outb, errb bytes.Buffer
	exec.Stdout = &outb
	exec.Stderr = &errb

	err := exec.Run()
	FfmpegStat[pos].EndTime = time.Now()

	if err != nil {
		fmt.Println("err:", errb.String())
		log.Println(err)
		FfmpegStat[pos].Status = "Error"
		return
	}

	FfmpegStat[pos].Status = "Success!"

	if remove {
		os.Remove(info.File)
		return
	}

	sep := string(filepath.Separator)
	root := strings.Split(file.Path, sep)[0]
	convPath := root + sep + ".ffmpeg"

	database.UpdateMediaName(oldName, newName)
	database.UpdateMediaPath(oldPath, newPath)

	os.Mkdir(convPath, 0755)
	os.Rename(info.File, convPath+sep+file.Name+file.Ext)
}

func unzipFFMPEG() {
	_, err := os.Stat(ffmpegPath)

	if os.IsNotExist(err) {
		fmt.Println("Unzipping " + ffmpegZip)
		exec := exec.Command("tar", "-xf", ffmpegZip, "--directory", "res/ffmpeg")
		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
		err := exec.Run()

		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

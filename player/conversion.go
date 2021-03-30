package player

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

const BROWSER_CODEC = "libx264"
const BROWSER_AUDIO = "libmp3lame"

var ffmpegPath string
var ffmpegZip string

var FfmpegStat []ConversionHistory
var ConversionPriorityFolder = ""

var codecFilter *regexp.Regexp
var audioFilter *regexp.Regexp

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
	codecFilter = regexp.MustCompile(`(?m)(Video: )([^\s]+)`)
	audioFilter = regexp.MustCompile(`(?m)(Audio: )([^\s]+)`)

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

	probeExec, _ := exec.Command(ffmpegPath, "-i", oldPath).CombinedOutput()
	codec := codecFilter.FindStringSubmatch(string(probeExec))[2]
	audio := audioFilter.FindStringSubmatch(string(probeExec))[2]

	var targVideo string
	var targAudio string

	fmt.Printf("Starting FFMPEG Conversion \n   > %s \n   > %s \n   > Codecs: %s / %s \n   > ", (file.Name + file.Ext), file.Path, codec, audio)

	switch codec {
	case "h264":
		fmt.Printf("Copying video container \n   > ")
		targVideo = "copy"
	case "hevc": // Full conversion
		fmt.Printf("Converting video to %s \n   > ", BROWSER_CODEC)
		targVideo = BROWSER_CODEC
	default: // Full conversion
		fmt.Printf("Default converting video to %s \n   > ", BROWSER_CODEC)
		targVideo = BROWSER_CODEC
	}

	switch audio {
	case "mp3":
		fmt.Printf("Copying audio container \n")
		targAudio = "copy"
	case "aac":
		fmt.Printf("Copying audio container \n")
		targAudio = "copy"
	default:
		fmt.Printf("Default converting audio to %s \n", BROWSER_AUDIO)
		targAudio = BROWSER_AUDIO
	}

	var ffmpeg = exec.Command(ffmpegPath, "-threads", "1", "-hide_banner", "-loglevel", "error", "-hwaccel", "cuda", "-y", "-i", oldPath, "-c:v", targVideo, "-c:a", targAudio, newPath)

	if stdout {
		ffmpeg.Stdout = os.Stdout
		ffmpeg.Stderr = os.Stderr
	}

	var outb, errb bytes.Buffer
	ffmpeg.Stdout = &outb
	ffmpeg.Stderr = &errb

	err := ffmpeg.Run()
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

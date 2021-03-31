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

var FfmpegStat []FfmpegMetrics
var ConversionPriorityFolder = ""

var codecFilter *regexp.Regexp
var audioFilter *regexp.Regexp

var NumFfmpegThreads int
var DisableFfmpeg = false

type FfmpegMetrics struct {
	File      string
	Status    string
	StartTime time.Time
	EndTime   time.Time
}

func (p *FfmpegMetrics) NowTime() string {
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

		// Recheck for new drives every n-seconds
		time.Sleep(20 * time.Second)
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

// ConvertToMP4 uses FFMPEG to convert media on tracked
// drives to MP4 by switching the container or changing
// the codec
func ConvertToMP4(file utils.File, stdout bool, remove bool) {

	if DisableFfmpeg == true {
		for {
			time.Sleep(30 * time.Second)

			if DisableFfmpeg == false {
				break
			}
		}
	}

	// Setup commands for file types
	if file.Ext != ".avi" && file.Ext != ".mkv" {
		return
	}

	var metrics = FfmpegMetrics{StartTime: time.Now()}
	metrics.File = file.Path + file.Name + file.Ext

	if !utils.IsLegalPath(metrics.File) {
		return
	}

	metrics.Status = "In Progress"
	FfmpegStat = append(FfmpegStat, metrics)
	pos := len(FfmpegStat) - 1

	mp4Name := file.Name + ".mp4"
	mp4Path := file.Path + mp4Name

	origName := file.Name + file.Ext
	origPath := file.Path + origName

	probeExec, _ := exec.Command(ffmpegPath, "-i", origPath).CombinedOutput()
	codec := codecFilter.FindStringSubmatch(string(probeExec))[2]
	audio := audioFilter.FindStringSubmatch(string(probeExec))[2]

	var targVideo string
	var targAudio string

	fmt.Printf("Starting FFMPEG (Threads: %d) \n   > %s \n   > %s \n   > Codecs: %s / %s \n   > ", NumFfmpegThreads, (file.Name + file.Ext), file.Path, codec, audio)

	// Setup video codec conversion
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

	// Setup audio codec conversion
	switch audio {
	case "aac":
		fmt.Printf("Copying audio container \n")
		targAudio = "copy"
	default:
		fmt.Printf("Default converting audio to %s \n", BROWSER_AUDIO)
		targAudio = BROWSER_AUDIO
	}

	// Run the command on terminal
	var ffmpeg = exec.Command(ffmpegPath, "-threads", fmt.Sprintf("%d", NumFfmpegThreads), "-hide_banner", "-loglevel", "error", "-hwaccel", "cuda", "-y", "-i", origPath, "-c:v", targVideo, "-c:a", targAudio, mp4Path)

	if stdout {
		ffmpeg.Stdout = os.Stdout
		ffmpeg.Stderr = os.Stderr
	}

	var outb, errb bytes.Buffer
	ffmpeg.Stdout = &outb
	ffmpeg.Stderr = &errb

	err := ffmpeg.Run()

	// Calculate duration of conversion
	FfmpegStat[pos].EndTime = time.Now()
	duration := fmt.Sprintf("%s", FfmpegStat[pos].EndTime.Sub(FfmpegStat[pos].StartTime).String())

	if err != nil {
		fmt.Println("err:", errb.String())
		log.Println(err)
		FfmpegStat[pos].Status = "Error"
		return
	}

	FfmpegStat[pos].Status = "Success!"

	// If set to delete file, bin it
	if remove {
		os.Remove(metrics.File)
		return
	}

	// Move the old file to .ffmpeg
	sep := string(filepath.Separator)
	root := strings.Split(file.Path, sep)[0]
	archiveFolder := root + sep + ".ffmpeg"
	archiveFile := archiveFolder + sep + file.Name + file.Ext

	// Make folder for .ffmpeg if doesn't exist
	os.Mkdir(archiveFolder, 0755)
	os.Rename(metrics.File, archiveFile)

	// Update ffmpeg + playhistoy database
	database.InsertFfmpeg(archiveFile, mp4Path, codec+" / "+audio, targVideo+" / "+targAudio, duration)
	database.UpdateMediaName(origName, mp4Name)
	database.UpdateMediaPath(origPath, mp4Path)
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

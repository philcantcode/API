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
	"time"

	"github.com/philcantcode/goApi/database"
	"github.com/philcantcode/goApi/utils"
)

const BROWSER_CODEC = "libx264"
const BROWSER_AUDIO = "libmp3lame"

var ffmpegPath string
var ffmpegZip string

var FfmpegStats []FfmpegMetrics
var ConversionPriorityFolder = ""

var codecFilter *regexp.Regexp
var audioFilter *regexp.Regexp

var NumFfmpegThreads int
var DisableFfmpeg = false

type FfmpegMetrics struct {
	File      utils.File
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
		drives := database.SelectFfmpegPriority()
		drives = append(drives, database.SelectRootDirectories()...)

		// There's a manual folder conversion priority set
		if ConversionPriorityFolder != "" {
			folder := ConversionPriorityFolder
			ConversionPriorityFolder = ""

			filepath.Walk(folder, convertWalkFunc)
		} else {
			for i := 0; i < len(drives); i++ {
				filepath.Walk(drives[i].AbsPath, convertWalkFunc)

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
		ConvertToMP4(path, false, false)

		if ConversionPriorityFolder != "" {
			return io.EOF
		}
	}

	return nil
}

// ConvertToMP4 uses FFMPEG to convert media on tracked
// drives to MP4 by switching the container or changing
// the codec
func ConvertToMP4(path string, stdout bool, remove bool) {
	if DisableFfmpeg == true {
		for {
			time.Sleep(30 * time.Second)

			if DisableFfmpeg == false {
				break
			}
		}
	}

	origPath := utils.ProcessFile(path)

	if !utils.IsLegalPath(origPath.AbsPath) {
		return
	}

	// Initial checks for the files
	if origPath.Ext != ".avi" && origPath.Ext != ".mkv" {
		return
	}

	origFile := database.FindOrCreatePlayback(path)
	mp4File := utils.ProcessFile(origPath.Path + origPath.FileName + ".mp4")

	// Metrics struct for tracking conversion progress
	var metrics = FfmpegMetrics{File: origPath, StartTime: time.Now(), Status: "In Progress"}

	for _, f := range FfmpegStats {
		if f.File.AbsPath == metrics.File.AbsPath {
			return
		}
	}

	FfmpegStats = append(FfmpegStats, metrics)
	pos := len(FfmpegStats) - 1

	probeExec, _ := exec.Command(ffmpegPath, "-i", origPath.AbsPath).CombinedOutput()
	codec := codecFilter.FindStringSubmatch(string(probeExec))[2]
	audio := audioFilter.FindStringSubmatch(string(probeExec))[2]

	var targVideo string
	var targAudio string

	fmt.Printf("Starting FFMPEG (Threads: %d) \n   > %s \n   > %s \n   > Codecs: %s / %s \n   > ", NumFfmpegThreads, (origPath.FileName + origPath.Ext), origPath.Path, codec, audio)

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
	var ffmpeg = exec.Command(ffmpegPath, "-threads", fmt.Sprintf("%d", NumFfmpegThreads), "-hide_banner", "-loglevel", "error", "-hwaccel", "cuda", "-y", "-i", origPath.AbsPath, "-c:v", targVideo, "-c:a", targAudio, mp4File.AbsPath)

	if stdout {
		ffmpeg.Stdout = os.Stdout
		ffmpeg.Stderr = os.Stderr
	}

	var outb, errb bytes.Buffer
	ffmpeg.Stdout = &outb
	ffmpeg.Stderr = &errb

	err := ffmpeg.Run()

	// Calculate duration of conversion
	FfmpegStats[pos].EndTime = time.Now()
	duration := fmt.Sprintf("%s", FfmpegStats[pos].EndTime.Sub(FfmpegStats[pos].StartTime).String())
	fmt.Printf("   > Duration %s \n", duration)

	if err != nil {
		fmt.Println("err:", errb.String())
		log.Println(err)
		FfmpegStats[pos].Status = "Error"
		return
	}

	FfmpegStats[pos].Status = "Success!"

	// If set to delete file, bin it
	if remove {
		os.Remove(metrics.File.AbsPath)
		return
	}

	// Move the old file to .ffmpeg
	sep := string(filepath.Separator)
	root := origPath.PathTokens[0]
	archiveFolder := root + sep + ".ffmpeg"
	archiveFile := archiveFolder + sep + origPath.FileName + origPath.Ext

	// Make folder for .ffmpeg if doesn't exist
	os.Mkdir(archiveFolder, 0755)
	os.Rename(metrics.File.AbsPath, archiveFile)
	mp4Hash, err := utils.MD5Hash(mp4File.AbsPath)
	utils.Error("Couldn't hash new MP4 file after conversion", err)

	database.InsertMediaHash(mp4Hash, mp4File.AbsPath, origFile.ID)
	database.InsertFfmpeg(mp4File.AbsPath, archiveFile, codec+" / "+audio, targVideo+" / "+targAudio, duration)
}

func unzipFFMPEG() {
	_, err := os.Stat(ffmpegPath)

	if os.IsNotExist(err) {
		fmt.Println("Unzipping " + ffmpegZip)
		exec := exec.Command("tar", "-xf", ffmpegZip, "--directory", "res/ffmpeg")
		exec.Stdout = os.Stdout
		exec.Stderr = os.Stderr
		err := exec.Run()

		utils.Error("unzipFFMPEG couldn't unzip the FFMPEG files", err)
	}
}

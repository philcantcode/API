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

	go generateExistingMediaHashes()
}

// ConvertTrackedMediaDrives should be run on a new thread
func ConvertTrackedMediaDrives() {
	for {
		drives := database.SelectFfmpegPriority()
		drives = append(drives, database.SelectDirectories()...)

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
	// Ensure the file exists in the playHistory
	targFile := database.FindOrCreateMedia(path).File
	mp4File := utils.ProcessFile(targFile.Path + targFile.Name + ".mp4")
	hash := database.SelectMediaHash(path)

	if DisableFfmpeg == true {
		for {
			time.Sleep(30 * time.Second)

			if DisableFfmpeg == false {
				break
			}
		}
	}

	// Initial checks for the files
	if targFile.Ext != ".avi" && targFile.Ext != ".mkv" {
		return
	}

	if !utils.IsLegalPath(targFile.AbsPath) {
		return
	}

	// Metrics struct for tracking conversion progress
	var metrics = FfmpegMetrics{StartTime: time.Now()}
	metrics.File = targFile

	metrics.Status = "In Progress"
	FfmpegStat = append(FfmpegStat, metrics)
	pos := len(FfmpegStat) - 1

	probeExec, _ := exec.Command(ffmpegPath, "-i", targFile.AbsPath).CombinedOutput()
	codec := codecFilter.FindStringSubmatch(string(probeExec))[2]
	audio := audioFilter.FindStringSubmatch(string(probeExec))[2]

	var targVideo string
	var targAudio string

	fmt.Printf("Starting FFMPEG (Threads: %d) \n   > %s \n   > %s \n   > Codecs: %s / %s \n   > ", NumFfmpegThreads, (targFile.Name + targFile.Ext), targFile.Path, codec, audio)

	// If the hash hasn't been created, generate it
	if len(hash) == 0 {
		fmt.Printf("Generating MD5 Hash\n   > ")
		hash, _ = utils.Hash(path)
		database.UpdateMediaHash(path, hash)
	}

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
	var ffmpeg = exec.Command(ffmpegPath, "-threads", fmt.Sprintf("%d", NumFfmpegThreads), "-hide_banner", "-loglevel", "error", "-hwaccel", "cuda", "-y", "-i", targFile.AbsPath, "-c:v", targVideo, "-c:a", targAudio, mp4File.AbsPath)

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
	fmt.Printf("   > Duration %s \n", duration)

	if err != nil {
		fmt.Println("err:", errb.String())
		log.Println(err)
		FfmpegStat[pos].Status = "Error"
		return
	}

	FfmpegStat[pos].Status = "Success!"

	// If set to delete file, bin it
	if remove {
		os.Remove(metrics.File.AbsPath)
		return
	}

	// Move the old file to .ffmpeg
	sep := string(filepath.Separator)
	root := strings.Split(targFile.Path, sep)[0]
	archiveFolder := root + sep + ".ffmpeg"
	archiveFile := archiveFolder + sep + targFile.Name + targFile.Ext

	// Make folder for .ffmpeg if doesn't exist
	os.Mkdir(archiveFolder, 0755)
	os.Rename(metrics.File.AbsPath, archiveFile)
	mp4Hash, _ := utils.Hash(mp4File.AbsPath)

	// Update ffmpeg + playhistoy database and generate altHash for new MP4 file
	database.InsertFfmpeg(archiveFile, mp4File.AbsPath, codec+" / "+audio, targVideo+" / "+targAudio, duration)
	database.UpdateMediaPathByHash(mp4File.AbsPath, hash)
	database.UpdateMediaAltHash(mp4File.AbsPath, mp4Hash)
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

// This function loops over all of the media items in the playHistory
// and attempts to generate any missing hashes
// Should be called in new background thread
func generateExistingMediaHashes() {
	print := true

	for {
		mediaList := database.SelectAllMedia()
		counter := 0

		for _, media := range mediaList {
			if media.Hash == "" {
				hash, err := utils.Hash(media.File.AbsPath)

				if err == nil {
					database.UpdateMediaHash(media.File.AbsPath, hash)
				}
			}

			counter++

			if counter%100 == 0 && print {
				fmt.Printf("Background Hasher: %d / %d\n", counter, len(mediaList))
			}

		}

		time.Sleep(120 * time.Second)
		print = false
	}
}

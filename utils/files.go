package utils

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var alphaNumFilter *regexp.Regexp

func init() {
	var err error
	alphaNumFilter, err = regexp.Compile("[^A-Za-z0-9]+")
	Error("Couldn't compile alphaNumFilter regex", err)
}

// File represents a file on OS
type File struct {
	FileName  string // No Extension
	PrintName string // No Extension, spaces
	Path      string // Folder
	Ext       string // .Extension

	AbsPath    string   // Path + Name + Ext
	PathTokens []string // Each piece of the file path
	Exists     bool     // If the file is found on disk
}

// GetFolderLayer returns a list of folders
func GetFolderLayer(path string) []File {
	files, _ := ioutil.ReadDir(path)
	var folders []File

	for _, file := range files {
		if file.IsDir() && IsLegalPath(file.Name()) {
			path := filepath.Join(path, file.Name())
			folders = append(folders, ProcessFile(path))
		}
	}

	return folders
}

// GetFilesLayer returns a list of files
func GetFilesLayer(path string) []File {
	files, _ := ioutil.ReadDir(path)
	var fileList []File

	for _, file := range files {
		if !file.IsDir() && IsLegalPath(file.Name()) {
			f := ProcessFile(filepath.Join(path, file.Name()))
			fileList = append(fileList, f)
		}
	}

	return fileList
}

func IsLegalPath(path string) bool {
	f := ProcessFile(path)

	if len(path) == 0 {
		return false
	}

	if f.FileName == "System Volume Information" {
		return false
	}

	// If any of the path tokens starts with illegal char
	for i := 0; i < len(f.PathTokens); i++ {
		switch string(f.PathTokens[i][0]) {
		case ".":
			return false
		case "$":
			return false
		}
	}

	return true
}

// ProcessFile extracts the file name from a path
func ProcessFile(path string) File {
	sep := string(filepath.Separator)
	file := File{}

	if len(path) == 0 {
		ErrorC("ProcessFile path length is 0: ")
	}

	// Split the file into parts by the system separator
	// e.g., /Users/Phil/Desktop/MediaTest becomes array['Users', 'Phil', 'Desktop', 'MediaTest']
	file.PathTokens = strings.Split(path, sep)

	// Last token becomes the temp file name
	tempFileName := file.PathTokens[len(file.PathTokens)-1]
	file.Path = strings.Join(file.PathTokens, sep)

	// Last token contains a . for EXT
	if strings.Contains(tempFileName, ".") {
		file.Path = strings.Join(file.PathTokens[0:len(file.PathTokens)-1], sep) + sep

		fileTokens := strings.Split(tempFileName, ".")
		fileExt := fileTokens[len(fileTokens)-1]

		file.Ext = "." + fileExt
		file.FileName = strings.Join(fileTokens[0:len(fileTokens)-1], ".")
	}

	file.PrintName = alphaNumFilter.ReplaceAllString(file.FileName, " ")
	file.AbsPath = file.Path + file.FileName + file.Ext

	// Reclice any file paths that are ''
	for i := 0; i < len(file.PathTokens); i++ {
		if file.PathTokens[i] == "" {
			file.PathTokens = RemoveIndex(file.PathTokens, i)
		}
	}

	if string(file.Path[len(file.Path)-1]) != sep {
		file.Path += sep
	}

	file.Exists = FileExists(file)

	return file
}

// GetDefaultSystemDrives returns a list of windows + Mac OS drives
func GetDefaultSystemDrives() []File {
	var drives []File

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			drives = append(drives, ProcessFile(string(drive)+":\\\\"))
			f.Close()
		}
	}

	usr, _ := user.Current()
	drives = append(drives, ProcessFile(usr.HomeDir))
	drives = append(drives, ProcessFile("/Volumes"))

	return drives
}

// GetNextMatchingOrderedFile Takes in a folder and file, returns the next ordered file or returns "" if none found
func GetNextMatchingOrderedFile(file File) string {
	files := GetFilesLayer(file.Path)

	for i := 0; i < len(files); i++ {
		if files[i].FileName == file.FileName {
			for j := i + 1; j < len(files); j++ {
				if files[j].Ext == file.Ext {
					return files[j].Path + files[j].FileName + files[j].Ext
				}
			}
		}
	}

	return ""
}

// FileExists returns true/false as to whether the file exists
func FileExists(f File) bool {
	_, err := os.Stat(f.AbsPath)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

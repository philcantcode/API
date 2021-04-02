package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var alphaNumFilter *regexp.Regexp

func init() {
	alphaNumFilter, _ = regexp.Compile("[^A-Za-z0-9]+")
}

// File represents a file on OS
type File struct {
	Name      string // No Extension
	PrintName string // No Extension, spaces
	Path      string // Folder
	Ext       string // .Extension

	AbsPath    string   // Path + Name + Ext
	PathTokens []string // Each piece of the file path
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

	if f.Name == "System Volume Information" {
		return false
	}

	// If any of the path tokens starts with illegal char
	for i := 0; i < len(f.PathTokens); i++ {
		fmt.Printf("::: %s\n", f.AbsPath)
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

	// If there is a system separator in the path
	if strings.Contains(path, sep) {
		tokens := strings.Split(path, sep)
		fileName := tokens[len(tokens)-1] // Last token = file name
		filePath := strings.Join(tokens[0:len(tokens)-1], sep)

		file.Path = filePath + sep
		file.PathTokens = tokens

		// If contains an extension
		if strings.Contains(fileName, ".") {
			fileTokens := strings.Split(fileName, ".")
			fileExt := fileTokens[len(fileTokens)-1]
			fileName := strings.Join(fileTokens[0:len(fileTokens)-1], ".")

			file.Ext = "." + fileExt
			file.Name = fileName

			fileName = alphaNumFilter.ReplaceAllString(fileName, " ")
			file.PrintName = fileName
		} else {
			file.Name = fileName
		}
	} else { // Single file or folder
		file.PathTokens = []string{path}

		if strings.Contains(path, ".") {
			fileTokens := strings.Split(path, ".")
			fileExt := fileTokens[len(fileTokens)-1]
			fileName := strings.Join(fileTokens[0:len(fileTokens)-1], ".")

			file.Ext = "." + fileExt
			file.Name = fileName

			fileName = alphaNumFilter.ReplaceAllString(fileName, " ")
			file.PrintName = fileName
		} else {
			file.Name = path
		}
	}

	file.AbsPath = file.Path + file.Name + file.Ext

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
		if files[i].Name == file.Name {
			for j := i + 1; j < len(files); j++ {
				if files[j].Ext == file.Ext {
					return files[j].Path + files[j].Name + files[j].Ext
				}
			}
		}
	}

	return ""
}

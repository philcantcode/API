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
	alphaNumFilter, _ = regexp.Compile("[^A-Za-z0-9]+")
}

// File represents a file on OS
type File struct {
	Name      string // No Extension
	PrintName string // No Extension, spaces
	Path      string // Folder
	Ext       string // .Extension
}

// GetFolderLayer returns a list of folders
func GetFolderLayer(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var folders []string

	for _, file := range files {
		if file.IsDir() && isLegalPath(file.Name()) {
			folders = append(folders, filepath.Join(path, file.Name()))
		}
	}

	return folders
}

// GetFilesLayer returns a list of files
func GetFilesLayer(path string) []File {
	files, _ := ioutil.ReadDir(path)
	var fileList []File

	for _, file := range files {
		if !file.IsDir() && isLegalPath(file.Name()) {
			f := ProcessFile(filepath.Join(path, file.Name()))
			fileList = append(fileList, f)
		}
	}

	return fileList
}

func isLegalPath(path string) bool {

	if len(path) == 0 {
		return false
	}

	if path == "System Volume Information" {
		return false
	}

	switch string(path[0]) {
	case ".":
		return false
	case "$":
		return false
	default:
		return true
	}
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

	}

	//fmt.Printf("File: %s [%s] \n\t Printable: %s \n\t Path: %s \n\t", file.Name, file.Ext, file.PrintName, file.Path)

	return file
}

// GetDrives returns a list of windows OS drives
func GetDrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive)+":\\\\")
			f.Close()
		}
	}

	usr, _ := user.Current()
	r = append(r, usr.HomeDir)
	r = append(r, "/Volumes")

	return
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

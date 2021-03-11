package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// File represents a file on OS
type File struct {
	NameWithExt    string
	NameWithoutExt string
	NameWithSpaces string
	PathWithName   string
	Path           string
	Extension      string
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
			f := File{NameWithExt: file.Name(), PathWithName: filepath.Join(path, file.Name())}
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

// ExtractFileName extracts the file name from a path
func ExtractFileName(fileName string) File {

	file := File{PathWithName: fileName}

	paths := strings.Split(fileName, string(filepath.Separator))
	fileName = paths[len(paths)-1]

	file.Path = strings.Join(paths[0:len(paths)-1], string(filepath.Separator))
	file.NameWithExt = fileName

	name := strings.Split(fileName, ".")

	file.NameWithoutExt = name[0]
	file.Extension = name[len(name)-1]
	file.NameWithSpaces = strings.ReplaceAll(fileName, ".", " ")

	fmt.Printf("New File: %+v", file)

	return file
}

// ExtractFolderName extracts the file name from a path
func ExtractFolderName(path string) string {
	paths := strings.Split(path, string(filepath.Separator))

	return strings.Join(paths[:len(paths)-1], string(filepath.Separator))
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

	return
}

// GetNextMatchingOrderedFile Takes in a folder and file, returns the next ordered file or returns "" if none found
func GetNextMatchingOrderedFile(folder string, file string) string {
	files := GetFilesLayer(folder)

	for i := 0; i < len(files); i++ {
		if files[i].PathWithName == file {
			for j := i + 1; j < len(files); j++ {
				if filepath.Ext(files[j].PathWithName) == filepath.Ext(file) {
					return files[j].PathWithName
				}
			}
		}
	}

	return ""
}

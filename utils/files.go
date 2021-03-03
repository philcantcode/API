package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// File represents a file on OS
type File struct {
	Name string
	Path string
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
			f := File{Name: file.Name(), Path: filepath.Join(path, file.Name())}
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
func ExtractFileName(path string) string {
	paths := strings.Split(path, string(filepath.Separator))
	path = paths[len(paths)-1]

	return strings.ReplaceAll(path, ".", " ")
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
	return
}

// GetNextMatchingOrderedFile Takes in a folder and file, returns the next ordered file or returns "" if none found
func GetNextMatchingOrderedFile(folder string, file string) string {
	files := GetFilesLayer(folder)

	for i := 0; i < len(files); i++ {
		if files[i].Path == file {
			for j := i + 1; j < len(files); j++ {
				if filepath.Ext(files[j].Path) == filepath.Ext(file) {
					return files[j].Path
				}
			}
		}
	}

	return ""
}

package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Node struct {
	Path  string
	Nodes []Node
}

type File struct {
	Name string
	Path string
}

// GetFolderTree returns a list of all folders
func GetFolderTree(parent Node) Node {
	files, err := ioutil.ReadDir(parent.Path)

	if err != nil {
		return parent
	}

	for _, file := range files {
		if file.IsDir() {

			var subdir Node
			subdir.Path = filepath.Join(parent.Path, file.Name())
			subdir.Nodes = make([]Node, 0)

			parent.Nodes = append(parent.Nodes, GetFolderTree(subdir))
		}
	}

	return parent
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
			var f File
			f.Name = file.Name()
			f.Path = filepath.Join(path, file.Name())

			fileList = append(fileList, f)
		}
	}

	return fileList
}

func isLegalPath(path string) bool {
	if strings.HasPrefix(path, ".") {
		return false
	}

	if strings.HasPrefix(path, "$") {
		return false
	}

	return true
}

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

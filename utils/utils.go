package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

const printError = false

// Fancy log function
func Error(info string, err error) {
	if err == nil {
		return
	}

	_, fn, line, _ := runtime.Caller(1)

	path := strings.Split(fn, string(filepath.Separator))
	fn = path[len(path)-1]

	if printError {
		fmt.Printf("%v\n\n\n[Error] %s in %s, Line: %d\n\n", err, info, fn, line)
	} else {
		fmt.Printf("[Error] %s in %s, Line: %d\n\n", info, fn, line)
	}

	os.Exit(0)
}

// Fancy log function
func ErrorC(info string) {
	_, fn, line, _ := runtime.Caller(1)

	path := strings.Split(fn, string(filepath.Separator))
	fn = path[len(path)-1]

	fmt.Printf("[Error] %s in %s, Line: %d\n\n", info, fn, line)

	os.Exit(0)
}

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MD5Hash(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)

	if err != nil {
		return returnMD5String, err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func RemoveStrIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func RemoveSocketIndex(s []*websocket.Conn, index int) []*websocket.Conn {
	return append(s[:index], s[index+1:]...)
}

// RandomString generates a random string
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

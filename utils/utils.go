package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// Err prints error messages
func Err(msg string, err error) {
	if err != nil {
		log.Fatalf("[%s] %s\n", msg, err)
	}
}

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Hash(filePath string) (string, error) {
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

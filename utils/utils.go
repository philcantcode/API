package utils

import (
	"log"
	"time"
)

// GetTime returns the current server time
func GetTime() int32 {
	return int32(time.Now().Unix())
}

// Err prints error messages
func Err(msg string, err error) {
	if err != nil {
		log.Fatalf("[%s] %s\n", msg, err)
	}
}

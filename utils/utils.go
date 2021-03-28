package utils

import (
	"log"
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

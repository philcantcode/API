package utils

import (
	"log"
	"strings"
	"time"
)

// GetTime returns the current server time
func GetTime() int32 {
	return int32(time.Now().Unix())
}

func Err(msg string, err error) {
	if err != nil {
		log.Fatalf("[%s] %s\n", msg, err)
	}
}

func JoinStr(s ...string) string {
	ret := ""

	ret = strings.Join(s, ":")
	ret = strings.TrimPrefix(ret, ":")
	ret = strings.TrimSuffix(ret, ":")

	return ret
}

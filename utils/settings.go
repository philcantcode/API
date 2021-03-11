package utils

import (
	"fmt"
	"log"
	"net"
	"runtime"
)

// Flags
var FilePath string
var Host string
var Port string
var DBLoc string
var FfmpegExe string

func init() {
	os := runtime.GOOS

	FilePath = "web/"
	Port = "9002"
	Host = GetOutboundIP().String()

	fmt.Printf("Server Launched at: http://%s:%s\n", Host, Port)

	switch os {
	case "windows":
		DBLoc = "C:/Users/Phil/Google Drive/elements.db"
	case "darwin":
		DBLoc = "/Users/Phil/Google Drive/elements.db"
		FfmpegExe = "/res/ffmpeg-osx"
	case "linux":
		fmt.Println("OS Not Supported")
	default:
		fmt.Printf("%s.\n", os)
	}
}

// GetOutboundIP Gets preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

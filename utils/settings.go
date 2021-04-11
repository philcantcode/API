package utils

import (
	"fmt"
	"log"
	"net"
	"os"
)

var FilePath string
var Host string
var Port string

func init() {
	FilePath = "web/"
	Port = "9002"
	Host = GetOutboundIP().String()

	fmt.Printf("Server Launched at: http://%s:%s\n", Host, Port)
	fmt.Printf("Server Launched at: http://%s.local:%s\n", GetHostname(), Port)

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

// GetHostname returns the computer host name
func GetHostname() string {
	hostName, err := os.Hostname()
	Error("Couldn't get computer hostname", err)

	return hostName
}

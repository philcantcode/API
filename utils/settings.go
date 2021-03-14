package utils

import (
	"fmt"
	"log"
	"net"
)

var FilePath string
var Host string
var Port string

func init() {
	FilePath = "web/"
	Port = "9002"
	Host = GetOutboundIP().String()

	fmt.Printf("Server Launched at: http://%s:%s\n", Host, Port)
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

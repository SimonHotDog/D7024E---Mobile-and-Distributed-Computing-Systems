package network

import (
	"log"
	"net"
)

// Get IP-address of this computer
func GetOutboundIP() string {
	// https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

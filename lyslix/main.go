package main

import (
	"fmt"
	"net"
	"os"
)

const LISTEN_ADDR string = "0.0.0.0"
const DEFAULT_PORT string = "56700"

func getUdpAddr() (*net.UDPAddr, error) {
	listenAddr := LISTEN_ADDR
	portNumber := os.Getenv("PORT")
	if len(portNumber) == 0 {
		portNumber = DEFAULT_PORT
	}
	strAddr := fmt.Sprintf("%s:%s", listenAddr, portNumber)
	return net.ResolveUDPAddr("udp", strAddr)
}

func main() {
	fmt.Println("lyslix 1.0")

	addr, err := getUdpAddr()
	if err != nil {
		return
	}

	_ = addr
}

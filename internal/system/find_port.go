package system

import (
	"fmt"
	"net"
)

const defaultPort = 7777

func FindPort() (int, int, bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", defaultPort))
	if err == nil {
		ln.Close()
		return defaultPort, defaultPort, false
	}

	ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return defaultPort, 0, false
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return defaultPort, port, true
}

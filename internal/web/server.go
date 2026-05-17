package web

import (
	"fmt"
	"net"
	"net/http"

	"statehound/internal/logger"
)

const defaultPort = 7777

func Start() {
	port, fallback := findPort()
	if port == 0 {
		logger.Failed("failed to find available port", nil)
		return
	}

	if fallback {
		logger.Warn(fmt.Sprintf("port %d is in use, using port %d instead", defaultPort, port))
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := "http://" + addr

	mux := http.NewServeMux()
	registerRoutes(mux)

	logger.Success("statehound web started at " + url)
	logger.Status("open in browser: " + url)

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Failed("web server error", err)
	}
}

func findPort() (int, bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", defaultPort))
	if err == nil {
		ln.Close()
		return defaultPort, false
	}

	ln, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, false
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port, true
}

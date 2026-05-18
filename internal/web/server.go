package web

import (
	"fmt"
	"net/http"

	"statehound/internal/logger"
	"statehound/internal/system"
)

func Start() {
	defaultPort, port, fallback := system.FindPort()
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

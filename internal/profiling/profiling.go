package profiling

import (
	_ "embed"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"statehound/internal/logger"
	"statehound/internal/system"
)

//go:embed static/profiling.html
var indexHTML string

func handleIndex(addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexHTML))
	}
}

func Start() string {
	_, port, fallback := system.FindPort()
	if port == 0 {
		return "failed to find available port for pprof"
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	mux := http.NewServeMux()

	// keep all pprof endpoints working
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	// pretty index
	mux.HandleFunc("/", handleIndex(addr))

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Failed("pprof server error", err)
		}
	}()

	msg := "pprof started at http://" + addr + "/"
	if fallback {
		msg += " (port 6060 was in use)"
	}
	return msg
}

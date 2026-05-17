package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/events", handleEvents)
	mux.HandleFunc("/api/snapshot", handleSnapshot)
	mux.HandleFunc("/api/hunt", handleHunt)

	sub, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/", http.FileServer(http.FS(sub)))
}

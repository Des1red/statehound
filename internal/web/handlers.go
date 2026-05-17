package web

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"statehound/internal/client"
	"statehound/internal/model"
	"statehound/internal/statehound/events"
)

func handleEvents(w http.ResponseWriter, r *http.Request) {
	var lines []string
	var newOffset int64

	offsetStr := r.URL.Query().Get("offset")

	if offsetStr == "" {
		// initial load — return all existing events + current offset
		content, err := events.FilterEvents("")
		if err != nil {
			http.Error(w, "failed to read events", http.StatusInternalServerError)
			return
		}
		if content != "" {
			for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
				if line != "" {
					lines = append(lines, line)
				}
			}
		}
		info, _ := os.Stat(model.EventPath)
		if info != nil {
			newOffset = info.Size()
		}
	} else {
		// subsequent poll — return only new lines since offset
		offset, _ := strconv.ParseInt(offsetStr, 10, 64)
		var err error
		lines, newOffset, err = events.EventsSince(offset)
		if err != nil {
			http.Error(w, "failed to read events", http.StatusInternalServerError)
			return
		}
	}

	if lines == nil {
		lines = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": lines,
		"offset": newOffset,
	})
}

func handleSnapshot(w http.ResponseWriter, r *http.Request) {
	resp, err := client.Send("SNAPSHOT")
	if err != nil {
		http.Error(w, "daemon not responding", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"snapshot": resp})
}

func handleHunt(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "missing target", http.StatusBadRequest)
		return
	}

	resp, err := client.Send("HUNT " + target)
	if err != nil {
		http.Error(w, "daemon not responding", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": resp})
}

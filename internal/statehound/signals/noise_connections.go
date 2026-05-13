package signals

import (
	"statehound/internal/statehound/collector"
	"strings"
)

func isNoiseConnection(conn collector.Connection) bool {
	process := strings.ToLower(conn.Process)
	switch process {
	case "firefox", "brave", "chrome", "chromium":
		return true
	}
	return false
}

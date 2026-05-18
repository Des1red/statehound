package signals

import (
	"statehound/internal/statehound/collector"
	"strings"
)

func isNoisePort(port collector.Port) bool {
	return isBrowserUDP(strings.ToLower(port.Process), port.Proto)
}

func isBrowserUDP(process, proto string) bool {
	if proto != "udp" {
		return false
	}

	switch process {
	case "firefox", "firefox-esr", "brave", "chrome", "chromium":
		return true
	default:
		return false
	}
}

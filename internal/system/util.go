package system

import (
	"os"
	"statehound/internal/command"
	"strings"
)

func IsHeadless() bool {
	if os.Getenv("DISPLAY") != "" {
		return false
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}

	out, err := command.Output("systemctl", "get-default")
	if err == nil {
		if strings.Contains(strings.TrimSpace(string(out)), "graphical") {
			return false
		}
	}

	return true
}

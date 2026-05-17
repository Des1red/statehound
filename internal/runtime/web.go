package runtime

import (
	"statehound/internal/logger"
	"statehound/internal/system"
	"statehound/internal/web"
)

func Web() {
	if system.IsHeadless() {
		logger.Status("headless system detected — no browser will open")
		logger.Status("to access the dashboard, SSH tunnel to this machine:")
		logger.Status("  ssh -L 7777:127.0.0.1:7777 user@<host>")
	}
	web.Start()
}

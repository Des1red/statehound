package viewer

import (
	"os/exec"
	"statehound/internal/logger"
	"statehound/internal/statehound/events"
	"strings"
)

func Show(urgency string) {
	content, err := events.FilterEvents(urgency)
	if err != nil {
		logger.Failed("viewer: failed to filter events", err)
		return
	}
	if content == "" {
		content = "no " + urgency + " events recorded"
	}

	cmd := exec.Command(
		"zenity",
		"--text-info",
		"--title=statehound — "+urgency+" events",
		"--width=900",
		"--height=600",
		"--font=monospace",
	)
	cmd.Stdin = strings.NewReader(content)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Failed("viewer: zenity failed: "+string(out), err)
	}
}

package runtime

import (
	"fmt"
	"os"

	"statehound/internal/logger"
	"statehound/internal/model"
	"statehound/internal/statehound/events"
	"statehound/internal/system"
	"statehound/internal/viewer"
)

func ShowEvents(filter string) {
	content, err := events.FilterEvents(filter)
	if err != nil {
		logger.Failed("failed to read statehound events", err)
		return
	}
	if content == "" {
		logger.Status("no statehound events recorded yet")
		return
	}
	fmt.Print(content)
}

func EventsGUI(filter string) {
	if system.IsHeadless() {
		logger.Warn("no display available — use --events instead")
		return
	}
	if filter == "" || (filter != "critical" && filter != "normal") {
		filter = "critical"
	}

	viewer.Show(filter)
}

func ClearEvents() {
	if err := os.WriteFile(model.EventPath, []byte{}, 0640); err != nil {
		logger.Failed("failed to clear statehound events", err)
		return
	}

	logger.Success("statehound events cleared")
}

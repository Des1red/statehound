package runtime

import (
	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Log() {
	if err := command.Run("journalctl", "-u", model.ServiceName, "--no-pager", "-n", "50"); err != nil {
		logger.Failed("failed to read statehound logs", err)
	}
}

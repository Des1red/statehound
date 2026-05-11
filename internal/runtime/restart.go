package runtime

import (
	"statehound/internal/client"
	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Restart() {
	if err := command.Run("systemctl", "restart", model.ServiceName); err != nil {
		logger.Failed("failed to restart statehound", err)
		return
	}

	if !client.IsRunning() {
		logger.Failed("statehound restart command was sent, but daemon is not responding", nil)
		return
	}

	logger.Success("statehound restarted")
}

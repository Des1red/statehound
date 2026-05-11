package runtime

import (
	"statehound/internal/client"
	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Start() {
	if client.IsRunning() {
		logger.Status("statehound is already running")
		return
	}

	if err := command.Run("systemctl", "start", model.ServiceName); err != nil {
		logger.Failed("failed to start statehound", err)
		return
	}

	if !client.IsRunning() {
		logger.Failed("statehound start command was sent, but daemon is not responding", nil)
		return
	}

	logger.Success("statehound started")
}

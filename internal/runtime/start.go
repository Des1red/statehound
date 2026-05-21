package runtime

import (
	"statehound/internal/client"
	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Start() {
	if !client.IsRunning() {
		if err := command.Run("systemctl", "start", model.ServiceName); err != nil {
			logger.Failed("failed to start statehound daemon", err)
			return
		}

		if !client.IsRunning() {
			logger.Failed("statehound start command was sent, but daemon is not responding", nil)
			return
		}
	} else {
		logger.Status("statehound daemon is already running")
	}

	if err := command.Run("systemctl", "--user", "start", model.NotifierServiceName); err != nil {
		logger.Warn("statehound notifier was not started")
		logger.Warn(err.Error())
	} else {
		logger.Status("statehound notifier started")
	}

	logger.Success("statehound started")
}

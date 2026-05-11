package runtime

import (
	"statehound/internal/client"
	"statehound/internal/logger"
)

func Status() {
	if !client.IsRunning() {
		logger.Warn("statehound daemon is not running or not responding")
		return
	}

	resp, err := client.Send("STATUS")
	if err != nil {
		logger.Failed("failed to read statehound status", err)
		return
	}

	logger.Success(resp)
}

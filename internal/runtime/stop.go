package runtime

import (
	"statehound/internal/client"
	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Stop() {
	if !client.IsRunning() {
		logger.Status("statehound is already stopped")
		return
	}

	if err := command.Run("systemctl", "stop", model.ServiceName); err != nil {
		logger.Failed("[!] failed to stop statehound:", err)
		return
	}

	if client.IsRunning() {
		logger.Failed("[!] statehound stop command was sent, but daemon is still responding", nil)
		return
	}

	logger.Success("[+] statehound stopped")
}

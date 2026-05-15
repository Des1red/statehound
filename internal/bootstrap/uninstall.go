package bootstrap

import (
	"os"

	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func Uninstall(purge bool) {
	logger.Status("uninstalling statehound")

	cleanupNotifierForSudoUser()

	if err := command.Run("systemctl", "stop", model.ServiceName); err != nil {
		logger.Warn("service was not running or could not be stopped")
	}

	if err := command.Run("systemctl", "disable", model.ServiceName); err != nil {
		logger.Warn("service was not enabled or could not be disabled")
	}

	if err := os.Remove(model.ServicePath); err != nil && !os.IsNotExist(err) {
		logger.Failed("failed to remove systemd service", err)
		return
	}

	if err := os.Remove(model.NotifierServicePath); err != nil && !os.IsNotExist(err) {
		logger.Failed("failed to remove notifier user service", err)
		return
	}

	if err := os.Remove(model.AliasPath); err != nil && !os.IsNotExist(err) {
		logger.Failed("failed to remove shound alias", err)
		return
	}

	if err := os.Remove(model.BinaryPath); err != nil && !os.IsNotExist(err) {
		logger.Failed("failed to remove statehound binary", err)
		return
	}

	if purge {
		removeInstalledDeps()

		if err := os.RemoveAll(model.ConfigDir); err != nil {
			logger.Failed("failed to remove statehound config directory", err)
			return
		}

		if err := os.RemoveAll(model.LogDir); err != nil {
			logger.Failed("failed to remove statehound log directory", err)
			return
		}

		logger.Status("removed statehound config and event logs")
	}

	if err := command.Run("systemctl", "daemon-reload"); err != nil {
		logger.Failed("failed to reload systemd", err)
		return
	}

	_ = command.Run("systemctl", "reset-failed", model.ServiceName)

	if purge {
		logger.Success("statehound uninstalled and purged")
		return
	}
	logger.Success("statehound uninstalled")
}

package bootstrap

import (
	"os"

	"statehound/internal/logger"
	"statehound/internal/model"
)

func VerifyInstallation() bool {
	required := []string{
		model.BinaryPath,
		model.AliasPath,
		model.ServicePath,
		model.NotifierServicePath,
		model.ConfigDir,
		model.LogDir,
		model.EventBackupDir,
		model.EventPath,
	}

	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			logger.Failed("statehound is not installed correctly", err)
			logger.Status("Please run: sudo go run . --install")
			return false
		}
	}

	ensureNotifierDeps()

	if err := setupNotifierForSudoUser(); err != nil {
		logger.Warn("desktop notifier was not enabled automatically")
		logger.Warn(err.Error())
		logger.Status("Please run: sudo go run . --install")
		return false
	}
	return true
}

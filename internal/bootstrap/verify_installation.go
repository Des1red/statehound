package bootstrap

import (
	"os"

	"statehound/internal/logger"
	"statehound/internal/model"
	"statehound/internal/system"
)

func VerifyInstallation() bool {
	required := []string{
		model.BinaryPath,
		model.AliasPath,
		model.ServicePath,
		model.ConfigDir,
		model.LogDir,
		model.EventBackupDir,
		model.EventPath,
	}

	if !system.IsHeadless() {
		required = append(required, model.NotifierServicePath)
	}

	for _, path := range required {
		if _, err := os.Stat(path); err != nil {
			logger.Failed("statehound is not installed correctly", err)
			logger.Status("Please run: sudo go run . --install")
			return false
		}
	}

	return true
}

package daemon

import (
	"time"

	"statehound/internal/logger"
	"statehound/internal/model"
	"statehound/internal/statehound"
	"statehound/internal/system"
)

func Run() {
	logger.Status("statehound daemon started")

	manager := statehound.NewManager(5 * time.Second)
	go manager.Run()

	if err := prepareRuntimeDir(); err != nil {
		logger.Failed("failed to prepare runtime directory", err)
		return
	}

	if !system.IsHeadless() {
		go func() {
			if err := startNotifySocket(model.NotifySocketPath, manager); err != nil {
				logger.Failed("statehound notify socket error", err)
			}
		}()
	}

	if err := startControlSocket(model.ControlSocketPath, manager); err != nil {
		logger.Failed("statehound socket error", err)
		return
	}
}

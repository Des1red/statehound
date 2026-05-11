package daemon

import (
	"time"

	"statehound/internal/logger"
	"statehound/internal/model"
	"statehound/internal/statehound"
)

func Run() {
	logger.Status("statehound daemon started")

	manager := statehound.NewManager(5 * time.Second)
	go manager.Run()

	if err := startSocket(model.SocketPath, manager); err != nil {
		logger.Failed("statehound socket error", err)
		return
	}
}

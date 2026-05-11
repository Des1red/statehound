package runtime

import (
	"fmt"
	"os"

	"statehound/internal/logger"
	"statehound/internal/model"
)

func Events() {
	data, err := os.ReadFile(model.EventPath)
	if err != nil {
		logger.Failed("failed to read statehound events", err)
		return
	}

	if len(data) == 0 {
		logger.Status("no statehound events recorded yet")
		return
	}

	fmt.Print(string(data))
}

func ClearEvents() {
	if err := os.WriteFile(model.EventPath, []byte{}, 0600); err != nil {
		logger.Failed("failed to clear statehound events", err)
		return
	}

	logger.Success("statehound events cleared")
}

package runtime

import (
	"fmt"
	"os"
	"strings"

	"statehound/internal/logger"
	"statehound/internal/model"
)

func ShowEvents(filter string) {
	data, err := os.ReadFile(model.EventPath)
	if err != nil {
		logger.Failed("failed to read statehound events", err)
		return
	}

	if len(data) == 0 {
		logger.Status("no statehound events recorded yet")
		return
	}

	if filter == "" {
		fmt.Print(string(data))
		return
	}

	tag := "[" + filter + "]"
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, tag) {
			fmt.Println(line)
		}
	}
}

func ClearEvents() {
	if err := os.WriteFile(model.EventPath, []byte{}, 0600); err != nil {
		logger.Failed("failed to clear statehound events", err)
		return
	}

	logger.Success("statehound events cleared")
}

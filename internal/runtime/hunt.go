package runtime

import (
	"statehound/internal/client"
	"statehound/internal/logger"
)

func Hunt(target string) {
	resp, err := client.Send("HUNT " + target)
	if err != nil {
		logger.Failed("failed to hunt target", err)
		return
	}

	if resp == "" {
		logger.Warn("no hunt result returned")
		return
	}

	logger.Success(resp)
}

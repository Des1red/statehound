package runtime

import (
	"statehound/internal/client"
	"statehound/internal/logger"
	"strings"
)

func Profile() {
	resp, err := client.Send("PROFILE")
	if err != nil {
		logger.Failed("failed to start profiling", err)
		return
	}
	logger.Success(resp)

	url := ""
	for _, part := range strings.Fields(resp) {
		if strings.HasPrefix(part, "http://") {
			url = part
			break
		}
	}

	if url != "" {
		logger.Status("heap:       go tool pprof " + url + "heap")
		logger.Status("cpu:        go tool pprof " + url + "profile?seconds=30")
		logger.Status("goroutines: go tool pprof " + url + "goroutine")
	}
}

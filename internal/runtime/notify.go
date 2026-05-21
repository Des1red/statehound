package runtime

import (
	"encoding/json"
	"os"
	"time"

	"statehound/internal/client"
	"statehound/internal/logger"
	"statehound/internal/statehound"
	"statehound/internal/system"
	"statehound/internal/viewer"
)

func Notify() {
	if system.IsHeadless() {
		return
	}

	if !waitForNotifySocket(10 * time.Second) {
		logger.Status("notifier: notify socket did not become ready, exiting")
		return
	}

	exitWhenNotifySocketDies()

	logger.Status("statehound desktop notifier started")

	for {
		resp, err := client.SendNotify("NOTIFICATIONS")
		if err != nil {
			logger.Failed("failed to pull notifications", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var notifications []statehound.Notification
		if json.Unmarshal([]byte(resp), &notifications) == nil {
			viewer.HandleNotifications(notifications)
		}

		time.Sleep(2 * time.Second)
	}
}

func exitWhenNotifySocketDies() {
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if !client.IsNotifyRunning() {
				logger.Status("notifier: notify socket is gone, exiting")
				os.Exit(0)
			}
		}
	}()
}

func waitForNotifySocket(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if client.IsNotifyRunning() {
			return true
		}

		time.Sleep(500 * time.Millisecond)
	}

	return false
}

package runtime

import (
	"encoding/json"
	"time"

	"statehound/internal/client"
	"statehound/internal/logger"
	"statehound/internal/statehound"
	"statehound/internal/viewer"
)

func Notify() {
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

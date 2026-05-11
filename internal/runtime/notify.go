// internal/runtime/notify.go
package runtime

import (
	"encoding/json"
	"statehound/internal/client"
	"statehound/internal/logger"
	"statehound/internal/notify"
	"time"
)

func Notify() {
	logger.Status("statehound desktop notifier started")

	for {
		resp, err := client.SendNotify("NOTIFICATIONS")
		if err != nil {
			logger.Failed("failed to pull notifications", err)
			time.Sleep(2 * time.Second)
			continue
		} else {
			var notifications []notify.Notification
			if json.Unmarshal([]byte(resp), &notifications) == nil {
				for _, n := range notifications {
					if err := notify.Desktop(n.Title, n.Message, n.Urgency); err != nil {
						logger.Failed("failed to send desktop notification", err)
					}
				}
			}
		}

		time.Sleep(2 * time.Second)
	}
}

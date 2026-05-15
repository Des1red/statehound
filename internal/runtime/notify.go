package runtime

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"statehound/internal/client"
	"statehound/internal/logger"
	"statehound/internal/notify"
	"statehound/internal/viewer"
)

type urgencyState struct {
	mu      sync.Mutex
	count   int
	session *notify.Session
}

func Notify() {
	logger.Status("statehound desktop notifier started")

	states := map[string]*urgencyState{
		"critical": {},
		"normal":   {},
	}

	for {
		resp, err := client.SendNotify("NOTIFICATIONS")
		if err != nil {
			logger.Failed("failed to pull notifications", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var notifications []notify.Notification
		if json.Unmarshal([]byte(resp), &notifications) != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		counts := map[string]int{}
		for _, n := range notifications {
			if n.Urgency != "" {
				counts[n.Urgency]++
			}
		}

		for urgency, count := range counts {
			state := states[urgency]
			state.mu.Lock()

			state.count += count

			title := fmt.Sprintf("%d %s event", state.count, urgency)
			if state.count > 1 {
				title += "s"
			}

			if state.session == nil {
				session, err := notify.Start(title, "Click to view", urgency)
				if err == nil {
					state.session = session
					u := urgency
					s := state
					go func() {
						actionFired := <-session.ActionCh()
						logger.Status(fmt.Sprintf("notification action fired: %v urgency: %s", actionFired, u))
						if actionFired {
							logger.Status("launching viewer for " + u)
							go viewer.Show(u)
						}
						s.mu.Lock()
						s.session = nil
						s.count = 0
						s.mu.Unlock()
					}()
				}
			}

			state.mu.Unlock()
		}

		time.Sleep(2 * time.Second)
	}
}

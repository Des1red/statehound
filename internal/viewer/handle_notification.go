package viewer

import (
	"fmt"
	"statehound/internal/statehound"
)

func HandleNotifications(notifications []statehound.Notification) {
	counts := map[string]int{}
	for _, n := range notifications {
		if n.Urgency != "" {
			counts[n.Urgency]++
		}
	}

	if IsOpen() {
		return
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
			session, err := start(title, "Click to view", urgency)
			if err == nil {
				state.session = session
				u := urgency
				s := state
				go func() {
					if <-session.ActionCh() {
						s.mu.Lock()
						s.count = 0
						s.mu.Unlock()
						go Show(u)
					}
					s.mu.Lock()
					s.session = nil
					s.mu.Unlock()
				}()
			}
		}

		state.mu.Unlock()
	}
}

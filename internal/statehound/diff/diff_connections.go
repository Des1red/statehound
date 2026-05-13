package diff

import (
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/signals"
)

func diffConnections(previous, current map[string]collector.Connection) []events.Event {
	out := []events.Event{}

	for key, conn := range current {
		if _, existed := previous[key]; !existed {
			event := events.Event{
				Time:       time.Now(),
				Type:       signals.ConnectionOpened,
				Message:    formatter.FormatConnection(conn),
				Connection: &conn,
			}

			signals.TagConnectionEvent(&event, conn)

			out = append(out, event)
		}
	}

	for key, conn := range previous {
		if _, exists := current[key]; !exists {
			event := events.Event{
				Time:       time.Now(),
				Type:       signals.ConnectionClosed,
				Message:    formatter.FormatConnection(conn),
				Connection: &conn,
			}

			signals.TagConnectionEvent(&event, conn)

			out = append(out, event)
		}
	}

	return out
}

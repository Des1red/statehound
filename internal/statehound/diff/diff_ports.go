package diff

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/signals"
	"time"
)

func diffListeningPorts(previous, current map[string]collector.Port) []events.Event {
	out := []events.Event{}

	for key, port := range current {
		if _, existed := previous[key]; !existed {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.PortOpened,
				Message: formatter.FormatPort(port),
				Port:    &port,
			}

			signals.TagPortEvent(&event, port)

			out = append(out, event)
		}
	}

	for key, port := range previous {
		if _, exists := current[key]; !exists {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.PortClosed,
				Message: formatter.FormatPort(port),
				Port:    &port,
			}

			signals.TagPortEvent(&event, port)

			out = append(out, event)
		}
	}

	return out
}

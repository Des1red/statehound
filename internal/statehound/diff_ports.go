package statehound

import (
	"time"
)

func diffListeningPorts(previous, current map[string]Port) []Event {
	events := []Event{}

	for key, port := range current {
		if _, existed := previous[key]; !existed {
			events = append(events, Event{
				Time:    time.Now(),
				Type:    "PORT_OPENED",
				Message: formatPort(port),
			})
		}
	}

	for key, port := range previous {
		if _, exists := current[key]; !exists {
			events = append(events, Event{
				Time:    time.Now(),
				Type:    "PORT_CLOSED",
				Message: formatPort(port),
			})
		}
	}

	return events
}

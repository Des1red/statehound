package statehound

import "time"

func diffSystemdServices(previous, current map[string]Service) []Event {
	events := []Event{}

	for name, currentService := range current {
		previousService, existed := previous[name]
		if !existed {
			events = append(events, Event{
				Time:    time.Now(),
				Type:    "UNIT_APPEARED",
				Message: formatService(currentService),
			})
			continue
		}

		if serviceStateChanged(previousService, currentService) {
			events = append(events, Event{
				Time:    time.Now(),
				Type:    serviceTransitionType(previousService, currentService),
				Message: formatServiceTransition(previousService, currentService),
			})
		}
	}

	for name, previousService := range previous {
		if _, exists := current[name]; !exists {
			events = append(events, Event{
				Time:    time.Now(),
				Type:    "UNIT_REMOVED",
				Message: formatService(previousService),
			})
		}
	}

	return events
}

func serviceStateChanged(previous, current Service) bool {
	return previous.LoadState != current.LoadState ||
		previous.ActiveState != current.ActiveState ||
		previous.SubState != current.SubState
}

func serviceTransitionType(previous, current Service) string {
	if current.ActiveState == "failed" || current.SubState == "failed" {
		return "SERVICE_FAILED"
	}

	if previous.ActiveState == "inactive" && current.ActiveState == "active" {
		return "SERVICE_STARTED"
	}

	if previous.ActiveState == "active" && current.ActiveState == "inactive" {
		return "SERVICE_STOPPED"
	}

	return "SERVICE_STATE_CHANGED"
}

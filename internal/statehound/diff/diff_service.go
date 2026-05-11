package diff

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/signals"
	"time"
)

func diffSystemdServices(previous, current map[string]collector.Service) []events.Event {
	out := []events.Event{}

	for name, currentService := range current {
		previousService, existed := previous[name]
		if !existed {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.UnitAppeared,
				Message: formatter.FormatService(currentService),
			}

			signals.TagServiceEvent(&event, currentService)

			out = append(out, event)
			continue
		}

		if serviceStateChanged(previousService, currentService) {
			event := events.Event{
				Time:    time.Now(),
				Type:    serviceTransitionType(previousService, currentService),
				Message: formatter.FormatServiceTransition(previousService, currentService),
			}

			signals.TagServiceEvent(&event, currentService)

			out = append(out, event)
		}
	}

	for name, previousService := range previous {
		if _, exists := current[name]; !exists {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.UnitRemoved,
				Message: formatter.FormatService(previousService),
			}

			signals.TagServiceEvent(&event, previousService)

			out = append(out, event)
		}
	}

	return out
}

func serviceStateChanged(previous, current collector.Service) bool {
	return previous.LoadState != current.LoadState ||
		previous.ActiveState != current.ActiveState ||
		previous.SubState != current.SubState
}

func serviceTransitionType(previous, current collector.Service) string {
	if current.ActiveState == "failed" || current.SubState == "failed" {
		return signals.ServiceFailed
	}

	if previous.ActiveState == "inactive" && current.ActiveState == "active" {
		return signals.ServiceStarted
	}

	if previous.ActiveState == "active" && current.ActiveState == "inactive" {
		return signals.ServiceStopped
	}

	return signals.ServiceStateChanged
}

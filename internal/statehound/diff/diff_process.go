package diff

import (
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/signals"
)

func diffSuspiciousProcesses(previous, current map[string]collector.Process) []events.Event {
	out := []events.Event{}

	for pid, currentProcess := range current {
		previousProcess, existed := previous[pid]
		if !existed {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.SuspiciousProcessStarted,
				Message: formatter.FormatProcess(currentProcess),
				Process: &currentProcess,
			}

			signals.TagProcessEvent(&event, currentProcess)

			out = append(out, event)
			continue
		}

		if suspiciousProcessChanged(previousProcess, currentProcess) {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.SuspiciousProcessChanged,
				Message: formatter.FormatProcess(currentProcess),
				Process: &currentProcess,
			}

			signals.TagProcessEvent(&event, currentProcess)

			out = append(out, event)
		}
	}

	for pid, previousProcess := range previous {
		if _, exists := current[pid]; !exists {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.SuspiciousProcessStopped,
				Message: formatter.FormatProcess(previousProcess),
				Process: &previousProcess,
			}

			signals.TagProcessEvent(&event, previousProcess)

			out = append(out, event)
		}
	}

	return out
}

func suspiciousProcessChanged(previous, current collector.Process) bool {
	return previous.Exe != current.Exe ||
		previous.Cmdline != current.Cmdline ||
		previous.Reason != current.Reason
}

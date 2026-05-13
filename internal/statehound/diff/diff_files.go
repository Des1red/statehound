package diff

import (
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/signals"
)

func diffWatchedFiles(previous, current map[string]collector.FileWatch) []events.Event {
	out := []events.Event{}

	for path, currentFile := range current {
		previousFile, existed := previous[path]
		if !existed || !previousFile.Exists && currentFile.Exists {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.FileAdded,
				Message: formatter.FormatFile(currentFile),
				File:    &currentFile,
			}

			signals.TagFileEvent(&event, currentFile)

			out = append(out, event)
			continue
		}

		if fileChanged(previousFile, currentFile) {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.FileChanged,
				Message: formatter.FormatFileTransition(previousFile, currentFile),
				File:    &currentFile,
			}

			signals.TagFileEvent(&event, currentFile)

			out = append(out, event)
		}
	}

	for path, previousFile := range previous {
		currentFile, exists := current[path]
		if !exists || previousFile.Exists && !currentFile.Exists {
			event := events.Event{
				Time:    time.Now(),
				Type:    signals.FileRemoved,
				Message: formatter.FormatFile(previousFile),
				File:    &previousFile,
			}

			signals.TagFileEvent(&event, previousFile)

			out = append(out, event)
		}
	}

	return out
}

func fileChanged(previous, current collector.FileWatch) bool {
	if !previous.Exists && !current.Exists {
		return false
	}

	if previous.Exists != current.Exists {
		return true
	}

	return previous.Hash != current.Hash ||
		previous.Size != current.Size ||
		previous.Mode != current.Mode
}

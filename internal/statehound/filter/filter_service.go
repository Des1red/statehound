package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// ServiceFilter interface
type ServiceFilter interface {
	MatchService(e events.Event) bool
}

// IgnoreSafeServices drops events for known safe services
type IgnoreSafeServices struct {
	Safe []string
}

func (f *IgnoreSafeServices) MatchService(e events.Event) bool {
	for _, name := range f.Safe {
		if e.HasTag(signals.TagSystemUnit) && containsServiceName(e.Message, name) {
			return false
		}
	}
	return true
}

func containsServiceName(msg, name string) bool {
	// crude check: can be improved with regex
	return len(msg) >= len(name) && msg[:len(name)] == name
}

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
type IgnoreSafeServices struct{}

func (f *IgnoreSafeServices) MatchService(e events.Event) bool {
	return !e.HasTag(signals.TagNoiseUnit)
}

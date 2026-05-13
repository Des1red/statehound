package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// ProcessFilter interface
type ProcessFilter interface {
	MatchProcess(e events.Event) bool
}

// SuspiciousProcessFilter drops unclassified processes, keeps suspicious ones
type SuspiciousProcessFilter struct{}

func (f *SuspiciousProcessFilter) MatchProcess(e events.Event) bool {
	return !e.HasTag(signals.TagNoiseProcess)
}

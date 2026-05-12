package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// PortFilter interface
type PortFilter interface {
	MatchPort(e events.Event) bool
}

// SuspiciousListenerFilter flags public shell tools
type SuspiciousListenerFilter struct{}

func (f *SuspiciousListenerFilter) MatchPort(e events.Event) bool {
	if e.HasTag(signals.TagPublicListener) && e.HasTag(signals.TagShellTool) {
		return true
	}
	return false
}

package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// ProcessFilter interface
type ProcessFilter interface {
	MatchProcess(e events.Event) bool
}

// SuspiciousProcessFilter flags suspicious temp or go-build executables
type SuspiciousProcessFilter struct{}

func (f *SuspiciousProcessFilter) MatchProcess(e events.Event) bool {
	if e.HasTag(signals.TagSuspiciousProcess) ||
		e.HasTag(signals.TagDeletedExecutable) ||
		e.HasTag(signals.TagTempExecutable) ||
		e.HasTag(signals.TagGoBuildExecutable) {
		return true
	}
	return false
}

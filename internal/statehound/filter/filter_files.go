package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// FileFilter interface
type FileFilter interface {
	MatchFile(e events.Event) bool
}

// PersistenceFileFilter flags critical file changes
type PersistenceFileFilter struct{}

func (f *PersistenceFileFilter) MatchFile(e events.Event) bool {
	return e.HasTag(signals.TagPersistenceFile)
}

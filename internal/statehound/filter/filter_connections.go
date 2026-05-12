package filter

import (
	"statehound/internal/statehound/events"
	"time"
)

// ConnectionFilter interface
type ConnectionFilter interface {
	MatchConnection(e events.Event) bool
}

// ConnectionRateFilter limits logging of rapid repeated connections
type ConnectionRateFilter struct {
	LastSeen map[string]time.Time
	Interval time.Duration
}

func (f *ConnectionRateFilter) MatchConnection(e events.Event) bool {
	if e.Connection == nil {
		return true // not a connection event
	}

	key := e.Connection.RemoteAddress + ":" + e.Connection.LocalAddress
	now := e.Time
	if f.LastSeen == nil {
		f.LastSeen = make(map[string]time.Time)
	}

	last, ok := f.LastSeen[key]
	f.LastSeen[key] = now

	if ok && now.Sub(last) < f.Interval {
		return false
	}
	return true
}

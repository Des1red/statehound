package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
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
		return true
	}

	if e.HasTag(signals.TagNoiseConnection) {
		return false
	}

	key := e.Connection.LocalAddress + ":" + e.Connection.LocalPort + ":" +
		e.Connection.RemoteAddress + ":" + e.Connection.RemotePort + ":" +
		e.Connection.PID

	now := e.Time

	if f.LastSeen == nil {
		f.LastSeen = make(map[string]time.Time)
	}

	for k, t := range f.LastSeen {
		if now.Sub(t) > f.Interval {
			delete(f.LastSeen, k)
		}
	}

	last, ok := f.LastSeen[key]
	f.LastSeen[key] = now

	if ok && now.Sub(last) < f.Interval {
		return false
	}

	return true
}

package filter

import (
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

// DiffFilter holds filters for each diff type
type DiffFilter struct {
	ServiceFilters    []ServiceFilter
	PortFilters       []PortFilter
	FileFilters       []FileFilter
	ConnectionFilters []ConnectionFilter
	ProcessFilters    []ProcessFilter
}

// Result holds filtered events and notifications
type Result struct {
	LogEvents     []events.Event
	Notifications []events.Event
}

func (df *DiffFilter) Filter(evts []events.Event) Result {
	res := Result{}

	for _, e := range evts {
		keep := true

		switch e.Type {
		// systemd/service events
		case signals.UnitAppeared,
			signals.UnitRemoved,
			signals.ServiceStarted,
			signals.ServiceStopped,
			signals.ServiceFailed,
			signals.ServiceStateChanged:
			for _, f := range df.ServiceFilters {
				if !f.MatchService(e) {
					keep = false
					break
				}
			}

		// port events
		case signals.PortOpened,
			signals.PortClosed:
			for _, f := range df.PortFilters {
				if !f.MatchPort(e) {
					keep = false
					break
				}
			}

		// file events
		case signals.FileAdded,
			signals.FileRemoved,
			signals.FileChanged:
			for _, f := range df.FileFilters {
				if !f.MatchFile(e) {
					keep = false
					break
				}
			}

		// outbound/inbound connections
		case signals.ConnectionOpened,
			signals.ConnectionClosed:
			for _, f := range df.ConnectionFilters {
				if !f.MatchConnection(e) {
					keep = false
					break
				}
			}

		// suspicious process events
		case signals.SuspiciousProcessStarted,
			signals.SuspiciousProcessStopped,
			signals.SuspiciousProcessChanged:
			for _, f := range df.ProcessFilters {
				if !f.MatchProcess(e) {
					keep = false
					break
				}
			}
		}

		if keep {
			res.LogEvents = append(res.LogEvents, e)

			// Use tags from signals instead of hard-coded strings
			if e.HasTag(signals.TagPersistenceFile) ||
				e.HasTag(signals.TagSuspiciousProcess) ||
				e.HasTag(signals.TagCronFile) {
				res.Notifications = append(res.Notifications, e)
			}
		}
	}

	return res
}

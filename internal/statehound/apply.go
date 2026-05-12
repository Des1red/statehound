package statehound

import (
	"statehound/internal/notify"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/filter"
	"time"
)

// applyFilters applies all your diff filters and returns filtered events + notifications
func applyFilters(evts []events.Event) filter.Result {
	df := filter.DiffFilter{
		ServiceFilters: []filter.ServiceFilter{
			&filter.IgnoreSafeServices{Safe: []string{"systemd", "NetworkManager", "cupsd"}},
		},
		PortFilters: []filter.PortFilter{
			&filter.SuspiciousListenerFilter{},
		},
		FileFilters: []filter.FileFilter{
			&filter.PersistenceFileFilter{},
		},
		ConnectionFilters: []filter.ConnectionFilter{
			&filter.ConnectionRateFilter{Interval: 2 * time.Second},
		},
		ProcessFilters: []filter.ProcessFilter{
			&filter.SuspiciousProcessFilter{},
		},
	}

	return df.Filter(evts)
}

func eventsToNotifications(evts []events.Event) []notify.Notification {
	var out []notify.Notification
	for _, e := range evts {
		out = append(out, notify.Notification{
			Time:    e.Time.Format("2006-01-02T15:04:05Z07:00"),
			Title:   e.Type,
			Message: e.Message,
			Urgency: "critical", // or pick based on tag/type
		})
	}
	return out
}

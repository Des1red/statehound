package statehound

import (
	"statehound/internal/logger"
	"statehound/internal/notify"
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/diff"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
	"strconv"
	"time"
)

func (m *Manager) tick() {
	current, err := collector.CollectSnapshot()
	if err != nil {
		logger.Failed("failed to collect snapshot", err)
		return
	}

	m.mu.RLock()
	isBaseline := m.previous == nil
	m.mu.RUnlock()

	if isBaseline {
		m.mu.Lock()
		m.previous = &current
		m.lastScan = current.Time
		m.mu.Unlock()

		activeServices := collector.CountActiveServices(current.Services)
		evts := []events.Event{
			{
				Time: time.Now(),
				Type: signals.BaselineCreated,
				Message: "baseline created with systemd_services=" +
					strconv.Itoa(len(current.Services)) +
					" active_services=" +
					strconv.Itoa(activeServices) +
					" listening_ports=" +
					strconv.Itoa(len(current.Ports)) +
					" watched_files=" +
					strconv.Itoa(len(current.Files)),
			},
		}

		if err := events.WriteEvents(evts); err != nil {
			logger.Failed("failed to write baseline event", err)
		}

		return
	}

	m.mu.RLock()
	previous := *m.previous
	m.mu.RUnlock()
	evts := diff.DiffSnapshots(previous, current)

	if err := events.WriteEvents(evts); err != nil {
		logger.Failed("failed to write events", err)
	}

	// this stage should happen in filtering
	for _, event := range evts {
		if event.HasTag(signals.TagPersistenceFile) {
			m.PushNotification(notify.Notification{
				Time:    event.Time.Format(time.RFC3339),
				Title:   "Statehound persistence change",
				Message: event.Message,
				Urgency: "critical",
			})
		}
	}
	////////////////////
	m.mu.Lock()
	m.lastScan = current.Time
	m.previous = &current
	m.mu.Unlock()
}

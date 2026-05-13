package statehound

import (
	"fmt"
	"statehound/internal/logger"
	"statehound/internal/notify"
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/filter"
	"statehound/internal/statehound/signals"
	"sync"
	"time"
)

type Manager struct {
	interval time.Duration

	mu       sync.RWMutex
	previous *collector.Snapshot
	lastScan time.Time

	notifications []notify.Notification
	filter        filter.DiffFilter
}

func NewManager(interval time.Duration) *Manager {
	return &Manager{
		interval: interval,
		filter:   buildFilter(),
	}
}

func (m *Manager) Run() {
	logger.Status("statehound manager started")

	m.writeStartupEvent()
	m.loop()
}

func (m *Manager) writeStartupEvent() {
	evts := []events.Event{
		{
			Time:    time.Now(),
			Type:    signals.ManagerStarted,
			Message: "statehound manager started",
		},
	}

	if err := events.WriteEvents(evts); err != nil {
		logger.Failed("failed to write event", err)
	}
}

func (m *Manager) loop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.tick()

	for {
		<-ticker.C
		m.tick()
	}
}

func (m *Manager) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.previous == nil {
		return "statehound daemon is running\nmanager=starting"
	}

	return fmt.Sprintf(
		"statehound daemon is running\nmanager=running\ninterval=%s\nlast_scan=%s\nsystemd_services=%d\nactive_services=%d\nlistening_ports=%d\nwatched_files=%d\noutbound_connections=%d\nsuspicious_processes=%d",
		m.interval,
		m.lastScan.Format(time.RFC3339),
		len(m.previous.Services),
		collector.CountActiveServices(m.previous.Services),
		len(m.previous.Ports),
		len(m.previous.Files),
		len(m.previous.Connections),
		len(m.previous.Processes),
	)
}

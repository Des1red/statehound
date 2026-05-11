package statehound

import (
	"fmt"
	"statehound/internal/logger"
	"strconv"
	"sync"
	"time"
)

type Manager struct {
	interval time.Duration

	mu       sync.RWMutex
	previous *Snapshot
	lastScan time.Time
}

func NewManager(interval time.Duration) *Manager {
	return &Manager{
		interval: interval,
	}
}

func (m *Manager) Run() {
	logger.Status("statehound manager started")

	m.writeStartupEvent()
	m.loop()
}

func (m *Manager) writeStartupEvent() {
	events := []Event{
		{
			Time:    time.Now(),
			Type:    "MANAGER_STARTED",
			Message: "statehound manager started",
		},
	}

	if err := writeEvents(events); err != nil {
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

func (m *Manager) tick() {
	current, err := collectSnapshot()
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

		activeServices := countActiveServices(current.Services)
		events := []Event{
			{
				Time: time.Now(),
				Type: "BASELINE_CREATED",
				Message: "baseline created with systemd_services=" +
					strconv.Itoa(len(current.Services)) +
					" active_services=" +
					strconv.Itoa(activeServices) +
					" listening_ports=" +
					strconv.Itoa(len(current.Ports)),
			},
		}

		if err := writeEvents(events); err != nil {
			logger.Failed("failed to write baseline event", err)
		}

		return
	}

	m.mu.RLock()
	previous := *m.previous
	m.mu.RUnlock()
	events := diffSnapshots(previous, current)

	if err := writeEvents(events); err != nil {
		logger.Failed("failed to write events", err)
	}

	m.mu.Lock()
	m.lastScan = current.Time
	m.previous = &current
	m.mu.Unlock()
}

func (m *Manager) Status() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.previous == nil {
		return "statehound daemon is running\nmanager=starting"
	}

	return fmt.Sprintf(
		"statehound daemon is running\nmanager=running\ninterval=%s\nlast_scan=%s\nsystemd_services=%d\nactive_services=%d\nlistening_ports=%d",
		m.interval,
		m.lastScan.Format(time.RFC3339),
		len(m.previous.Services),
		countActiveServices(m.previous.Services),
		len(m.previous.Ports),
	)
}

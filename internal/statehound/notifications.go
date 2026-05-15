package statehound

import (
	"encoding/json"
	"statehound/internal/notify"
	"statehound/internal/statehound/events"
)

func (m *Manager) PushNotification(n notify.Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.notifications = append(m.notifications, n)
}

func (m *Manager) DrainNotifications() []notify.Notification {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := append([]notify.Notification(nil), m.notifications...)
	m.notifications = nil

	return out
}

func (m *Manager) NotificationsJSON() string {
	notifications := m.DrainNotifications()

	data, err := json.Marshal(notifications)
	if err != nil {
		return "[]"
	}

	return string(data)
}

func eventsToNotifications(evts []events.Event) []notify.Notification {
	var out []notify.Notification

	for _, e := range evts {
		if e.Urgency == "" {
			continue
		}

		out = append(out, notify.Notification{
			Time:    e.Time.Format("2006-01-02T15:04:05Z07:00"),
			Title:   e.Type,
			Message: e.Message,
			Urgency: e.Urgency,
		})
	}

	return out
}

package statehound

import (
	"encoding/json"
	"statehound/internal/statehound/events"
)

type Notification struct {
	Time    string
	Title   string
	Message string
	Urgency string
}

func (m *Manager) PushNotification(n Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.notifications = append(m.notifications, n)
}

func (m *Manager) DrainNotifications() []Notification {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := append([]Notification(nil), m.notifications...)
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

func eventsToNotifications(evts []events.Event) []Notification {
	var out []Notification

	for _, e := range evts {
		if e.Urgency == "" {
			continue
		}

		out = append(out, Notification{
			Time:    e.Time.Format("2006-01-02T15:04:05Z07:00"),
			Title:   e.Type,
			Message: e.Message,
			Urgency: e.Urgency,
		})
	}

	return out
}

package statehound

import (
	"encoding/json"
	"statehound/internal/notify"
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

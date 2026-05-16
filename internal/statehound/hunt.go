package statehound

import (
	"statehound/internal/statehound/hunt"
	"strings"
)

func (m *Manager) Hunt(target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return "hunt target is empty"
	}

	m.mu.RLock()
	if m.previous == nil {
		m.mu.RUnlock()
		return "manager snapshot is not ready yet"
	}

	services := hunt.CopyServices(m.previous.Services)
	ports := hunt.CopyPorts(m.previous.Ports)
	m.mu.RUnlock()

	targetLower := strings.ToLower(target)

	var out []string
	out = append(out, "hunt target: "+target)

	serviceMatches := hunt.HuntServices(services, targetLower)
	portMatches := hunt.HuntPorts(ports, targetLower)

	if len(serviceMatches) > 0 {
		out = append(out, "")
		out = append(out, "Services:")
		out = append(out, serviceMatches...)
	}

	if len(portMatches) > 0 {
		out = append(out, "")
		out = append(out, "Listening ports:")
		out = append(out, portMatches...)
	}

	if len(serviceMatches) == 0 && len(portMatches) == 0 {
		out = append(out, "no live matches found")
	}

	return strings.Join(out, "\n")
}

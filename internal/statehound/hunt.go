package statehound

import (
	"fmt"
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/formatter"
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

	services := copyServices(m.previous.Services)
	ports := copyPorts(m.previous.Ports)
	m.mu.RUnlock()

	targetLower := strings.ToLower(target)

	var out []string
	out = append(out, "hunt target: "+target)

	serviceMatches := huntServices(services, targetLower)
	portMatches := huntPorts(ports, targetLower)

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

func copyServices(src map[string]collector.Service) map[string]collector.Service {
	out := make(map[string]collector.Service, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

func copyPorts(src map[string]collector.Port) map[string]collector.Port {
	out := make(map[string]collector.Port, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

func huntServices(services map[string]collector.Service, target string) []string {
	var matches []string

	for name, service := range services {
		if serviceMatchesTarget(name, target) {
			line := "  " + formatter.FormatService(service)

			details := collector.CollectServiceDetails(name)
			if details.MainPID != "" && details.MainPID != "0" {
				line += " pid=" + details.MainPID
			}
			if details.FragmentPath != "" {
				line += " fragment=" + details.FragmentPath
			}

			matches = append(matches, line)
		}
	}

	return matches
}

func serviceMatchesTarget(name, target string) bool {
	name = strings.ToLower(name)
	target = strings.ToLower(target)

	return name == target ||
		name == target+".service" ||
		strings.TrimSuffix(name, ".service") == target
}

func huntPorts(ports map[string]collector.Port, target string) []string {
	var matches []string

	for _, port := range ports {
		if portMatchesTarget(port, target) {
			matches = append(matches, "  "+formatter.FormatPort(port))
		}
	}

	return matches
}

func portMatchesTarget(p collector.Port, target string) bool {
	values := []string{
		p.Proto,
		p.Address,
		p.Port,
		p.Scope,
		p.Process,
		p.PID,
		p.Exe,
		p.Cmdline,
		fmt.Sprintf("%s:%s", p.Address, p.Port),
	}

	for _, value := range values {
		if strings.Contains(strings.ToLower(value), target) {
			return true
		}
	}

	return false
}

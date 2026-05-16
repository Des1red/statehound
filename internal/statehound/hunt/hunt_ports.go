package hunt

import (
	"fmt"
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/formatter"
	"strings"
)

func HuntPorts(ports map[string]collector.Port, target string) []string {
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

func CopyPorts(src map[string]collector.Port) map[string]collector.Port {
	out := make(map[string]collector.Port, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

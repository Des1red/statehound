package hunt

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/formatter"
	"strings"
)

func HuntServices(services map[string]collector.Service, target string) []string {
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

func CopyServices(src map[string]collector.Service) map[string]collector.Service {
	out := make(map[string]collector.Service, len(src))

	for key, value := range src {
		out[key] = value
	}

	return out
}

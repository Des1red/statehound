package signals

import (
	"statehound/internal/statehound/collector"
	"strings"
)

var noiseServiceNames = []string{
	"systemd-journald",
	"systemd-logind",
	"systemd-udevd",
	"systemd-resolved",
	"systemd-timesyncd",
	"NetworkManager",
	"NetworkManager-dispatcher",
	"avahi-daemon",
	"ModemManager",
	"cupsd",
}

func isNoiseService(service collector.Service) bool {
	name := strings.TrimSuffix(service.Name, ".service")
	for _, noise := range noiseServiceNames {
		if name == strings.TrimSuffix(noise, ".service") {
			return true
		}
	}
	return false
}

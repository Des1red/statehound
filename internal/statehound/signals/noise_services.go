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
	"user-runtime-dir@",
	"user@",
}

func isNoiseService(service collector.Service) bool {
	name := strings.TrimSuffix(service.Name, ".service")

	for _, noise := range noiseServiceNames {
		noise = strings.TrimSuffix(noise, ".service")

		if strings.HasSuffix(noise, "@") {
			if strings.HasPrefix(name, noise) {
				return true
			}
			continue
		}

		if name == noise {
			return true
		}
	}

	return false
}

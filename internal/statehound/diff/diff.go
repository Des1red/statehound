package diff

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
)

func DiffSnapshots(previous, current collector.Snapshot) []events.Event {
	out := []events.Event{}

	out = append(out, diffSystemdServices(previous.Services, current.Services)...)
	out = append(out, diffListeningPorts(previous.Ports, current.Ports)...)

	return out
}

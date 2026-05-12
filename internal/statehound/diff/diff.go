package diff

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
)

func DiffSnapshots(previous, current collector.Snapshot) []events.Event {
	out := []events.Event{}

	out = append(out, diffSystemdServices(previous.Services, current.Services)...)
	out = append(out, diffListeningPorts(previous.Ports, current.Ports)...)
	out = append(out, diffWatchedFiles(previous.Files, current.Files)...)
	out = append(out, diffConnections(previous.Connections, current.Connections)...)
	out = append(out, diffSuspiciousProcesses(previous.SuspiciousProcesses, current.SuspiciousProcesses)...)

	return out
}

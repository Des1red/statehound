package statehound

import (
	"fmt"
	"strings"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/formatter"
	"statehound/internal/statehound/utils"
)

func (m *Manager) Snapshot() string {
	m.mu.RLock()
	if m.previous == nil {
		m.mu.RUnlock()
		return "snapshot is not ready yet"
	}

	snapshot := *m.previous
	lastScan := m.lastScan
	m.mu.RUnlock()

	var out []string

	out = append(out, "Statehound Snapshot")
	out = append(out, "")
	out = append(out, "Summary:")
	out = append(out, fmt.Sprintf("  last_scan=%s", lastScan.Format("2006-01-02T15:04:05Z07:00")))
	out = append(out, fmt.Sprintf("  systemd_services=%d", len(snapshot.Services)))
	out = append(out, fmt.Sprintf("  active_services=%d", collector.CountActiveServices(snapshot.Services)))
	out = append(out, fmt.Sprintf("  listening_ports=%d", len(snapshot.Ports)))
	out = append(out, fmt.Sprintf("  watched_files=%d", len(snapshot.Files)))
	out = append(out, fmt.Sprintf("  connections=%d", len(snapshot.Connections)))
	out = append(out, fmt.Sprintf("  processes=%d", len(snapshot.Processes)))

	out = append(out, "")
	out = append(out, "Active services:")
	activeServices := utils.SortedActiveServices(snapshot.Services)
	if len(activeServices) == 0 {
		out = append(out, "  none")
	} else {
		for _, service := range activeServices {
			out = append(out, "  "+formatter.FormatService(service))
		}
	}

	out = append(out, "")
	out = append(out, "Listening ports:")
	ports := utils.SortedPorts(snapshot.Ports)
	if len(ports) == 0 {
		out = append(out, "  none")
	} else {
		for _, port := range ports {
			out = append(out, "  "+formatter.FormatPort(port))
		}
	}

	out = append(out, "")
	out = append(out, "Watched files:")
	files := utils.SortedFiles(snapshot.Files)
	if len(files) == 0 {
		out = append(out, "  none")
	} else {
		for _, file := range files {
			out = append(out, "  "+formatter.FormatFile(file))
		}
	}

	out = append(out, "")
	out = append(out, "Connections:")
	connections := utils.SortedConnections(snapshot.Connections)
	if len(connections) == 0 {
		out = append(out, "  none")
	} else {
		for _, conn := range connections {
			out = append(out, "  "+formatter.FormatConnection(conn))
		}
	}

	out = append(out, "")
	out = append(out, "Processes:")
	processes := utils.SortedProcesses(snapshot.Processes)
	if len(processes) == 0 {
		out = append(out, "  none")
	} else {
		for _, proc := range processes {
			out = append(out, "  "+formatter.FormatProcess(proc))
		}
	}

	return strings.Join(out, "\n")
}

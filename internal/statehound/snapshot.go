package statehound

import (
	"fmt"
	"sort"
	"strings"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/formatter"
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
	activeServices := sortedActiveServices(snapshot.Services)
	if len(activeServices) == 0 {
		out = append(out, "  none")
	} else {
		for _, service := range activeServices {
			out = append(out, "  "+formatter.FormatService(service))
		}
	}

	out = append(out, "")
	out = append(out, "Listening ports:")
	ports := sortedPorts(snapshot.Ports)
	if len(ports) == 0 {
		out = append(out, "  none")
	} else {
		for _, port := range ports {
			out = append(out, "  "+formatter.FormatPort(port))
		}
	}

	out = append(out, "")
	out = append(out, "Watched files:")
	files := sortedFiles(snapshot.Files)
	if len(files) == 0 {
		out = append(out, "  none")
	} else {
		for _, file := range files {
			out = append(out, "  "+formatter.FormatFile(file))
		}
	}

	out = append(out, "")
	out = append(out, "Connections:")
	connections := sortedConnections(snapshot.Connections)
	if len(connections) == 0 {
		out = append(out, "  none")
	} else {
		for _, conn := range connections {
			out = append(out, "  "+formatter.FormatConnection(conn))
		}
	}

	out = append(out, "")
	out = append(out, "Processes:")
	processes := sortedProcesses(snapshot.Processes)
	if len(processes) == 0 {
		out = append(out, "  none")
	} else {
		for _, proc := range processes {
			out = append(out, "  "+formatter.FormatProcess(proc))
		}
	}

	return strings.Join(out, "\n")
}

func sortedActiveServices(services map[string]collector.Service) []collector.Service {
	out := make([]collector.Service, 0)

	for _, service := range services {
		if service.ActiveState == "active" {
			out = append(out, service)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})

	return out
}

func sortedPorts(ports map[string]collector.Port) []collector.Port {
	out := make([]collector.Port, 0, len(ports))

	for _, port := range ports {
		out = append(out, port)
	}

	sort.Slice(out, func(i, j int) bool {
		a := out[i].Proto + out[i].Address + out[i].Port
		b := out[j].Proto + out[j].Address + out[j].Port
		return a < b
	})

	return out
}

func sortedFiles(files map[string]collector.FileWatch) []collector.FileWatch {
	out := make([]collector.FileWatch, 0, len(files))

	for _, file := range files {
		out = append(out, file)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})

	return out
}

func sortedConnections(connections map[string]collector.Connection) []collector.Connection {
	out := make([]collector.Connection, 0, len(connections))

	for _, conn := range connections {
		out = append(out, conn)
	}

	sort.Slice(out, func(i, j int) bool {
		a := out[i].LocalAddress + out[i].LocalPort + out[i].RemoteAddress + out[i].RemotePort + out[i].PID
		b := out[j].LocalAddress + out[j].LocalPort + out[j].RemoteAddress + out[j].RemotePort + out[j].PID
		return a < b
	})

	return out
}

func sortedProcesses(processes map[string]collector.Process) []collector.Process {
	out := make([]collector.Process, 0, len(processes))

	for _, process := range processes {
		out = append(out, process)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].PID < out[j].PID
	})

	return out
}

package utils

import (
	"sort"
	"statehound/internal/statehound/collector"
)

func SortedActiveServices(services map[string]collector.Service) []collector.Service {
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

func SortedPorts(ports map[string]collector.Port) []collector.Port {
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

func SortedFiles(files map[string]collector.FileWatch) []collector.FileWatch {
	out := make([]collector.FileWatch, 0, len(files))

	for _, file := range files {
		out = append(out, file)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})

	return out
}

func SortedConnections(connections map[string]collector.Connection) []collector.Connection {
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

func SortedProcesses(processes map[string]collector.Process) []collector.Process {
	out := make([]collector.Process, 0, len(processes))

	for _, process := range processes {
		out = append(out, process)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].PID < out[j].PID
	})

	return out
}

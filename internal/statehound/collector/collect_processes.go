package collector

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func collectProcesses() (map[string]Process, error) {
	processes := make(map[string]Process)

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid := entry.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		proc := collectProcess(pid)
		if proc.PID == "" || proc.Exe == "" {
			continue
		}

		processes[pid] = proc
	}

	return processes, nil
}

func collectProcess(pid string) Process {
	status := readProcessStatus(pid)

	proc := Process{
		PID:     pid,
		PPID:    status["PPid"],
		UID:     firstStatusValue(status["Uid"]),
		Name:    status["Name"],
		Exe:     readProcessExe(pid),
		Cmdline: readProcessCmdline(pid),
	}

	if proc.Cmdline == "" {
		proc.Cmdline = proc.Name
	}

	return proc
}

func readProcessStatus(pid string) map[string]string {
	out := make(map[string]string)

	data, err := os.ReadFile(filepath.Join("/proc", pid, "status"))
	if err != nil {
		return out
	}

	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		out[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return out
}

func firstStatusValue(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

package collector

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func collectSuspiciousProcesses() (map[string]Process, error) {
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
		if proc.PID == "" {
			continue
		}

		if reason := suspiciousProcessReason(proc); reason != "" {
			proc.Reason = reason
			processes[pid] = proc
		}
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

func suspiciousProcessReason(p Process) string {
	exe := strings.ToLower(p.Exe)
	cmd := strings.ToLower(p.Cmdline)
	name := strings.ToLower(p.Name)

	switch {
	case strings.Contains(exe, " (deleted)"):
		return "deleted_executable"

	case strings.HasPrefix(exe, "/tmp/"),
		strings.HasPrefix(exe, "/var/tmp/"),
		strings.HasPrefix(exe, "/dev/shm/"):
		return "temp_executable"

	case strings.Contains(exe, "/.cache/go-build/"):
		return "go_build_cache_executable"

	case isSuspiciousShell(name, cmd):
		return "suspicious_shell"

	case isNetworkTool(name, exe, cmd):
		return "network_tool"

	case isScriptServer(name, cmd):
		return "script_server"

	default:
		return ""
	}
}

func isSuspiciousShell(name, cmd string) bool {
	if name != "bash" && name != "sh" && name != "zsh" {
		return false
	}

	return strings.Contains(cmd, "/dev/tcp") ||
		strings.Contains(cmd, "bash -i") ||
		strings.Contains(cmd, "sh -i") ||
		strings.Contains(cmd, "0>&1") ||
		strings.Contains(cmd, "2>&1")
}

func isNetworkTool(name, exe, cmd string) bool {
	values := []string{name, exe, cmd}
	needles := []string{"nc", "ncat", "netcat", "socat"}

	for _, value := range values {
		for _, needle := range needles {
			if value == needle ||
				strings.HasSuffix(value, "/"+needle) ||
				strings.Contains(value, " "+needle+" ") {
				return true
			}
		}
	}

	return false
}

func isScriptServer(name, cmd string) bool {
	switch name {
	case "python", "python3":
		return strings.Contains(cmd, "http.server") ||
			strings.Contains(cmd, "simplehttpserver")
	case "php":
		return strings.Contains(cmd, "-s ")
	case "ruby":
		return strings.Contains(cmd, "webrick")
	default:
		return false
	}
}

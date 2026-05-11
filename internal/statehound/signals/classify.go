package signals

import (
	"strings"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
)

func TagPortEvent(event *events.Event, port collector.Port) {
	switch port.Scope {
	case "public":
		event.Tag(TagPublicListener)
	case "local":
		event.Tag(TagLocalListener)
	case "interface":
		event.Tag(TagInterfaceListener)
	}

	switch port.Proto {
	case "tcp":
		event.Tag(TagTCPListener)
	case "udp":
		event.Tag(TagUDPListener)
	}

	process := strings.ToLower(port.Process)
	exe := strings.ToLower(port.Exe)
	cmd := strings.ToLower(port.Cmdline)

	if isShellTool(process, exe, cmd) {
		event.Tag(TagShellTool)
	}

	if isBrowserUDP(process, port.Proto) {
		event.Tag(TagBrowserUDP)
	}

	switch classifyExePath(exe) {
	case "system":
		event.Tag(TagSystemBinary)
	case "user":
		event.Tag(TagUserBinary)
	case "temp":
		event.Tag(TagTempBinary)
	default:
		event.Tag(TagUnknownBinary)
	}
}

func TagServiceEvent(event *events.Event, service collector.Service) {
	if service.ActiveState == "failed" || service.SubState == "failed" {
		event.Tag(TagServiceFailed)
	}

	if strings.HasPrefix(service.Name, "user@") ||
		strings.HasPrefix(service.Name, "user-runtime-dir@") {
		event.Tag(TagUserUnit)
		return
	}

	event.Tag(TagSystemUnit)
}

func isShellTool(process, exe, cmd string) bool {
	values := []string{process, exe, cmd}

	needles := []string{
		"nc",
		"ncat",
		"netcat",
		"socat",
		"bash",
		"sh",
		"zsh",
	}

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

func isBrowserUDP(process, proto string) bool {
	if proto != "udp" {
		return false
	}

	switch process {
	case "firefox", "brave", "chrome", "chromium":
		return true
	default:
		return false
	}
}

func classifyExePath(exe string) string {
	if exe == "" {
		return "unknown"
	}

	switch {
	case strings.HasPrefix(exe, "/usr/bin/"),
		strings.HasPrefix(exe, "/usr/sbin/"),
		strings.HasPrefix(exe, "/usr/lib"),
		strings.HasPrefix(exe, "/bin/"),
		strings.HasPrefix(exe, "/sbin/"),
		strings.HasPrefix(exe, "/opt/"):
		return "system"

	case strings.HasPrefix(exe, "/home/"):
		return "user"

	case strings.HasPrefix(exe, "/tmp/"),
		strings.HasPrefix(exe, "/var/tmp/"),
		strings.HasPrefix(exe, "/dev/shm/"):
		return "temp"

	default:
		return "unknown"
	}
}

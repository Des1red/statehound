package signals

import (
	"os"
	"statehound/internal/statehound/collector"
	"strings"
)

var standardRootPaths = []string{
	"/usr/bin/",
	"/usr/sbin/",
	"/usr/lib/",
	"/usr/lib64/",
	"/usr/libexec/",
	"/usr/local/",
	"/bin/",
	"/sbin/",
	"/lib/",
	"/lib64/",
	"/opt/",
	"/snap/",
}

var suspiciousParents = []string{
	"nginx", "apache2", "httpd", "lighttpd", "caddy",
	"php-fpm", "php", "gunicorn", "uwsgi",
	"mysqld", "mariadb", "postgres", "mongod",
	"redis-server", "node", "ruby",
}

var knownShells = []string{"bash", "sh", "zsh", "dash", "fish"}

func classifyProcess(p collector.Process) string {
	switch {
	case strings.Contains(p.Exe, " (deleted)"):
		return TagDeletedExecutable
	case strings.HasPrefix(p.Exe, "/tmp/"),
		strings.HasPrefix(p.Exe, "/var/tmp/"),
		strings.HasPrefix(p.Exe, "/dev/shm/"):
		return TagTempExecutable
	case isRootNonstandardPath(p):
		return TagRootNonstandardPath
	case isShellFromSuspiciousParent(p):
		return TagShellFromSuspiciousParent
	case isUserLocalExecutable(p):
		return TagUserLocalExecutable
	case isHomeHiddenExecutable(p):
		return TagHomeHiddenExecutable
	case isScriptServer(p):
		return TagScriptServer
	case isNetworkTool(p):
		return TagNetworkTool
	default:
		return ""
	}
}

func isHomeHiddenExecutable(p collector.Process) bool {
	if !strings.HasPrefix(p.Exe, "/home/") {
		return false
	}
	parts := strings.Split(p.Exe, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func isUserLocalExecutable(p collector.Process) bool {
	if !strings.HasPrefix(p.Exe, "/home/") {
		return false
	}
	return strings.Contains(p.Exe, "/.local/bin/")
}

func isScriptServer(p collector.Process) bool {
	name := strings.ToLower(p.Name)
	cmd := strings.ToLower(p.Cmdline)

	switch name {
	case "python", "python3":
		return strings.Contains(cmd, "http.server") ||
			strings.Contains(cmd, "simplehttpserver")
	case "php":
		return strings.Contains(cmd, "-s ")
	case "ruby":
		return strings.Contains(cmd, "webrick")
	}

	return false
}

func isNetworkTool(p collector.Process) bool {
	name := strings.ToLower(p.Name)
	exe := strings.ToLower(p.Exe)

	needles := []string{"nc", "ncat", "netcat", "socat"}

	for _, needle := range needles {
		if name == needle ||
			strings.HasSuffix(exe, "/"+needle) {
			return true
		}
	}

	return false
}

func isRootNonstandardPath(p collector.Process) bool {
	if p.UID != "0" || p.Exe == "" {
		return false
	}

	for _, path := range standardRootPaths {
		if strings.HasPrefix(p.Exe, path) {
			return false
		}
	}

	return true
}

func isShellFromSuspiciousParent(p collector.Process) bool {
	name := strings.ToLower(p.Name)

	isShell := false
	for _, shell := range knownShells {
		if name == shell {
			isShell = true
			break
		}
	}

	if !isShell {
		return false
	}

	if p.PPID == "" || p.PPID == "0" || p.PPID == "1" {
		return false
	}

	parentName := strings.ToLower(readParentName(p.PPID))
	if parentName == "" {
		return false
	}

	for _, parent := range suspiciousParents {
		if parentName == parent || strings.HasPrefix(parentName, parent) {
			return true
		}
	}

	return false
}

func readParentName(pid string) string {
	data, err := os.ReadFile("/proc/" + pid + "/comm")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
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

func classifyExePath(exe string) string {
	if exe == "" {
		return "unknown"
	}

	switch {
	case strings.HasPrefix(exe, "/usr/bin/"),
		strings.HasPrefix(exe, "/usr/sbin/"),
		strings.HasPrefix(exe, "/usr/lib/"),
		strings.HasPrefix(exe, "/usr/lib64/"),
		strings.HasPrefix(exe, "/usr/libexec/"),
		strings.HasPrefix(exe, "/bin/"),
		strings.HasPrefix(exe, "/sbin/"),
		strings.HasPrefix(exe, "/lib/"),
		strings.HasPrefix(exe, "/lib64/"),
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

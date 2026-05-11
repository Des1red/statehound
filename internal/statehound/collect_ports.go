package statehound

import (
	"regexp"
	"strings"

	"statehound/internal/command"
)

// collectListeningPorts collects currently listening TCP/UDP sockets.
//
// This catches exposed ports even when the owning process is not a systemd service,
// for example: python3 -m http.server, nc -lvnp, php -S, custom binaries, etc.
func collectListeningPorts() (map[string]Port, error) {
	out, err := command.Output("ss", "-H", "-ltnup")
	if err != nil {
		return nil, err
	}

	ports := make(map[string]Port)

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		port, ok := parseSSLine(line)
		if !ok {
			continue
		}

		key := portKey(port)
		ports[key] = port
	}

	return ports, nil
}

func portKey(p Port) string {
	return p.Proto + ":" + p.Address + ":" + p.Port
}

func parseSSLine(line string) (Port, bool) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return Port{}, false
	}

	proto := fields[0]
	local := fields[4]

	address, port, ok := splitAddressPort(local)
	if !ok {
		return Port{}, false
	}

	process, pid := parseProcessInfo(line)

	portInfo := Port{
		Proto:   proto,
		Address: address,
		Port:    port,
		Scope:   classifyPortScope(address),
		Process: process,
		PID:     pid,
	}

	return enrichPortProcess(portInfo), true
}

func splitAddressPort(local string) (string, string, bool) {
	idx := strings.LastIndex(local, ":")
	if idx == -1 || idx == len(local)-1 {
		return "", "", false
	}

	address := local[:idx]
	port := local[idx+1:]

	address = strings.Trim(address, "[]")

	return address, port, true
}

var procRe = regexp.MustCompile(`users:\(\("([^"]+)",pid=([0-9]+),`)

func parseProcessInfo(line string) (string, string) {
	match := procRe.FindStringSubmatch(line)
	if len(match) != 3 {
		return "", ""
	}

	return match[1], match[2]
}

func classifyPortScope(address string) string {
	switch address {
	case "127.0.0.1", "localhost", "::1":
		return "local"

	case "0.0.0.0", "*", "::":
		return "public"

	default:
		return "interface"
	}
}

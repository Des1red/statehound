package collector

import (
	"regexp"
	"strings"

	"statehound/internal/command"
)

// collectOutboundConnections collects established TCP connections.
//
// This is not alerting yet. It only snapshots active TCP connections so
// the manager can later diff or classify suspicious outbound activity.
func collectOutboundConnections() (map[string]Connection, error) {
	out, err := command.Output("ss", "-H", "-t", "-n", "-p", "-o", "state", "established")
	if err != nil {
		return nil, err
	}

	connections := make(map[string]Connection)

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		conn, ok := parseSSConnectionLine(line)
		if !ok {
			continue
		}

		key := connectionKey(conn)
		connections[key] = conn
	}

	return connections, nil
}

func connectionKey(c Connection) string {
	return c.Proto + ":" +
		c.LocalAddress + ":" + c.LocalPort + "->" +
		c.RemoteAddress + ":" + c.RemotePort + ":" +
		c.PID
}

func parseSSConnectionLine(line string) (Connection, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return Connection{}, false
	}

	// ss -H -t -n -p state established usually gives:
	// Recv-Q Send-Q Local Address:Port Peer Address:Port Process
	local := fields[2]
	remote := fields[3]

	localAddress, localPort, ok := splitAddressPort(local)
	if !ok {
		return Connection{}, false
	}

	remoteAddress, remotePort, ok := splitAddressPort(remote)
	if !ok {
		return Connection{}, false
	}

	process, pid := parseConnectionProcessInfo(line)

	conn := Connection{
		Proto:         "tcp",
		State:         "established",
		LocalAddress:  localAddress,
		LocalPort:     localPort,
		RemoteAddress: remoteAddress,
		RemotePort:    remotePort,
		Process:       process,
		PID:           pid,
	}

	return enrichConnectionProcess(conn), true
}

var connProcRe = regexp.MustCompile(`users:\(\("([^"]+)",pid=([0-9]+),`)

func parseConnectionProcessInfo(line string) (string, string) {
	match := connProcRe.FindStringSubmatch(line)
	if len(match) != 3 {
		return "", ""
	}

	return match[1], match[2]
}

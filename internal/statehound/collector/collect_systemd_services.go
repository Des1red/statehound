package collector

import (
	"strings"

	"statehound/internal/command"
)

// collectSystemdServices collects loaded systemd service units and their current states.
//
// It tracks service state transitions such as active/running, active/exited,
// inactive/dead, failed/failed, activating/start, and deactivating/stop.
//
// It does not collect every process on the host.
// Manual processes are handled through the listening-port snapshot.
func collectSystemdServices() (map[string]Service, error) {
	out, err := command.Output(
		"systemctl",
		"list-units",
		"--type=service",
		"--all",
		"--no-pager",
		"--no-legend",
	)
	if err != nil {
		return nil, err
	}

	services := make(map[string]Service)

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		service, ok := parseSystemdServiceLine(line)
		if !ok {
			continue
		}

		services[service.Name] = service
	}

	return services, nil
}

func parseSystemdServiceLine(line string) (Service, bool) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return Service{}, false
	}

	name := fields[0]
	loadState := fields[1]
	activeState := fields[2]
	subState := fields[3]

	description := ""
	if len(fields) > 4 {
		description = strings.Join(fields[4:], " ")
	}

	return Service{
		Name:        name,
		LoadState:   loadState,
		ActiveState: activeState,
		SubState:    subState,
		Description: description,
	}, true
}

package collector

import "time"

func CollectSnapshot() (Snapshot, error) {
	services, err := collectSystemdServices()
	if err != nil {
		return Snapshot{}, err
	}

	ports, err := collectListeningPorts()
	if err != nil {
		return Snapshot{}, err
	}

	files, err := collectWatchedFiles()
	if err != nil {
		return Snapshot{}, err
	}

	connections, err := collectOutboundConnections()
	if err != nil {
		return Snapshot{}, err
	}

	processes, err := collectSuspiciousProcesses()
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Time:                time.Now(),
		Services:            services,
		Ports:               ports,
		Files:               files,
		Connections:         connections,
		SuspiciousProcesses: processes,
	}, nil
}

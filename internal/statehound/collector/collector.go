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

	return Snapshot{
		Time:     time.Now(),
		Services: services,
		Ports:    ports,
	}, nil
}

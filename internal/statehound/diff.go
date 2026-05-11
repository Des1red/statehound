package statehound

func diffSnapshots(previous, current Snapshot) []Event {
	events := []Event{}

	events = append(events, diffSystemdServices(previous.Services, current.Services)...)
	events = append(events, diffListeningPorts(previous.Ports, current.Ports)...)

	return events
}

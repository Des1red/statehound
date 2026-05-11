package statehound

import "time"

type Snapshot struct {
	Time     time.Time
	Services map[string]Service
	Ports    map[string]Port
}

type Service struct {
	Name        string
	LoadState   string
	ActiveState string
	SubState    string
	Description string
}

type Port struct {
	Proto   string
	Address string
	Port    string
	Scope   string

	Process string
	PID     string
	Exe     string
	Cmdline string
}

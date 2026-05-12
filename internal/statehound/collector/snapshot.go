package collector

import "time"

type Snapshot struct {
	Time     time.Time
	Services map[string]Service
	Ports    map[string]Port
	Files    map[string]FileWatch
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

type FileWatch struct {
	Path    string
	Exists  bool
	Size    int64
	Mode    string
	ModTime time.Time
	Hash    string
}

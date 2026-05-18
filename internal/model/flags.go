package model

type CLIFlags struct {
	Help    bool
	Version bool

	Install   bool
	Uninstall bool
	Purge     bool

	Start    bool
	Stop     bool
	Restart  bool
	Status   bool
	Logs     bool
	Snapshot bool

	Hunt string

	Events      bool
	EventsGUI   bool
	ClearEvents bool
	Filter      string

	Web bool

	Profile bool
	Verbose bool // unused
}

type RuntimeFlags struct {
	Daemon  bool
	Notify  bool
	Verbose bool
}

type Flags struct {
	CLI     CLIFlags
	Runtime RuntimeFlags
}

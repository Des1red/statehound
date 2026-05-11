package model

type CLIFlags struct {
	Help    bool
	Version bool

	Install   bool
	Uninstall bool
	Purge     bool

	Start   bool
	Stop    bool
	Restart bool
	Status  bool
	Logs    bool
	Events  bool

	Hunt string

	ClearEvents bool

	Verbose bool
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

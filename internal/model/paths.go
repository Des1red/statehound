package model

const ServiceName = "statehound.service"

const (
	usrDir     = "/usr/local/bin"
	BinaryPath = usrDir + "/statehound"
	AliasPath  = usrDir + "/shound"

	ServicePath = "/etc/systemd/system/statehound.service"
	ConfigDir   = "/etc/statehound"
	LogDir      = "/var/log/statehound"
	EventPath   = LogDir + "/events.log"

	SocketPath = "/run/statehound/statehound.sock"
)

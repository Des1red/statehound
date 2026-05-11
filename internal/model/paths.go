package model

const (
	GroupName = "statehound"

	usrDir     = "/usr/local/bin"
	BinaryPath = usrDir + "/" + GroupName
	AliasPath  = usrDir + "/shound"

	ServiceName         = GroupName + ".service"
	NotifierServiceName = GroupName + "-notifier.service"

	ServicePath         = "/etc/systemd/system/" + ServiceName
	NotifierServicePath = "/etc/systemd/user/" + NotifierServiceName

	ConfigDir      = "/etc/" + GroupName
	LogDir         = "/var/log/" + GroupName
	EventBackupDir = LogDir + "/backups"
	EventPath      = LogDir + "/events.log"

	RuntimeDir        = "/run/" + GroupName
	ControlSocketPath = RuntimeDir + "/" + GroupName + ".sock"
	NotifySocketPath  = RuntimeDir + "/notify.sock"
)
const (
	MaxEventLogSizeBytes  = 5 * 1024 * 1024 // 5 MB
	EventBackupMaxAgeDays = 30
)

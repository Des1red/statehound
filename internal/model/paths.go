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

const (
	SystemCronFile = "/etc/crontab"
	SystemCronDir  = "/etc/cron.d"

	SystemdSystemDir     = "/etc/systemd/system"
	SystemdUserGlobalDir = "/etc/systemd/user"

	XDGAutostartDir = "/etc/xdg/autostart"

	RootAuthorizedKeysGlob = "/root/.ssh/authorized_keys"
	HomeAuthorizedKeysGlob = "/home/*/.ssh/authorized_keys"

	HomeUserSystemdGlob = "/home/*/.config/systemd/user/*.service"
	HomeAutostartGlob   = "/home/*/.config/autostart/*.desktop"
	HomeBashrcGlob      = "/home/*/.bashrc"
	HomeProfileGlob     = "/home/*/.profile"
	HomeZshrcGlob       = "/home/*/.zshrc"
	RootBashrc          = "/root/.bashrc"
	RootProfile         = "/root/.profile"
)

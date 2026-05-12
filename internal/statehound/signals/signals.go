package signals

const (
	UnitAppeared = "UNIT_APPEARED"
	UnitRemoved  = "UNIT_REMOVED"

	ServiceStarted      = "SERVICE_STARTED"
	ServiceStopped      = "SERVICE_STOPPED"
	ServiceFailed       = "SERVICE_FAILED"
	ServiceStateChanged = "SERVICE_STATE_CHANGED"

	PortOpened       = "PORT_OPENED"
	PortClosed       = "PORT_CLOSED"
	PortOwnerChanged = "PORT_OWNER_CHANGED"

	ManagerStarted  = "MANAGER_STARTED"
	BaselineCreated = "BASELINE_CREATED"

	FileAdded   = "FILE_ADDED"
	FileRemoved = "FILE_REMOVED"
	FileChanged = "FILE_CHANGED"
)

const (
	TagPublicListener    = "public_listener"
	TagLocalListener     = "local_listener"
	TagInterfaceListener = "interface_listener"

	TagTCPListener = "tcp_listener"
	TagUDPListener = "udp_listener"

	TagShellTool     = "shell_tool"
	TagBrowserUDP    = "browser_udp"
	TagSystemBinary  = "system_binary"
	TagUserBinary    = "user_binary"
	TagTempBinary    = "temp_binary"
	TagUnknownBinary = "unknown_binary"

	TagServiceFailed = "service_failed"
	TagUserUnit      = "user_unit"
	TagSystemUnit    = "system_unit"

	TagCronFile        = "cron_file"
	TagSystemdUnitFile = "systemd_unit_file"
	TagSSHKeysFile     = "ssh_keys_file"
	TagPersistenceFile = "persistence_file"

	TagOutboundConnection = "outbound_connection"
	TagInboundConnection  = "inbound_connection"
	TagExternalRemote     = "external_remote"
	TagLocalRemote        = "local_remote"

	TagSuspiciousProcess = "suspicious_process"
	TagDeletedExecutable = "deleted_executable"
	TagTempExecutable    = "temp_executable"
	TagGoBuildExecutable = "go_build_cache_executable"
)

const (
	ConnectionOpened = "CONNECTION_OPENED"
	ConnectionClosed = "CONNECTION_CLOSED"
)

const (
	SuspiciousProcessStarted = "SUSPICIOUS_PROCESS_STARTED"
	SuspiciousProcessStopped = "SUSPICIOUS_PROCESS_STOPPED"
	SuspiciousProcessChanged = "SUSPICIOUS_PROCESS_CHANGED"
)

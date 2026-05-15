package signals

const (
	TagUrgencyCritical = "critical"
	TagUrgencyNormal   = "normal"
)

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

	TagSuspiciousProcess         = "suspicious_process"
	TagHomeHiddenExecutable      = "home_hidden_executable"
	TagUserLocalExecutable       = "user_local_executable"
	TagDeletedExecutable         = "deleted_executable"
	TagTempExecutable            = "temp_executable"
	TagRootNonstandardPath       = "root_nonstandard_path"
	TagShellFromSuspiciousParent = "shell_from_suspicious_parent"
	TagScriptServer              = "script_server"
	TagNetworkTool               = "network_tool"
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

const (
	TagNoiseUnit       = "noise_unit"
	TagNoiseProcess    = "noise_process"
	TagNoiseConnection = "noise_connection"
	TagNoisePort       = "noise_port"
)

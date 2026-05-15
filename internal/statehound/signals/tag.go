package signals

import (
	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"strings"
)

func TagFileEvent(event *events.Event, file collector.FileWatch) {
	path := strings.ToLower(file.Path)

	// every watched file is a persistence candidate
	event.Tag(TagPersistenceFile)

	switch {
	case strings.Contains(path, "/cron"):
		event.Tag(TagCronFile)
		event.Tag(TagPersistenceFile)

	case strings.Contains(path, "/systemd/system/") && strings.HasSuffix(path, ".service"):
		event.Tag(TagSystemdUnitFile)
		event.Tag(TagPersistenceFile)

	case strings.HasSuffix(path, "/authorized_keys"):
		event.Tag(TagSSHKeysFile)
		event.Tag(TagPersistenceFile)
	}

	if event.HasTag(TagPersistenceFile) {
		event.TagUrgency(TagUrgencyCritical)
	}
}

func TagConnectionEvent(event *events.Event, conn collector.Connection) {
	event.Tag(TagOutboundConnection)

	remote := strings.ToLower(conn.RemoteAddress)

	switch remote {
	case "127.0.0.1", "localhost", "::1":
		event.Tag(TagLocalRemote)
	default:
		event.Tag(TagExternalRemote)
	}

	if isNoiseConnection(conn) {
		event.Tag(TagNoiseConnection)
		return
	}

	process := strings.ToLower(conn.Process)
	exe := strings.ToLower(conn.Exe)
	cmd := strings.ToLower(conn.Cmdline)

	if isShellTool(process, exe, cmd) {
		event.Tag(TagShellTool)
	}

	switch classifyExePath(exe) {
	case "system":
		event.Tag(TagSystemBinary)
	case "user":
		event.Tag(TagUserBinary)
	case "temp":
		event.Tag(TagTempBinary)
	default:
		event.Tag(TagUnknownBinary)
	}

	if event.HasTag(TagShellTool) && event.HasTag(TagExternalRemote) {
		event.TagUrgency(TagUrgencyCritical)
	}
}

func TagProcessEvent(event *events.Event, process collector.Process) {
	reason := classifyProcess(process)

	if reason == "" {
		event.Tag(TagNoiseProcess)
		return
	}

	if event.Process != nil {
		event.Process.Reason = reason
	}

	event.Tag(TagSuspiciousProcess)

	switch reason {
	case TagDeletedExecutable:
		event.Tag(TagDeletedExecutable)
	case TagTempExecutable:
		event.Tag(TagTempExecutable)
	case TagRootNonstandardPath:
		event.Tag(TagRootNonstandardPath)
	case TagShellFromSuspiciousParent:
		event.Tag(TagShellFromSuspiciousParent)
	case TagScriptServer:
		event.Tag(TagScriptServer)
	case TagNetworkTool:
		event.Tag(TagNetworkTool)
	case TagHomeHiddenExecutable:
		event.Tag(TagHomeHiddenExecutable)
	case TagUserLocalExecutable:
		event.Tag(TagUserLocalExecutable)
	}
	// urgency
	switch reason {
	case TagDeletedExecutable, TagTempExecutable,
		TagHomeHiddenExecutable, TagRootNonstandardPath,
		TagShellFromSuspiciousParent:
		event.TagUrgency(TagUrgencyCritical)
	case TagUserLocalExecutable, TagScriptServer, TagNetworkTool:
		event.TagUrgency(TagUrgencyNormal)
	}
}

func TagPortEvent(event *events.Event, port collector.Port) {
	switch port.Scope {
	case "public":
		event.Tag(TagPublicListener)
	case "local":
		event.Tag(TagLocalListener)
	case "interface":
		event.Tag(TagInterfaceListener)
	}

	switch port.Proto {
	case "tcp":
		event.Tag(TagTCPListener)
	case "udp":
		event.Tag(TagUDPListener)
	}

	process := strings.ToLower(port.Process)
	exe := strings.ToLower(port.Exe)
	cmd := strings.ToLower(port.Cmdline)

	if isShellTool(process, exe, cmd) {
		event.Tag(TagShellTool)
	}

	if isNoisePort(port) {
		event.Tag(TagNoisePort)
		return
	}

	switch classifyExePath(exe) {
	case "system":
		event.Tag(TagSystemBinary)
	case "user":
		event.Tag(TagUserBinary)
	case "temp":
		event.Tag(TagTempBinary)
	default:
		event.Tag(TagUnknownBinary)
	}

	if event.HasTag(TagPublicListener) && event.HasTag(TagShellTool) {
		event.TagUrgency(TagUrgencyCritical)
	} else if event.HasTag(TagPublicListener) {
		event.TagUrgency(TagUrgencyNormal)
	}
}

func TagServiceEvent(event *events.Event, service collector.Service) {
	if service.ActiveState == "failed" || service.SubState == "failed" {
		event.Tag(TagServiceFailed)
	}

	if isNoiseService(service) {
		event.Tag(TagNoiseUnit)
	}

	if strings.HasPrefix(service.Name, "user@") ||
		strings.HasPrefix(service.Name, "user-runtime-dir@") {
		event.Tag(TagUserUnit)
		return
	}

	event.Tag(TagSystemUnit)

	if event.HasTag(TagServiceFailed) {
		event.TagUrgency(TagUrgencyNormal)
	}
}

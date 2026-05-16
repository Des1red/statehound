package unit_test

import (
	"testing"
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

func newProcessEvent(p collector.Process) events.Event {
	e := events.Event{
		Time:    time.Now(),
		Type:    signals.SuspiciousProcessStarted,
		Process: &p,
	}
	return e
}

func newServiceEvent(s collector.Service) events.Event {
	return events.Event{
		Time:    time.Now(),
		Type:    signals.UnitAppeared,
		Service: &s,
	}
}

func newPortEvent(p collector.Port) events.Event {
	return events.Event{
		Time: time.Now(),
		Type: signals.PortOpened,
		Port: &p,
	}
}

func newConnectionEvent(c collector.Connection) events.Event {
	return events.Event{
		Time:       time.Now(),
		Type:       signals.ConnectionOpened,
		Connection: &c,
	}
}

func newFileEvent(f collector.FileWatch) events.Event {
	return events.Event{
		Time: time.Now(),
		Type: signals.FileAdded,
		File: &f,
	}
}

// --- Process ---

func TestTagProcessEvent_DeletedExecutable(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/usr/bin/bash (deleted)",
		UID: "0",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagDeletedExecutable)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagProcessEvent_TempExecutable(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/tmp/evil",
		UID: "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagTempExecutable)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagProcessEvent_RootNonstandardPath(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/home/user/evil",
		UID: "0",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagRootNonstandardPath)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagProcessEvent_HomeHiddenExecutable(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/home/user/.config/evil",
		UID: "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagHomeHiddenExecutable)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagProcessEvent_UserLocalExecutable(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/home/user/.local/bin/tool",
		UID: "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagUserLocalExecutable)
	assertUrgency(t, e, signals.TagUrgencyNormal)
}

func TestTagProcessEvent_ScriptServer(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID:     "1234",
		Exe:     "/usr/bin/python3",
		Name:    "python3",
		Cmdline: "python3 -m http.server 8080",
		UID:     "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagScriptServer)
	assertUrgency(t, e, signals.TagUrgencyNormal)
}

func TestTagProcessEvent_NetworkTool(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID:  "1234",
		Exe:  "/usr/bin/nc",
		Name: "nc",
		UID:  "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagSuspiciousProcess)
	assertHasTag(t, e, signals.TagNetworkTool)
	assertUrgency(t, e, signals.TagUrgencyNormal)
}

func TestTagProcessEvent_NormalProcess(t *testing.T) {
	e := newProcessEvent(collector.Process{
		PID: "1234",
		Exe: "/usr/bin/ls",
		UID: "1000",
	})
	signals.TagProcessEvent(&e, *e.Process)

	assertHasTag(t, e, signals.TagNoiseProcess)
	assertNoTag(t, e, signals.TagSuspiciousProcess)
}

// --- Service ---

func TestTagServiceEvent_NoiseService(t *testing.T) {
	e := newServiceEvent(collector.Service{Name: "NetworkManager.service"})
	signals.TagServiceEvent(&e, *e.Service)
	assertHasTag(t, e, signals.TagNoiseUnit)
}

func TestTagServiceEvent_NoiseServicePrefix(t *testing.T) {
	e := newServiceEvent(collector.Service{Name: "user@1000.service"})
	signals.TagServiceEvent(&e, *e.Service)
	assertHasTag(t, e, signals.TagUserUnit)
}

func TestTagServiceEvent_FailedService(t *testing.T) {
	e := newServiceEvent(collector.Service{
		Name:        "myapp.service",
		ActiveState: "failed",
		SubState:    "failed",
	})
	signals.TagServiceEvent(&e, *e.Service)
	assertHasTag(t, e, signals.TagServiceFailed)
	assertUrgency(t, e, signals.TagUrgencyNormal)
}

func TestTagServiceEvent_SystemUnit(t *testing.T) {
	e := newServiceEvent(collector.Service{
		Name:        "sshd.service",
		ActiveState: "active",
	})
	signals.TagServiceEvent(&e, *e.Service)
	assertHasTag(t, e, signals.TagSystemUnit)
	assertNoTag(t, e, signals.TagNoiseUnit)
}

// --- Port ---

func TestTagPortEvent_PublicShellTool(t *testing.T) {
	e := newPortEvent(collector.Port{
		Proto:   "tcp",
		Scope:   "public",
		Process: "nc",
		Exe:     "/usr/bin/nc",
	})
	signals.TagPortEvent(&e, *e.Port)
	assertHasTag(t, e, signals.TagPublicListener)
	assertHasTag(t, e, signals.TagShellTool)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagPortEvent_PublicNonShell(t *testing.T) {
	e := newPortEvent(collector.Port{
		Proto:   "tcp",
		Scope:   "public",
		Process: "nginx",
		Exe:     "/usr/bin/nginx",
	})
	signals.TagPortEvent(&e, *e.Port)
	assertHasTag(t, e, signals.TagPublicListener)
	assertUrgency(t, e, signals.TagUrgencyNormal)
}

func TestTagPortEvent_BrowserUDP(t *testing.T) {
	e := newPortEvent(collector.Port{
		Proto:   "udp",
		Scope:   "public",
		Process: "firefox",
	})
	signals.TagPortEvent(&e, *e.Port)
	assertHasTag(t, e, signals.TagNoisePort)
}

// --- Connection ---

func TestTagConnectionEvent_BrowserNoise(t *testing.T) {
	e := newConnectionEvent(collector.Connection{
		Process:       "firefox",
		RemoteAddress: "1.2.3.4",
	})
	signals.TagConnectionEvent(&e, *e.Connection)
	assertHasTag(t, e, signals.TagNoiseConnection)
}

func TestTagConnectionEvent_ShellExternal(t *testing.T) {
	e := newConnectionEvent(collector.Connection{
		Process:       "bash",
		Exe:           "/usr/bin/bash",
		RemoteAddress: "1.2.3.4",
	})
	signals.TagConnectionEvent(&e, *e.Connection)
	assertHasTag(t, e, signals.TagShellTool)
	assertHasTag(t, e, signals.TagExternalRemote)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

// --- File ---

func TestTagFileEvent_CronFile(t *testing.T) {
	e := newFileEvent(collector.FileWatch{Path: "/etc/cron.d/backdoor"})
	signals.TagFileEvent(&e, *e.File)
	assertHasTag(t, e, signals.TagCronFile)
	assertHasTag(t, e, signals.TagPersistenceFile)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagFileEvent_AuthorizedKeys(t *testing.T) {
	e := newFileEvent(collector.FileWatch{Path: "/home/user/.ssh/authorized_keys"})
	signals.TagFileEvent(&e, *e.File)
	assertHasTag(t, e, signals.TagSSHKeysFile)
	assertHasTag(t, e, signals.TagPersistenceFile)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

func TestTagFileEvent_SystemdUnit(t *testing.T) {
	e := newFileEvent(collector.FileWatch{Path: "/etc/systemd/system/evil.service"})
	signals.TagFileEvent(&e, *e.File)
	assertHasTag(t, e, signals.TagSystemdUnitFile)
	assertHasTag(t, e, signals.TagPersistenceFile)
	assertUrgency(t, e, signals.TagUrgencyCritical)
}

// --- Helpers ---

func assertHasTag(t *testing.T, e events.Event, tag string) {
	t.Helper()
	if !e.HasTag(tag) {
		t.Errorf("expected tag %q not found in %v", tag, e.Tags)
	}
}

func assertNoTag(t *testing.T, e events.Event, tag string) {
	t.Helper()
	if e.HasTag(tag) {
		t.Errorf("unexpected tag %q found in %v", tag, e.Tags)
	}
}

func assertUrgency(t *testing.T, e events.Event, urgency string) {
	t.Helper()
	if e.Urgency != urgency {
		t.Errorf("expected urgency %q got %q", urgency, e.Urgency)
	}
}

package unit_test

import (
	"testing"
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/diff"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/signals"
)

func emptySnapshot() collector.Snapshot {
	return collector.Snapshot{
		Time:        time.Now(),
		Services:    map[string]collector.Service{},
		Ports:       map[string]collector.Port{},
		Files:       map[string]collector.FileWatch{},
		Connections: map[string]collector.Connection{},
		Processes:   map[string]collector.Process{},
	}
}

func assertEventType(t *testing.T, evts []events.Event, eventType string) {
	t.Helper()
	for _, e := range evts {
		if e.Type == eventType {
			return
		}
	}
	t.Errorf("expected event type %q not found in %d events", eventType, len(evts))
}

func TestDiff_ServiceAppeared(t *testing.T) {
	prev := emptySnapshot()
	curr := emptySnapshot()
	curr.Services["evil.service"] = collector.Service{
		Name: "evil.service", ActiveState: "active", SubState: "running",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.UnitAppeared)
}

func TestDiff_ServiceRemoved(t *testing.T) {
	prev := emptySnapshot()
	prev.Services["evil.service"] = collector.Service{Name: "evil.service"}
	curr := emptySnapshot()
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.UnitRemoved)
}

func TestDiff_ServiceStarted(t *testing.T) {
	prev := emptySnapshot()
	prev.Services["myapp.service"] = collector.Service{
		Name: "myapp.service", ActiveState: "inactive", SubState: "dead",
	}
	curr := emptySnapshot()
	curr.Services["myapp.service"] = collector.Service{
		Name: "myapp.service", ActiveState: "active", SubState: "running",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.ServiceStarted)
}

func TestDiff_ServiceFailed(t *testing.T) {
	prev := emptySnapshot()
	prev.Services["myapp.service"] = collector.Service{
		Name: "myapp.service", ActiveState: "active", SubState: "running",
	}
	curr := emptySnapshot()
	curr.Services["myapp.service"] = collector.Service{
		Name: "myapp.service", ActiveState: "failed", SubState: "failed",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.ServiceFailed)
}

func TestDiff_PortOpened(t *testing.T) {
	prev := emptySnapshot()
	curr := emptySnapshot()
	curr.Ports["tcp:0.0.0.0:4444"] = collector.Port{
		Proto: "tcp", Address: "0.0.0.0", Port: "4444", Scope: "public",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.PortOpened)
}

func TestDiff_PortClosed(t *testing.T) {
	prev := emptySnapshot()
	prev.Ports["tcp:0.0.0.0:4444"] = collector.Port{
		Proto: "tcp", Address: "0.0.0.0", Port: "4444", Scope: "public",
	}
	curr := emptySnapshot()
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.PortClosed)
}

func TestDiff_FileAdded(t *testing.T) {
	prev := emptySnapshot()
	curr := emptySnapshot()
	curr.Files["/etc/cron.d/backdoor"] = collector.FileWatch{
		Path: "/etc/cron.d/backdoor", Exists: true, Hash: "abc123",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.FileAdded)
}

func TestDiff_FileRemoved(t *testing.T) {
	prev := emptySnapshot()
	prev.Files["/etc/cron.d/backdoor"] = collector.FileWatch{
		Path: "/etc/cron.d/backdoor", Exists: true, Hash: "abc123",
	}
	curr := emptySnapshot()
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.FileRemoved)
}

func TestDiff_FileChanged(t *testing.T) {
	prev := emptySnapshot()
	prev.Files["/root/.bashrc"] = collector.FileWatch{
		Path: "/root/.bashrc", Exists: true, Hash: "aaa",
	}
	curr := emptySnapshot()
	curr.Files["/root/.bashrc"] = collector.FileWatch{
		Path: "/root/.bashrc", Exists: true, Hash: "bbb",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.FileChanged)
}

func TestDiff_ConnectionOpened(t *testing.T) {
	prev := emptySnapshot()
	curr := emptySnapshot()
	curr.Connections["tcp:127.0.0.1:9000->1.2.3.4:443:1234"] = collector.Connection{
		Proto: "tcp", LocalAddress: "127.0.0.1", LocalPort: "9000",
		RemoteAddress: "1.2.3.4", RemotePort: "443", PID: "1234",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.ConnectionOpened)
}

func TestDiff_ConnectionClosed(t *testing.T) {
	prev := emptySnapshot()
	prev.Connections["tcp:127.0.0.1:9000->1.2.3.4:443:1234"] = collector.Connection{
		Proto: "tcp", LocalAddress: "127.0.0.1", LocalPort: "9000",
		RemoteAddress: "1.2.3.4", RemotePort: "443", PID: "1234",
	}
	curr := emptySnapshot()
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.ConnectionClosed)
}

func TestDiff_ProcessStarted(t *testing.T) {
	prev := emptySnapshot()
	curr := emptySnapshot()
	curr.Processes["1234"] = collector.Process{
		PID: "1234", Exe: "/tmp/evil", UID: "0",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.SuspiciousProcessStarted)
}

func TestDiff_ProcessStopped(t *testing.T) {
	prev := emptySnapshot()
	prev.Processes["1234"] = collector.Process{
		PID: "1234", Exe: "/tmp/evil", UID: "0",
	}
	curr := emptySnapshot()
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.SuspiciousProcessStopped)
}

func TestDiff_ProcessChanged(t *testing.T) {
	prev := emptySnapshot()
	prev.Processes["1234"] = collector.Process{
		PID: "1234", Exe: "/tmp/evil", UID: "0",
	}
	curr := emptySnapshot()
	curr.Processes["1234"] = collector.Process{
		PID: "1234", Exe: "/tmp/evil2", UID: "0",
	}
	assertEventType(t, diff.DiffSnapshots(prev, curr), signals.SuspiciousProcessChanged)
}

func TestDiff_NoChange(t *testing.T) {
	snap := emptySnapshot()
	snap.Services["sshd.service"] = collector.Service{
		Name: "sshd.service", ActiveState: "active", SubState: "running",
	}
	evts := diff.DiffSnapshots(snap, snap)
	if len(evts) != 0 {
		t.Errorf("expected no events for identical snapshots, got %d", len(evts))
	}
}

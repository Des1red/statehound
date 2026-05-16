package unit_test

import (
	"testing"
	"time"

	"statehound/internal/statehound/collector"
	"statehound/internal/statehound/events"
	"statehound/internal/statehound/filter"
	"statehound/internal/statehound/signals"
)

func buildTestFilter() filter.DiffFilter {
	return filter.BuildFilter()
}

func serviceEvent(name, activeState, subState string, tags ...string) events.Event {
	e := events.Event{
		Time: time.Now(),
		Type: signals.ServiceStarted,
		Service: &collector.Service{
			Name:        name,
			ActiveState: activeState,
			SubState:    subState,
		},
	}
	for _, tag := range tags {
		e.Tag(tag)
	}
	return e
}

func portEvent(scope, process string, tags ...string) events.Event {
	e := events.Event{
		Time: time.Now(),
		Type: signals.PortOpened,
		Port: &collector.Port{
			Scope:   scope,
			Process: process,
		},
	}
	for _, tag := range tags {
		e.Tag(tag)
	}
	return e
}

func fileEvent(path string, tags ...string) events.Event {
	e := events.Event{
		Time: time.Now(),
		Type: signals.FileAdded,
		File: &collector.FileWatch{Path: path},
	}
	for _, tag := range tags {
		e.Tag(tag)
	}
	return e
}

func connectionEvent(remote string, tags ...string) events.Event {
	e := events.Event{
		Time: time.Now(),
		Type: signals.ConnectionOpened,
		Connection: &collector.Connection{
			RemoteAddress: remote,
			RemotePort:    "443",
			LocalAddress:  "127.0.0.1",
			LocalPort:     "9000",
			PID:           "1234",
		},
	}
	e.Time = time.Now()
	for _, tag := range tags {
		e.Tag(tag)
	}
	return e
}

func processEvent(exe string, tags ...string) events.Event {
	e := events.Event{
		Time:    time.Now(),
		Type:    signals.SuspiciousProcessStarted,
		Process: &collector.Process{PID: "1234", Exe: exe},
	}
	for _, tag := range tags {
		e.Tag(tag)
	}
	return e
}

// --- Service filter ---

func TestFilter_NoiseServiceDropped(t *testing.T) {
	f := buildTestFilter()
	e := serviceEvent("NetworkManager.service", "active", "running", signals.TagNoiseUnit)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 0 {
		t.Errorf("expected noise service to be dropped, got %d events", len(res.LogEvents))
	}
}

func TestFilter_RealServiceKept(t *testing.T) {
	f := buildTestFilter()
	e := serviceEvent("evil.service", "active", "running", signals.TagSystemUnit)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected real service to be kept, got %d events", len(res.LogEvents))
	}
}

// --- Port filter ---

func TestFilter_NoisePortDropped(t *testing.T) {
	f := buildTestFilter()
	e := portEvent("public", "firefox", signals.TagNoisePort)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 0 {
		t.Errorf("expected noise port to be dropped, got %d events", len(res.LogEvents))
	}
}

func TestFilter_PublicListenerKept(t *testing.T) {
	f := buildTestFilter()
	e := portEvent("public", "nginx", signals.TagPublicListener)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected public listener to be kept, got %d events", len(res.LogEvents))
	}
}

func TestFilter_ShellToolKept(t *testing.T) {
	f := buildTestFilter()
	e := portEvent("local", "nc", signals.TagShellTool)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected shell tool to be kept, got %d events", len(res.LogEvents))
	}
}

// --- File filter ---

func TestFilter_PersistenceFileKept(t *testing.T) {
	f := buildTestFilter()
	e := fileEvent("/etc/cron.d/backdoor", signals.TagPersistenceFile, signals.TagCronFile)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected persistence file to be kept, got %d events", len(res.LogEvents))
	}
}

func TestFilter_NonPersistenceFileDropped(t *testing.T) {
	f := buildTestFilter()
	e := fileEvent("/tmp/something")
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 0 {
		t.Errorf("expected non-persistence file to be dropped, got %d events", len(res.LogEvents))
	}
}

// --- Connection filter ---

func TestFilter_NoiseConnectionDropped(t *testing.T) {
	f := buildTestFilter()
	e := connectionEvent("1.2.3.4", signals.TagNoiseConnection)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 0 {
		t.Errorf("expected noise connection to be dropped, got %d events", len(res.LogEvents))
	}
}

func TestFilter_RealConnectionKept(t *testing.T) {
	f := buildTestFilter()
	e := connectionEvent("1.2.3.4", signals.TagExternalRemote)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected real connection to be kept, got %d events", len(res.LogEvents))
	}
}

func TestFilter_RapidConnectionRateLimited(t *testing.T) {
	f := buildTestFilter()
	e := connectionEvent("1.2.3.4", signals.TagExternalRemote)

	res1 := f.Filter([]events.Event{e})
	e.Time = e.Time.Add(1 * time.Second)
	res2 := f.Filter([]events.Event{e})

	if len(res1.LogEvents) != 1 {
		t.Errorf("expected first connection to be kept")
	}
	if len(res2.LogEvents) != 0 {
		t.Errorf("expected rapid repeat connection to be rate limited")
	}
}

// --- Process filter ---

func TestFilter_NoiseProcessDropped(t *testing.T) {
	f := buildTestFilter()
	e := processEvent("/usr/bin/ls", signals.TagNoiseProcess)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 0 {
		t.Errorf("expected noise process to be dropped, got %d events", len(res.LogEvents))
	}
}

func TestFilter_SuspiciousProcessKept(t *testing.T) {
	f := buildTestFilter()
	e := processEvent("/tmp/evil", signals.TagSuspiciousProcess, signals.TagTempExecutable)
	res := f.Filter([]events.Event{e})
	if len(res.LogEvents) != 1 {
		t.Errorf("expected suspicious process to be kept, got %d events", len(res.LogEvents))
	}
}

// --- Notifications ---

func TestFilter_CriticalEventBecomesNotification(t *testing.T) {
	f := buildTestFilter()
	e := fileEvent("/etc/cron.d/backdoor", signals.TagPersistenceFile)
	e.Urgency = signals.TagUrgencyCritical
	res := f.Filter([]events.Event{e})
	if len(res.Notifications) != 1 {
		t.Errorf("expected critical event in notifications, got %d", len(res.Notifications))
	}
}

func TestFilter_NoUrgencyNoNotification(t *testing.T) {
	f := buildTestFilter()
	e := serviceEvent("myapp.service", "active", "running", signals.TagSystemUnit)
	res := f.Filter([]events.Event{e})
	if len(res.Notifications) != 0 {
		t.Errorf("expected no notification for event without urgency, got %d", len(res.Notifications))
	}
}

package viewer

import (
	"io"
	"os/exec"
	"sync"
)

type tabSession struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	lastLines int
}

type viewerSession struct {
	notebook *exec.Cmd
	tabs     map[string]*tabSession
	updateCh chan tabUpdate
}

type tabUpdate struct {
	urgency string
	lines   []string
}

type tabConfig struct {
	urgency string
	label   string
	tabnum  int
}

var (
	viewerMu sync.Mutex
	session  *viewerSession
)

var tabs = []tabConfig{
	{"critical", "Critical", 1},
	{"normal", "Normal", 2},
}

type urgencyState struct {
	mu      sync.Mutex
	count   int
	session *Session
}

var states = map[string]*urgencyState{
	"critical": {},
	"normal":   {},
}

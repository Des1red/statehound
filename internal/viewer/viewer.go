package viewer

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"statehound/internal/logger"
	"statehound/internal/statehound/events"
)

type viewerSession struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	lastLines int
}

var (
	sessions  = map[string]*viewerSession{}
	sessionMu sync.Mutex
)

func Show(urgency string) {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	content, err := events.FilterEvents(urgency)
	if err != nil {
		logger.Failed("viewer: failed to filter events", err)
		return
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	if content == "" {
		lines = []string{"no " + urgency + " events recorded"}
	}

	session, exists := sessions[urgency]

	if !exists {
		cmd := exec.Command(
			"yad",
			"--text-info",
			"--listen",
			"--tail",
			"--title=statehound — "+urgency+" events",
			"--width=900",
			"--height=600",
			"--fontname=monospace",
		)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			logger.Failed("viewer: failed to create stdin pipe", err)
			return
		}

		if err := cmd.Start(); err != nil {
			logger.Failed("viewer: failed to start yad", err)
			return
		}

		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				line = " " + line
			}
			fmt.Fprintln(stdin, line)
		}

		session = &viewerSession{
			cmd:       cmd,
			stdin:     stdin,
			lastLines: len(lines),
		}
		sessions[urgency] = session

		go func() {
			cmd.Wait()
			sessionMu.Lock()
			delete(sessions, urgency)
			sessionMu.Unlock()
		}()

	} else if len(lines) > session.lastLines {
		for _, line := range lines[session.lastLines:] {
			if strings.HasPrefix(line, "#") {
				line = " " + line
			}
			fmt.Fprintln(session.stdin, line)
		}
		session.lastLines = len(lines)
	}
}

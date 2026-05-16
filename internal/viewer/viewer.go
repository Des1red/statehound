package viewer

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"statehound/internal/logger"
	"statehound/internal/statehound/events"
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

var (
	viewerMu sync.Mutex
	session  *viewerSession
)

var tabs = []struct {
	urgency string
	label   string
	tabnum  int
}{
	{"critical", "Critical", 1},
	{"normal", "Normal", 2},
}

func writeLine(stdin io.WriteCloser, line string) {
	if strings.HasPrefix(line, "#") {
		line = " " + line
	}
	fmt.Fprintln(stdin, line)
}

func Show(urgency string) {
	if !acquireLock() {
		return
	}
	defer releaseLock()

	viewerMu.Lock()
	if session != nil {
		viewerMu.Unlock()
		return
	}

	key := rand.Intn(0x7FFFFFFF) + 1
	tabSessions := make(map[string]*tabSession)

	for _, t := range tabs {
		cmd := exec.Command(
			"yad",
			fmt.Sprintf("--plug=%d", key),
			fmt.Sprintf("--tabnum=%d", t.tabnum),
			"--text-info",
			"--listen",
			"--tail",
			"--fontname=monospace",
		)
		cmd.Env = buildEnv()

		stdin, err := cmd.StdinPipe()
		if err != nil {
			viewerMu.Unlock()
			killTabs(tabSessions)
			logger.Failed("viewer: failed to create stdin pipe for "+t.urgency, err)
			return
		}

		if err := cmd.Start(); err != nil {
			viewerMu.Unlock()
			killTabs(tabSessions)
			logger.Failed("viewer: failed to start tab for "+t.urgency, err)
			return
		}

		content, _ := events.FilterEvents(t.urgency)
		lines := strings.Split(strings.TrimSpace(content), "\n")
		if content == "" {
			lines = []string{"no " + t.urgency + " events recorded"}
		}

		for _, line := range lines {
			writeLine(stdin, line)
		}

		tabSessions[t.urgency] = &tabSession{
			cmd:       cmd,
			stdin:     stdin,
			lastLines: len(lines),
		}
	}

	time.Sleep(200 * time.Millisecond)

	activeTab := 1
	for _, t := range tabs {
		if t.urgency == urgency {
			activeTab = t.tabnum
			break
		}
	}

	notebook := exec.Command(
		"yad",
		"--notebook",
		fmt.Sprintf("--key=%d", key),
		"--tab=Critical",
		"--tab=Normal",
		fmt.Sprintf("--active-tab=%d", activeTab),
		"--width=900",
		"--height=600",
		"--title=statehound events",
		"--no-buttons",
	)
	notebook.Env = buildEnv()
	notebook.Stderr = os.Stderr

	if err := notebook.Start(); err != nil {
		viewerMu.Unlock()
		killTabs(tabSessions)
		logger.Failed("viewer: failed to start notebook", err)
		return
	}

	updateCh := make(chan tabUpdate, 32)
	done := make(chan struct{})

	session = &viewerSession{
		notebook: notebook,
		tabs:     tabSessions,
		updateCh: updateCh,
	}
	viewerMu.Unlock()

	// update writer
	go func() {
		for update := range updateCh {
			viewerMu.Lock()
			tab, ok := session.tabs[update.urgency]
			viewerMu.Unlock()
			if ok {
				for _, line := range update.lines {
					writeLine(tab.stdin, line)
				}
			}
		}
	}()

	// self-updating watcher — works regardless of which process opened the viewer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				viewerMu.Lock()
				if session == nil {
					viewerMu.Unlock()
					return
				}
				var updates []tabUpdate
				for _, t := range tabs {
					tab := session.tabs[t.urgency]
					if tab == nil {
						continue
					}
					content, _ := events.FilterEvents(t.urgency)
					lines := strings.Split(strings.TrimSpace(content), "\n")
					if len(lines) > tab.lastLines {
						newLines := make([]string, len(lines)-tab.lastLines)
						copy(newLines, lines[tab.lastLines:])
						tab.lastLines = len(lines)
						updates = append(updates, tabUpdate{urgency: t.urgency, lines: newLines})
					}
				}
				ch := session.updateCh
				viewerMu.Unlock()
				for _, u := range updates {
					select {
					case ch <- u:
					case <-done:
						return
					}
				}
			}
		}
	}()

	notebook.Wait()
	close(done)
	wg.Wait() // ensure watcher exits before closing channel

	// close stdin pipes first — sends EOF to plug processes so they exit cleanly
	for _, tab := range tabSessions {
		_ = tab.stdin.Close()
	}

	cleanupSharedMemory(key)
	killTabs(tabSessions)
	close(updateCh)

	viewerMu.Lock()
	session = nil
	viewerMu.Unlock()
}

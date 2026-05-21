package viewer

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"statehound/internal/statehound/events"
)

func startTabs(key int) (map[string]*tabSession, error) {
	tabSessions := make(map[string]*tabSession)

	for _, t := range tabs {
		tab, err := startTab(key, t)
		if err != nil {
			killTabs(tabSessions)
			return nil, err
		}

		tabSessions[t.urgency] = tab
	}

	return tabSessions, nil
}

func writeInitialLines(tabSessions map[string]*tabSession) {
	for _, t := range tabs {
		tab, ok := tabSessions[t.urgency]
		if !ok {
			continue
		}

		lines := initialLines(t.urgency)
		for _, line := range lines {
			writeLine(tab.stdin, line)
		}
	}
}

func startTab(key int, t tabConfig) (*tabSession, error) {
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
		return nil, fmt.Errorf("failed to create stdin pipe for %s: %w", t.urgency, err)
	}

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("failed to start tab for %s: %w", t.urgency, err)
	}

	return &tabSession{
		cmd:   cmd,
		stdin: stdin,
	}, nil
}

func initialLines(urgency string) []string {
	content, _ := events.FilterEvents(urgency)
	if content == "" {
		return []string{"no " + urgency + " events recorded"}
	}

	return strings.Split(strings.TrimSpace(content), "\n")
}

func writeLine(stdin io.WriteCloser, line string) {
	if strings.HasPrefix(line, "#") {
		line = " " + line
	}

	fmt.Fprintln(stdin, line)
}

func activeTabNum(urgency string) int {
	activeTab := 1

	for _, t := range tabs {
		if t.urgency == urgency {
			activeTab = t.tabnum
			break
		}
	}

	return activeTab
}

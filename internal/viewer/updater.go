package viewer

import (
	"strings"
	"time"

	"statehound/internal/statehound/events"
)

func startUpdateWriter(updateCh <-chan tabUpdate, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return

			case update, ok := <-updateCh:
				if !ok {
					return
				}

				viewerMu.Lock()
				s := session
				if s == nil {
					viewerMu.Unlock()
					continue
				}

				tab, ok := s.tabs[update.urgency]
				viewerMu.Unlock()

				if !ok {
					continue
				}

				for _, line := range update.lines {
					writeLine(tab.stdin, line)
				}
			}
		}
	}()
}

func startEventWatcher(done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return

			case <-ticker.C:
				updates := collectTabUpdates()
				sendTabUpdates(updates, done)
			}
		}
	}()
}

func collectTabUpdates() []tabUpdate {
	viewerMu.Lock()
	defer viewerMu.Unlock()

	s := session
	if s == nil {
		return nil
	}

	var updates []tabUpdate

	for _, t := range tabs {
		tab := s.tabs[t.urgency]
		if tab == nil {
			continue
		}

		content, _ := events.FilterEvents(t.urgency)
		lines := strings.Split(strings.TrimSpace(content), "\n")

		if len(lines) > tab.lastLines {
			newLines := make([]string, len(lines)-tab.lastLines)
			copy(newLines, lines[tab.lastLines:])
			tab.lastLines = len(lines)

			updates = append(updates, tabUpdate{
				urgency: t.urgency,
				lines:   newLines,
			})
		}
	}

	return updates
}

func sendTabUpdates(updates []tabUpdate, done <-chan struct{}) {
	viewerMu.Lock()
	s := session
	if s == nil {
		viewerMu.Unlock()
		return
	}

	ch := s.updateCh
	viewerMu.Unlock()

	for _, update := range updates {
		select {
		case ch <- update:
		case <-done:
			return
		}
	}
}

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
	s := session
	if s == nil {
		viewerMu.Unlock()
		return nil
	}

	offset := s.fileOffset
	viewerMu.Unlock()

	lines, newOffset, err := events.EventsSince(offset)
	if err != nil || len(lines) == 0 {
		return nil
	}

	updatesByUrgency := make(map[string][]string)

	for _, line := range lines {
		for _, t := range tabs {
			tag := "[" + t.urgency + "]"
			if strings.Contains(line, tag) {
				updatesByUrgency[t.urgency] = append(updatesByUrgency[t.urgency], line)
			}
		}
	}

	viewerMu.Lock()
	if session != nil {
		session.fileOffset = newOffset
	}
	viewerMu.Unlock()

	var updates []tabUpdate
	for _, t := range tabs {
		if tabLines := updatesByUrgency[t.urgency]; len(tabLines) > 0 {
			updates = append(updates, tabUpdate{
				urgency: t.urgency,
				lines:   tabLines,
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

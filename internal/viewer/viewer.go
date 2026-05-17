package viewer

import (
	"math/rand"
	"time"

	"statehound/internal/logger"
)

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
	viewerMu.Unlock()

	key := rand.Intn(0x7FFFFFFF) + 1

	tabSessions, err := startTabs(key)
	if err != nil {
		logger.Failed("viewer: failed to start tabs", err)
		return
	}

	time.Sleep(200 * time.Millisecond)

	notebook, err := startNotebook(key, activeTabNum(urgency))
	if err != nil {
		killTabs(tabSessions)
		logger.Failed("viewer: failed to start notebook", err)
		return
	}

	updateCh := make(chan tabUpdate, 32)

	viewerMu.Lock()
	session = &viewerSession{
		notebook: notebook,
		tabs:     tabSessions,
		updateCh: updateCh,
	}
	viewerMu.Unlock()

	done := make(chan struct{})
	stop := makeStopper(key, notebook, tabSessions, done)

	startUpdateWriter(updateCh, done)
	watchProcessExits(notebook, tabSessions, stop)
	startEventWatcher(done)

	<-done
}

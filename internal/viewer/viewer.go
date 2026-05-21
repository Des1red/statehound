package viewer

import (
	"math/rand"
	"os"
	"time"

	"statehound/internal/logger"
	"statehound/internal/model"
)

func Show(urgency string) {
	lockFile, ok := acquireLock()
	if !ok {
		return
	}
	defer releaseLock(lockFile)

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
	defer func() {
		if session == nil {
			killTabs(tabSessions)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	notebook, err := startNotebook(key, activeTabNum(urgency))
	if err != nil {
		killTabs(tabSessions)
		logger.Failed("viewer: failed to start notebook", err)
		return
	}
	go writeInitialLines(tabSessions)
	updateCh := make(chan tabUpdate, 32)

	info, err := os.Stat(model.EventPath)
	if err != nil {
		killTabs(tabSessions)
		logger.Failed("viewer: failed to stat event file", err)
		return
	}
	viewerMu.Lock()
	session = &viewerSession{
		notebook:   notebook,
		tabs:       tabSessions,
		updateCh:   updateCh,
		fileOffset: info.Size(),
	}
	viewerMu.Unlock()

	done := make(chan struct{})
	stop := makeStopper(key, notebook, tabSessions, done)

	startUpdateWriter(updateCh, done)
	watchProcessExits(notebook, tabSessions, stop)
	startEventWatcher(done)

	<-done
}

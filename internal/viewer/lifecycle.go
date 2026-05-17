package viewer

import (
	"os/exec"
	"sync"
)

func makeStopper(key int, notebook *exec.Cmd, tabSessions map[string]*tabSession, done chan struct{}) func() {
	var doneOnce sync.Once

	return func() {
		doneOnce.Do(func() {
			close(done)

			for _, tab := range tabSessions {
				_ = tab.stdin.Close()
			}

			if notebook != nil && notebook.Process != nil {
				_ = notebook.Process.Kill()
			}

			killTabs(tabSessions)
			cleanupSharedMemory(key)

			viewerMu.Lock()
			if session != nil {
				close(session.updateCh)
				session = nil
			}
			viewerMu.Unlock()
		})
	}
}

func watchProcessExits(notebook *exec.Cmd, tabSessions map[string]*tabSession, stop func()) {
	go func() {
		_ = notebook.Wait()
		stop()
	}()

	for _, tab := range tabSessions {
		tab := tab

		go func() {
			_ = tab.cmd.Wait()
			stop()
		}()
	}
}

package viewer

import (
	"errors"
	"os"
	"statehound/internal/logger"
	"syscall"
)

const lockPath = "/tmp/statehound-viewer.lock"

func acquireLock() (*os.File, bool) {
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Failed("viewer lock: failed to open lock file", err)
		return nil, false
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = f.Close()

		if errors.Is(err, syscall.EWOULDBLOCK) {
			logger.Status("viewer lock: already locked")
			return nil, false
		}

		logger.Failed("viewer lock: failed to acquire lock", err)
		return nil, false
	}

	logger.Status("viewer lock: acquired")
	return f, true
}

func releaseLock(f *os.File) {
	if f == nil {
		return
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		logger.Failed("viewer lock: failed to unlock", err)
	}

	if err := f.Close(); err != nil {
		logger.Failed("viewer lock: failed to close lock file", err)
	}

	if err := os.Remove(lockPath); err != nil && !os.IsNotExist(err) {
		logger.Failed("viewer lock: failed to remove lock file", err)
		return
	}

	logger.Status("viewer lock: released")
}

func isAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

func killTabs(tabSessions map[string]*tabSession) {
	for _, tab := range tabSessions {
		if tab.cmd.Process != nil {
			_ = tab.cmd.Process.Kill()
		}
	}
}

func IsOpen() bool {
	viewerMu.Lock()
	defer viewerMu.Unlock()
	return session != nil
}

func buildEnv() []string {
	env := append(os.Environ(), "GTK_THEME=Adwaita:dark")
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		env = append(env, "GDK_BACKEND=x11")
	}
	return env
}

func cleanupSharedMemory(key int) {
	id, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), 0, 0)
	if errno != 0 {
		return
	}
	_, _, _ = syscall.Syscall(syscall.SYS_SHMCTL, id, 0, 0)
}

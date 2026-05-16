package viewer

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

const lockPath = "/tmp/statehound-viewer.lock"

func acquireLock() bool {
	data, err := os.ReadFile(lockPath)
	if err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil && isAlive(pid) {
			return false
		}
		os.Remove(lockPath)
	}
	return os.WriteFile(lockPath, []byte(strconv.Itoa(os.Getpid())), 0644) == nil
}

func releaseLock() {
	os.Remove(lockPath)
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

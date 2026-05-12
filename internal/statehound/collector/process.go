package collector

import (
	"os"
	"strings"
)

func enrichPortProcess(p Port) Port {
	if p.PID == "" {
		return p
	}

	p.Exe = readProcessExe(p.PID)
	p.Cmdline = readProcessCmdline(p.PID)

	return p
}

func readProcessExe(pid string) string {
	exe, err := os.Readlink("/proc/" + pid + "/exe")
	if err != nil {
		return ""
	}

	return exe
}

func enrichConnectionProcess(c Connection) Connection {
	if c.PID == "" {
		return c
	}

	c.Exe = readProcessExe(c.PID)
	c.Cmdline = readProcessCmdline(c.PID)

	return c
}

func readProcessCmdline(pid string) string {
	data, err := os.ReadFile("/proc/" + pid + "/cmdline")
	if err != nil {
		return ""
	}

	cmdline := strings.ReplaceAll(string(data), "\x00", " ")
	return strings.TrimSpace(cmdline)
}

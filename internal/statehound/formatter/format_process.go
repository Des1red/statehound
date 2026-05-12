package formatter

import "statehound/internal/statehound/collector"

func FormatProcess(p collector.Process) string {
	msg := "pid=" + p.PID

	if p.PPID != "" {
		msg += " ppid=" + p.PPID
	}

	if p.UID != "" {
		msg += " uid=" + p.UID
	}

	if p.Name != "" {
		msg += " name=" + p.Name
	}

	if p.Exe != "" {
		msg += " exe=" + p.Exe
	}

	if p.Cmdline != "" {
		msg += " cmdline=\"" + shortCmdline(p.Cmdline) + "\""
	}

	if p.Reason != "" {
		msg += " reason=" + p.Reason
	}

	return msg
}

package formatter

import (
	"fmt"
	"statehound/internal/statehound/collector"
)

func FormatPort(p collector.Port) string {
	msg := fmt.Sprintf("%s %s:%s", p.Proto, p.Address, p.Port)

	if p.Scope != "" {
		msg += " scope=" + p.Scope
	}

	if p.Process != "" {
		msg += " process=" + p.Process
	}

	if p.PID != "" {
		msg += " pid=" + p.PID
	}

	if p.Exe != "" {
		msg += " exe=" + p.Exe
	}

	if p.Cmdline != "" {
		msg += " cmdline=\"" + shortCmdline(p.Cmdline) + "\""
	}

	return msg
}

func shortCmdline(s string) string {
	const max = 160

	if len(s) <= max {
		return s
	}

	return s[:max] + "..."
}

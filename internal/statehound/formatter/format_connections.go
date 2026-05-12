package formatter

import (
	"fmt"

	"statehound/internal/statehound/collector"
)

func FormatConnection(c collector.Connection) string {
	msg := fmt.Sprintf(
		"%s %s:%s -> %s:%s state=%s",
		c.Proto,
		c.LocalAddress,
		c.LocalPort,
		c.RemoteAddress,
		c.RemotePort,
		c.State,
	)

	if c.Process != "" {
		msg += " process=" + c.Process
	}

	if c.PID != "" {
		msg += " pid=" + c.PID
	}

	if c.Exe != "" {
		msg += " exe=" + c.Exe
	}

	if c.Cmdline != "" {
		msg += " cmdline=\"" + shortCmdline(c.Cmdline) + "\""
	}

	return msg
}

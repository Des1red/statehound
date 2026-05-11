package statehound

import "fmt"

func formatServiceTransition(previous, current Service) string {
	return formatService(current) +
		" transition=" +
		previous.ActiveState + "/" + previous.SubState +
		"->" +
		current.ActiveState + "/" + current.SubState
}

func formatService(s Service) string {
	msg := s.Name

	if s.Description != "" {
		msg += ` desc="` + s.Description + `"`
	}

	if s.LoadState != "" {
		msg += " load=" + s.LoadState
	}

	if s.ActiveState != "" || s.SubState != "" {
		msg += " state=" + s.ActiveState + "/" + s.SubState
	}

	return msg
}

func formatPort(p Port) string {
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

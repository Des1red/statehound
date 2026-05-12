package formatter

import "statehound/internal/statehound/collector"

func FormatService(s collector.Service) string {
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

func FormatServiceTransition(previous, current collector.Service) string {
	return FormatService(current) +
		" transition=" +
		previous.ActiveState + "/" + previous.SubState +
		"->" +
		current.ActiveState + "/" + current.SubState
}

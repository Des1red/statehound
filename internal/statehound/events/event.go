package events

import (
	"fmt"
	"os"
	"strings"
	"time"

	"statehound/internal/model"
	"statehound/internal/statehound/collector"
)

type Event struct {
	Time    time.Time
	Type    string
	Message string

	Urgency string
	Tags    []string

	Connection *collector.Connection // optional, only set for CONNECTION_* events
	Process    *collector.Process    // optional, for suspicious processes
	File       *collector.FileWatch  // optional, for FILE_* events
	Service    *collector.Service    // optional, for SERVICE_* events
	Port       *collector.Port       // optional, for PORT_* events
}

func (e *Event) TagUrgency(level string) {
	level = strings.TrimSpace(level)
	if level == "" {
		return
	}
	e.Urgency = level
}

func (e *Event) Tag(tag string) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return
	}

	for _, existing := range e.Tags {
		if existing == tag {
			return
		}
	}

	e.Tags = append(e.Tags, tag)
}

func (e Event) HasTag(tag string) bool {
	for _, existing := range e.Tags {
		if existing == tag {
			return true
		}
	}

	return false
}

func (e Event) string() string {
	if e.Urgency != "" {
		return fmt.Sprintf(
			"[%s] [%s] [%s] %s\n",
			e.Time.Format(time.RFC3339),
			e.Urgency,
			e.Type,
			e.Message,
		)
	}

	return fmt.Sprintf(
		"[%s] %s %s\n",
		e.Time.Format(time.RFC3339),
		e.Type,
		e.Message,
	)
}

func WriteEvents(events []Event) error {
	if len(events) == 0 {
		return nil
	}

	cleanupEventLogIfNeeded()

	file, err := os.OpenFile(model.EventPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open event log: %w", err)
	}
	defer file.Close()

	for _, event := range events {
		if _, err := file.WriteString(event.string()); err != nil {
			return fmt.Errorf("failed to write event: %w", err)
		}
	}

	return nil
}

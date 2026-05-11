package statehound

import (
	"fmt"
	"os"
	"time"

	"statehound/internal/model"
)

type Event struct {
	Time    time.Time
	Type    string
	Message string
}

func (e Event) string() string {
	return fmt.Sprintf(
		"[%s] %s %s\n",
		e.Time.Format(time.RFC3339),
		e.Type,
		e.Message,
	)
}

func writeEvents(events []Event) error {
	if len(events) == 0 {
		return nil
	}

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

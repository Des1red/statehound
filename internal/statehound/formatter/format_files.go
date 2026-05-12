package formatter

import (
	"fmt"

	"statehound/internal/statehound/collector"
)

func FormatFile(file collector.FileWatch) string {
	msg := file.Path

	if file.Exists {
		msg += fmt.Sprintf(
			" size=%d mode=%s modtime=%s",
			file.Size,
			file.Mode,
			file.ModTime.Format("2006-01-02T15:04:05Z07:00"),
		)
	}

	if file.Hash != "" {
		msg += " hash=" + file.Hash
	}

	return msg
}

func FormatFileTransition(previous, current collector.FileWatch) string {
	return FormatFile(current) +
		" previous_hash=" + previous.Hash +
		" current_hash=" + current.Hash
}

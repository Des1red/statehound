package events

import (
	"io"
	"os"
	"strings"

	"statehound/internal/model"
)

func FilterEvents(filter string) (string, error) {
	data, err := os.ReadFile(model.EventPath)
	if err != nil {
		return "", err
	}

	if filter == "" {
		return string(data), nil
	}

	tag := "[" + filter + "]"
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, tag) {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n"), nil
}

func EventsSince(offset int64) ([]string, int64, error) {
	f, err := os.Open(model.EventPath)
	if err != nil {
		return nil, offset, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, offset, err
	}

	// Log was rotated/truncated.
	if info.Size() < offset {
		offset = 0
	}

	if info.Size() == offset {
		return nil, offset, nil
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, offset, err
	}

	newBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, offset, err
	}

	newOffset := info.Size()

	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(string(newBytes)), "\n") {
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}

	return lines, newOffset, nil
}

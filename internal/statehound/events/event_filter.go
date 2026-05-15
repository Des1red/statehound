package events

import (
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

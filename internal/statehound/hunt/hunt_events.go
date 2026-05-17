package hunt

import (
	"os"
	"strings"

	"statehound/internal/model"
)

func HuntEvents(target string) []string {
	data, err := os.ReadFile(model.EventPath)
	if err != nil {
		return nil
	}

	var matches []string
	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		if strings.Contains(strings.ToLower(line), target) {
			matches = append(matches, "  "+line)
		}
	}

	return matches
}

package collector

import (
	"strings"

	"statehound/internal/command"
)

type ServiceDetails struct {
	Name         string
	Description  string
	LoadState    string
	ActiveState  string
	SubState     string
	MainPID      string
	FragmentPath string
}

func collectSystemdServiceDetails(name string) ServiceDetails {
	out, err := command.Output(
		"systemctl",
		"show",
		name,
		"--property=Id,Description,LoadState,ActiveState,SubState,MainPID,FragmentPath",
		"--no-pager",
	)
	if err != nil {
		return ServiceDetails{Name: name}
	}

	values := parseSystemdShow(out)

	return ServiceDetails{
		Name:         first(values["Id"], name),
		Description:  values["Description"],
		LoadState:    values["LoadState"],
		ActiveState:  values["ActiveState"],
		SubState:     values["SubState"],
		MainPID:      values["MainPID"],
		FragmentPath: values["FragmentPath"],
	}
}

func parseSystemdShow(out []byte) map[string]string {
	values := make(map[string]string)

	for _, line := range strings.Split(string(out), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		values[key] = value
	}

	return values
}

func first(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}

func CountActiveServices(services map[string]Service) int {
	count := 0

	for _, service := range services {
		if service.ActiveState == "active" {
			count++
		}
	}

	return count
}

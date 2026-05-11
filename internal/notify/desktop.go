package notify

import "os/exec"

func Desktop(title, message, urgency string) error {
	if urgency == "" {
		urgency = "normal"
	}

	cmd := exec.Command(
		"notify-send",
		"-u",
		urgency,
		title,
		message,
	)

	return cmd.Run()
}

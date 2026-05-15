package notify

import (
	"bufio"
	"os/exec"
	"statehound/internal/logger"
	"strings"
)

type Session struct {
	ID     string
	action chan bool
}

func (s *Session) ActionCh() <-chan bool {
	return s.action
}

func Start(title, message, urgency string) (*Session, error) {
	cmd := exec.Command(
		"notify-send",
		"--print-id",
		"--wait",
		"--action=view=View events",
		"-u", urgency,
		title,
		message,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	idCh := make(chan string, 1)
	actionCh := make(chan bool, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		if scanner.Scan() {
			id := strings.TrimSpace(scanner.Text())
			logger.Status("notification id: " + id)
			idCh <- id
		} else {
			idCh <- ""
		}

		action := ""
		if scanner.Scan() {
			action = strings.TrimSpace(scanner.Text())
			logger.Status("notification action received: " + action)
		}

		cmd.Wait()
		actionCh <- action == "view"
	}()

	return &Session{
		ID:     <-idCh,
		action: actionCh,
	}, nil
}

func Send(title, message, urgency string) (string, error) {
	out, err := exec.Command(
		"notify-send",
		"--print-id",
		"-u", urgency,
		title,
		message,
	).Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func Replace(id, title, message, urgency string) (string, error) {
	out, err := exec.Command(
		"notify-send",
		"--print-id",
		"--replace-id="+id,
		"-u", urgency,
		title,
		message,
	).Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

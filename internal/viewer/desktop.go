package viewer

import (
	"bufio"
	"os/exec"
	"strings"
)

type Session struct {
	ID     string
	action chan bool
}

func (s *Session) ActionCh() <-chan bool {
	return s.action
}

func start(title, message, urgency string) (*Session, error) {
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
			idCh <- id
		} else {
			idCh <- ""
		}

		action := ""
		if scanner.Scan() {
			action = strings.TrimSpace(scanner.Text())
		}

		cmd.Wait()
		actionCh <- action == "view"
	}()

	return &Session{
		ID:     <-idCh,
		action: actionCh,
	}, nil
}

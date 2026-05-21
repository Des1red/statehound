package client

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"statehound/internal/model"
)

func Send(command string) (string, error) {
	return SendTo(model.ControlSocketPath, command)
}

func SendNotify(command string) (string, error) {
	return SendTo(model.NotifySocketPath, command)
}

func IsNotifyRunning() bool {
	resp, err := SendNotify("PING")
	return err == nil && resp == "PONG"
}

func SendTo(path, command string) (string, error) {
	conn, err := net.DialTimeout("unix", path, 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to statehound socket %s: %w", path, err)
	}
	defer conn.Close()

	if _, err := fmt.Fprintln(conn, command); err != nil {
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	data, err := io.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read daemon response: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func IsRunning() bool {
	resp, err := Send("PING")
	return err == nil && resp == "PONG"
}

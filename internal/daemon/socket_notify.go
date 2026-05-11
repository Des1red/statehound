package daemon

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"

	"statehound/internal/logger"
	"statehound/internal/model"
	"statehound/internal/statehound"
)

func startNotifySocket(path string, manager *statehound.Manager) error {
	_ = os.Remove(path)

	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	defer listener.Close()
	defer os.Remove(path)

	if err := setNotifySocketPermissions(path); err != nil {
		return err
	}

	logger.Status("statehound notify socket listening: " + path)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleNotifyConnection(conn, manager)
	}
}

func setNotifySocketPermissions(path string) error {
	group, err := user.LookupGroup(model.GroupName)
	if err != nil {
		return fmt.Errorf("failed to lookup group %s: %w", model.GroupName, err)
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return fmt.Errorf("invalid group gid %q: %w", group.Gid, err)
	}

	if err := os.Chown(path, 0, gid); err != nil {
		return err
	}

	if err := os.Chmod(path, 0660); err != nil {
		return err
	}

	return nil
}

func handleNotifyConnection(conn net.Conn, manager *statehound.Manager) {
	defer conn.Close()

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	msg = strings.TrimSpace(msg)

	switch msg {
	case "PING":
		writeResponse(conn, "PONG")
	case "NOTIFICATIONS":
		writeResponse(conn, manager.NotificationsJSON())
	default:
		writeResponse(conn, "unknown command")
	}
}

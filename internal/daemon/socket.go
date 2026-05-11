package daemon

import (
	"bufio"
	"net"
	"os"
	"statehound/internal/logger"
	"statehound/internal/statehound"
	"strings"
)

func startSocket(path string, manager *statehound.Manager) error {
	_ = os.Remove(path)

	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	defer listener.Close()
	defer os.Remove(path)

	if err := os.Chmod(path, 0600); err != nil {
		return err
	}

	logger.Status("statehound socket listening: " + path)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn, manager)
	}
}

func handleConnection(conn net.Conn, manager *statehound.Manager) {
	defer conn.Close()

	if err := verifyPeer(conn); err != nil {
		writeResponse(conn, "unauthorized")
		return
	}

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	msg = strings.TrimSpace(msg)

	switch {
	case msg == "PING":
		writeResponse(conn, "PONG")
	case msg == "STATUS":
		writeResponse(conn, manager.Status())
	case strings.HasPrefix(msg, "HUNT "):
		target := strings.TrimSpace(strings.TrimPrefix(msg, "HUNT "))
		writeResponse(conn, manager.Hunt(target))
	default:
		writeResponse(conn, "unknown command")
	}
}

func writeResponse(conn net.Conn, msg string) {
	_, _ = conn.Write([]byte(msg + "\n"))
}

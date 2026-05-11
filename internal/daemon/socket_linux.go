package daemon

import (
	"fmt"
	"net"
	"syscall"
)

func verifyPeer(conn net.Conn) error {
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return fmt.Errorf("connection is not unix connection")
	}

	rawConn, err := unixConn.SyscallConn()
	if err != nil {
		return fmt.Errorf("failed to get raw connection: %w", err)
	}

	var cred *syscall.Ucred
	var sockErr error

	err = rawConn.Control(func(fd uintptr) {
		cred, sockErr = syscall.GetsockoptUcred(
			int(fd),
			syscall.SOL_SOCKET,
			syscall.SO_PEERCRED,
		)
	})
	if err != nil {
		return fmt.Errorf("failed socket control: %w", err)
	}

	if sockErr != nil {
		return fmt.Errorf("failed to get peer credentials: %w", sockErr)
	}

	if cred == nil {
		return fmt.Errorf("missing peer credentials")
	}

	if cred.Uid != 0 {
		return fmt.Errorf("unauthorized peer uid: %d", cred.Uid)
	}

	return nil
}

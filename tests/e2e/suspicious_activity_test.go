// to run this test :
// sudo STATEHOUND_E2E=1 go test ./tests/e2e -v
package e2e

import (
	"net"
	"os"
	"os/user"
	"testing"
	"time"
)

const cronTestPath = "/etc/cron.d/statehound-e2e-test"

func TestStatehoundSuspiciousActivity(t *testing.T) {
	if os.Getenv("STATEHOUND_E2E") != "1" {
		t.Skip("set STATEHOUND_E2E=1 to run this integration test")
	}

	u, err := user.Current()
	if err != nil {
		t.Fatalf("failed to get current user: %v", err)
	}

	if u.Uid != "0" {
		t.Skip("this test must run as root because it writes /etc/cron.d")
	}

	t.Cleanup(func() {
		_ = os.Remove(cronTestPath)
	})

	t.Run("cron persistence file add/remove", func(t *testing.T) {
		data := []byte("# statehound e2e cron test\n")

		if err := os.WriteFile(cronTestPath, data, 0644); err != nil {
			t.Fatalf("failed to create cron test file: %v", err)
		}

		time.Sleep(7 * time.Second)

		if err := os.Remove(cronTestPath); err != nil {
			t.Fatalf("failed to remove cron test file: %v", err)
		}

		time.Sleep(7 * time.Second)
	})

	t.Run("public tcp listener", func(t *testing.T) {
		ln, err := net.Listen("tcp", "0.0.0.0:4444")
		if err != nil {
			t.Fatalf("failed to open public listener: %v", err)
		}

		time.Sleep(7 * time.Second)

		if err := ln.Close(); err != nil {
			t.Fatalf("failed to close listener: %v", err)
		}

		time.Sleep(7 * time.Second)
	})
}

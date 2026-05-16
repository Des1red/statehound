package e2e

import (
	"os"
	"statehound/internal/client"
	"strings"
	"testing"
)

func TestDeameonComms(t *testing.T) {
	if os.Getenv("STATEHOUND_E2E") != "1" {
		t.Skip("set STATEHOUND_E2E=1 to run this integration test")
	}
	if !client.IsRunning() {
		t.Skip("statehound daemon not running — run sudo shound --install first")
	}
	t.Run("socket responds to ping", func(t *testing.T) {
		resp, err := client.Send("PING")
		if err != nil {
			t.Fatalf("failed to ping daemon: %v", err)
		}
		if resp != "PONG" {
			t.Errorf("expected PONG got %q", resp)
		}
	})

	t.Run("socket responds to status", func(t *testing.T) {
		resp, err := client.Send("STATUS")
		if err != nil {
			t.Fatalf("failed to get status: %v", err)
		}
		if !strings.Contains(resp, "manager=") {
			t.Errorf("unexpected status response: %q", resp)
		}
	})

	t.Run("unauthorized command rejected", func(t *testing.T) {
		resp, err := client.Send("unknown command here")
		if err != nil {
			t.Fatalf("failed to send: %v", err)
		}
		if resp != "unknown command" {
			t.Errorf("expected unknown command rejection got %q", resp)
		}
	})
}

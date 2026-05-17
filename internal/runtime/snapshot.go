package runtime

import (
	"fmt"

	"statehound/internal/client"
	"statehound/internal/logger"
	"statehound/internal/system"
	"statehound/internal/viewer"
)

func Snapshot() {
	if system.IsHeadless() {
		snapshotHeadless()
		return
	}
	snapshotGUI()
}

func snapshotHeadless() {
	resp, err := client.Send("SNAPSHOT")
	if err != nil {
		logger.Failed("failed to read statehound snapshot", err)
		return
	}

	fmt.Println(resp)
}

func snapshotGUI() {
	viewer.ShowSnapshot(func() (string, error) {
		return client.Send("SNAPSHOT")
	})
}

package runtime

import (
	"statehound/internal/client"
	"statehound/internal/viewer"
)

func SnapshotGUI() {
	viewer.ShowSnapshot(func() (string, error) {
		return client.Send("SNAPSHOT")
	})
}

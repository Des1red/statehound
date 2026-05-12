package runtime

import (
	"fmt"

	"statehound/internal/client"
	"statehound/internal/logger"
)

func Snapshot() {
	resp, err := client.Send("SNAPSHOT")
	if err != nil {
		logger.Failed("failed to read statehound snapshot", err)
		return
	}

	fmt.Println(resp)
}

package runtime

import (
	"fmt"

	"statehound/internal/client"
	"statehound/internal/command"
	"statehound/internal/model"
)

func Stop() {
	if !client.IsRunning() {
		fmt.Println("[*] statehound is already stopped")
		return
	}

	if err := command.Run("systemctl", "stop", model.ServiceName); err != nil {
		fmt.Println("[!] failed to stop statehound:", err)
		return
	}

	if client.IsRunning() {
		fmt.Println("[!] statehound stop command was sent, but daemon is still responding")
		return
	}

	fmt.Println("[+] statehound stopped")
}

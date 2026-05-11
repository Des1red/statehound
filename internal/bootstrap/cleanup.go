package bootstrap

import (
	"os"
	"os/user"
	"statehound/internal/command"
	"statehound/internal/model"
)

func cleanupExistingInstall() {
	cleanupNotifierForSudoUser()

	_ = command.Run("systemctl", "stop", model.ServiceName)
	_ = command.Run("systemctl", "disable", model.ServiceName)
	_ = os.Remove(model.ServicePath)
	_ = os.Remove(model.NotifierServicePath)
	_ = os.Remove(model.AliasPath)
	_ = os.Remove(model.BinaryPath)
	_ = command.Run("systemctl", "daemon-reload")
	_ = command.Run("systemctl", "reset-failed", model.ServiceName)
}

func cleanupNotifierForSudoUser() {
	username := os.Getenv("SUDO_USER")
	if username == "" || username == "root" {
		return
	}

	u, err := user.Lookup(username)
	if err != nil {
		return
	}

	runtimeDir := "/run/user/" + u.Uid
	if _, err := os.Stat(runtimeDir); err != nil {
		return
	}

	envPrefix := "XDG_RUNTIME_DIR=" + runtimeDir

	_ = command.Run(
		"runuser",
		"-u",
		username,
		"--",
		"env",
		envPrefix,
		"systemctl",
		"--user",
		"disable",
		"--now",
		model.NotifierServiceName,
	)

	_ = command.Run(
		"runuser",
		"-u",
		username,
		"--",
		"env",
		envPrefix,
		"systemctl",
		"--user",
		"daemon-reload",
	)
}

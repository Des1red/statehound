package bootstrap

import (
	"os"
	"os/user"
	"statehound/internal/command"
	"statehound/internal/model"
	"statehound/internal/system"
)

func cleanupExistingInstall() {
	cleanupNotifierForSudoUser()
	_ = command.RunSilent("systemctl", "stop", model.ServiceName)
	_ = command.RunSilent("systemctl", "disable", model.ServiceName)
	_ = os.Remove(model.ServicePath)
	if !system.IsHeadless() {
		_ = os.Remove(model.NotifierServicePath)
	}
	_ = os.Remove(model.AliasPath)
	_ = os.Remove(model.BinaryPath)
	_ = command.RunSilent("systemctl", "daemon-reload")
	_ = command.RunSilent("systemctl", "reset-failed", model.ServiceName)
}

func cleanupNotifierForSudoUser() {
	if system.IsHeadless() {
		return
	}

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

	_ = command.RunSilent(
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

	_ = command.RunSilent(
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

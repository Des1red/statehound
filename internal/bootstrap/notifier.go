package bootstrap

import (
	"fmt"
	"os"
	"os/user"

	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

func setupNotifierForSudoUser() error {
	username := os.Getenv("SUDO_USER")
	if username == "" || username == "root" {
		return fmt.Errorf("could not detect desktop user from SUDO_USER")
	}

	u, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("failed to lookup sudo user %q: %w", username, err)
	}

	if err := command.Run("usermod", "-aG", model.GroupName, username); err != nil {
		return fmt.Errorf("failed to add user %q to %s group: %w", username, model.GroupName, err)
	}

	runtimeDir := "/run/user/" + u.Uid

	if _, err := os.Stat(runtimeDir); err != nil {
		return fmt.Errorf("user runtime dir not available: %s", runtimeDir)
	}

	if err := grantNotifierSocketAccess(username); err != nil {
		return err
	}

	envPrefix := "XDG_RUNTIME_DIR=" + runtimeDir

	if err := command.Run(
		"runuser",
		"-u",
		username,
		"--",
		"env",
		envPrefix,
		"systemctl",
		"--user",
		"daemon-reload",
	); err != nil {
		return fmt.Errorf("failed to reload user systemd for %q: %w", username, err)
	}

	if err := command.Run(
		"runuser",
		"-u",
		username,
		"--",
		"env",
		envPrefix,
		"systemctl",
		"--user",
		"enable",
		"--now",
		model.NotifierServiceName,
	); err != nil {
		return fmt.Errorf("failed to enable/start notifier for %q: %w", username, err)
	}

	logger.Success("statehound desktop notifier enabled for " + username)
	return nil
}

func grantNotifierSocketAccess(username string) error {
	if err := command.Run("setfacl", "-m", "u:"+username+":rx", model.RuntimeDir); err != nil {
		return fmt.Errorf("failed to grant runtime dir access to %q: %w", username, err)
	}

	if err := command.Run("setfacl", "-m", "u:"+username+":rw", model.NotifySocketPath); err != nil {
		return fmt.Errorf("failed to grant notify socket access to %q: %w", username, err)
	}

	return nil
}

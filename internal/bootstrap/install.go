package bootstrap

import (
	"fmt"
	"os"

	"statehound/internal/command"
	"statehound/internal/logger"
	"statehound/internal/model"
)

const unit = `[Unit]
Description=statehound Host Monitor
After=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=` + model.BinaryPath + ` --daemon
Restart=on-failure
RestartSec=3
RuntimeDirectory=statehound
RuntimeDirectoryMode=0750

[Install]
WantedBy=multi-user.target
`

const notifierUnit = `[Unit]
Description=Statehound Desktop Notifier
After=graphical-session.target

[Service]
Type=simple
ExecStart=` + model.BinaryPath + ` --notify
Restart=on-failure
RestartSec=3

[Install]
WantedBy=default.target
`

func Install() {
	logger.Status("installing statehound")

	cleanupExistingInstall()
	ensureNotifierDeps()

	if err := command.Run("go", "build", "-o", model.BinaryPath, "."); err != nil {
		logger.Failed("failed to build/install binary", err)
		return
	}

	if err := ensureGroup(); err != nil {
		logger.Failed("", err)
		return
	}

	if err := createFiles(); err != nil {
		logger.Failed("", err)
		return
	}

	if err := enforcePermissions(); err != nil {
		logger.Failed("", err)
		return
	}

	if err := createAlias(); err != nil {
		logger.Failed("", err)
		return
	}

	if err := finalizeInstall(); err != nil {
		logger.Failed("", err)
		return
	}

	if err := setupNotifierForSudoUser(); err != nil {
		logger.Warn("desktop notifier was not enabled automatically")
		logger.Warn(err.Error())
	}

	logger.Success("statehound installed, enabled, and started")
}

func createAlias() error {
	_ = os.Remove(model.AliasPath)

	if err := os.Symlink(model.BinaryPath, model.AliasPath); err != nil {
		return fmt.Errorf("failed to create shound alias: %w", err)
	}

	if err := os.Lchown(model.AliasPath, 0, 0); err != nil {
		return fmt.Errorf("failed to set shound alias owner: %w", err)
	}

	return nil
}

func finalizeInstall() error {
	if err := command.Run("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	if err := command.Run("systemctl", "enable", model.ServiceName); err != nil {
		return fmt.Errorf("failed to enable statehound service: %w", err)
	}

	if !VerifyInstallation() {
		return fmt.Errorf("install completed but verification failed")
	}

	if err := command.Run("systemctl", "restart", model.ServiceName); err != nil {
		return fmt.Errorf("statehound installed but failed to start service: %w", err)
	}

	return nil
}

func ensureGroup() error {
	if err := command.Run("getent", "group", model.GroupName); err == nil {
		return nil
	}

	if err := command.Run("groupadd", "--system", model.GroupName); err != nil {
		return fmt.Errorf("failed to create statehound group: %w", err)
	}

	return nil
}

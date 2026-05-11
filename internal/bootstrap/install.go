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
RuntimeDirectoryMode=0700

[Install]
WantedBy=multi-user.target
`

func Install() {
	logger.Status("installing statehound")

	cleanupExistingInstall()

	if err := command.Run("go", "build", "-o", model.BinaryPath, "."); err != nil {
		logger.Failed("failed to build/install binary", err)
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

	logger.Success("statehound installed, enabled, and started")
}

func cleanupExistingInstall() {
	_ = command.Run("systemctl", "stop", model.ServiceName)
	_ = command.Run("systemctl", "disable", model.ServiceName)
	_ = os.Remove(model.ServicePath)
	_ = os.Remove(model.AliasPath)
	_ = os.Remove(model.BinaryPath)
	_ = command.Run("systemctl", "daemon-reload")
	_ = command.Run("systemctl", "reset-failed", model.ServiceName)
}

func createFiles() error {
	if err := os.WriteFile(model.ServicePath, []byte(unit), 0644); err != nil {
		return fmt.Errorf("failed to write systemd service: %w", err)
	}

	if err := os.MkdirAll(model.ConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.MkdirAll(model.LogDir, 0700); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	if _, err := os.Stat(model.EventPath); os.IsNotExist(err) {
		if err := os.WriteFile(model.EventPath, []byte{}, 0600); err != nil {
			return fmt.Errorf("failed to create event log file: %w", err)
		}
	}

	if err := os.MkdirAll(model.EventBackupDir, 0700); err != nil {
		return fmt.Errorf("failed to create event backup directory: %w", err)
	}

	return nil
}

func enforcePermissions() error {
	if err := os.Chown(model.BinaryPath, 0, 0); err != nil {
		return fmt.Errorf("failed to set binary owner: %w", err)
	}

	if err := os.Chmod(model.BinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set binary permissions: %w", err)
	}

	if err := os.Chown(model.ServicePath, 0, 0); err != nil {
		return fmt.Errorf("failed to set service owner: %w", err)
	}

	if err := os.Chmod(model.ServicePath, 0644); err != nil {
		return fmt.Errorf("failed to set service permissions: %w", err)
	}

	if err := os.Chown(model.ConfigDir, 0, 0); err != nil {
		return fmt.Errorf("failed to set config directory owner: %w", err)
	}

	if err := os.Chmod(model.ConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to set config directory permissions: %w", err)
	}

	if err := os.Chown(model.LogDir, 0, 0); err != nil {
		return fmt.Errorf("failed to set log directory owner: %w", err)
	}

	if err := os.Chmod(model.LogDir, 0700); err != nil {
		return fmt.Errorf("failed to set log directory permissions: %w", err)
	}

	if err := os.Chown(model.EventPath, 0, 0); err != nil {
		return fmt.Errorf("failed to set event log owner: %w", err)
	}

	if err := os.Chmod(model.EventPath, 0600); err != nil {
		return fmt.Errorf("failed to set event log permissions: %w", err)
	}

	if err := os.Chown(model.EventBackupDir, 0, 0); err != nil {
		return fmt.Errorf("failed to set event backup directory owner: %w", err)
	}

	if err := os.Chmod(model.EventBackupDir, 0600); err != nil {
		return fmt.Errorf("failed to set event backup directory permissions: %w", err)
	}

	return nil
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

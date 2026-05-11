package events

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"statehound/internal/logger"
	"statehound/internal/model"
)

func cleanupEventLogIfNeeded() {
	if err := cleanupOldEventBackups(); err != nil {
		logger.Failed("failed to clean old event backups", err)
	}

	info, err := os.Stat(model.EventPath)
	if err != nil {
		return
	}

	if info.Size() < model.MaxEventLogSizeBytes {
		return
	}

	if err := backupAndTruncateEventLog(); err != nil {
		logger.Failed("failed to rotate event log", err)
	}
}

func backupAndTruncateEventLog() error {
	if err := os.MkdirAll(model.EventBackupDir, 0700); err != nil {
		return err
	}

	data, err := os.ReadFile(model.EventPath)
	if err != nil {
		return err
	}

	name := "events-" + time.Now().Format("2006-01-02T15-04-05") + ".log"
	backupPath := filepath.Join(model.EventBackupDir, name)

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return err
	}

	if err := os.Chown(backupPath, 0, 0); err != nil {
		return err
	}

	if err := os.WriteFile(model.EventPath, []byte{}, 0600); err != nil {
		return err
	}

	return nil
}

func cleanupOldEventBackups() error {
	entries, err := os.ReadDir(model.EventBackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -model.EventBackupMaxAgeDays)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(model.EventBackupDir, entry.Name())

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove old backup %s: %w", path, err)
			}
		}
	}

	return nil
}

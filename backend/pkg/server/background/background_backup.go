package background

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func StartDatabaseBackup(
	ctx context.Context, logger *slog.Logger, db database.Storage, cfg *config.Config,
) <-chan struct{} {
	logger.Info("Starting database backup task...")

	done := make(chan struct{})

	go func() {
		defer close(done)

		// Run immediately on start
		performBackup(logger, db, cfg.DBPath)

		// Check every hour if we need to run backup (or just simple 24h ticker as per plan)
		// Plan said "Schedule a ticker for 24 hours".
		// To be more robust about "daily", we could calculate time until next midnight,
		// but 24h ticker starting from launch is acceptable MVP as per plan.
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopped database backup task")
				return
			case <-ticker.C:
				performBackup(logger, db, cfg.DBPath)
			}
		}
	}()

	return done
}

func performBackup(logger *slog.Logger, db database.Storage, originalDBPath string) {
	logger.Info("Running database backup...")

	dir := filepath.Dir(originalDBPath)
	filename := filepath.Base(originalDBPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	dateSuffix := time.Now().Format("2006-01-02")
	backupFilename := fmt.Sprintf("%s_%s%s", nameWithoutExt, dateSuffix, ext)
	backupPath := filepath.Join(dir, backupFilename)

	// Check if exists
	if _, err := os.Stat(backupPath); err == nil {
		logger.Info("Backup file already exists, skipping", "path", backupPath)
		return
	}

	if err := db.Backup(backupPath); err != nil {
		logger.Error("Failed to create database backup", "error", err, "path", backupPath)
		return
	}

	logger.Info("Database backup created successfully", "path", backupPath)
}

package background

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func TestStartDatabaseBackup(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "geekbudget_bg_backup_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a dummy DB file
	file, err := os.Create(dbPath)
	require.NoError(t, err)
	file.Close()

	// Initialize real storage with this DB
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{
		DBPath:  dbPath,
		Verbose: true,
	}

	storage := database.NewStorage(logger, cfg)
	err = storage.Open()
	require.NoError(t, err)
	// Create tables so VACUUM INTO works on a valid DB structure
	// (NewStorage/Open already does autoMigrateModels)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start backup task
	done := StartDatabaseBackup(ctx, logger, storage, cfg)

	// Wait a bit for the immediate backup to finish
	// In a real scenario we might want to sync differently, but for this simple test sleep is okay
	// or we can loop checking file existence

	dateSuffix := time.Now().Format("2006-01-02")
	expectedBackupPath := filepath.Join(tmpDir, "test_"+dateSuffix+".db")

	assert.Eventually(t, func() bool {
		_, err := os.Stat(expectedBackupPath)
		return err == nil
	}, 2*time.Second, 100*time.Millisecond, "Backup file should be created immediately on start")

	// Clean up
	cancel()
	<-done
}

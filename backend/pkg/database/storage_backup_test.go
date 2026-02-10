package database

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

func TestStorage_Backup(t *testing.T) {
	// Create a temporary directory for the test database
	tmpDir, err := os.MkdirTemp("", "geekbudget_backup_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	backupPath := filepath.Join(tmpDir, "backup.db")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{
		DBPath:  dbPath,
		Verbose: true,
	}

	storage := NewStorage(logger, cfg)
	err = storage.Open()
	require.NoError(t, err)
	defer storage.Close()

	// Insert some data to verify backup content (optional, but good practice)
	// For simplicity, we just rely on open/close and VACUUM success.
	// But let's create a user to be sure DB is initialized
	_, err = storage.CreateUser("testuser", "password")
	require.NoError(t, err)

	// Perform backup
	err = storage.Backup(backupPath)
	assert.NoError(t, err)

	// Verify backup file exists
	_, err = os.Stat(backupPath)
	assert.NoError(t, err)

	// Verify backup file is a valid SQLite database by opening it
	backupCfg := &config.Config{
		DBPath:  backupPath,
		Verbose: true,
	}
	backupStorage := NewStorage(logger, backupCfg)
	err = backupStorage.Open()
	assert.NoError(t, err)
	defer backupStorage.Close()

	// Check if user exists in backup
	_, err = backupStorage.GetUserID("testuser")
	assert.NoError(t, err)
}

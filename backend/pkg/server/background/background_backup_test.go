package background

import (
	"archive/tar"
	"compress/gzip"
	"io"
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

func TestDatabaseBackupRun(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	storage := database.NewStorage(logger, &config.Config{DBPath: dbPath})
	require.NoError(t, storage.Open())

	// Create a bank-importer-files directory with a dummy file
	bankDir := filepath.Join(tmpDir, "bank-importer-files")
	require.NoError(t, os.MkdirAll(bankDir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(bankDir, "statement.csv"), []byte("date,amount\n"), 0o640))

	cfg := &config.Config{
		DBPath:                dbPath,
		BankImporterFilesPath: bankDir,
		BackupInterval:        "24h",
		BackupMaxCount:        10,
		// BackupPath empty → derived to tmpDir/geekbudget-backups/
	}
	task := newDatabaseBackupTask(logger, cfg)
	task.run()

	expectedBackupDir := filepath.Join(tmpDir, "geekbudget-backups")
	expectedFile := filepath.Join(expectedBackupDir,
		"geekbudget-backup-"+time.Now().Format("2006-01-02")+".tar.gz")

	_, err := os.Stat(expectedFile)
	require.NoError(t, err, "backup archive should exist")

	assertArchiveContains(t, expectedFile, "geekbudget.db", "bank-importer-files/statement.csv")

	// Second run same day — should skip (idempotent)
	task.run()
	entries, err := os.ReadDir(expectedBackupDir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "second run same day should not create a second file")
}

func TestDatabaseBackupPruning(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "geekbudget-backups")
	require.NoError(t, os.MkdirAll(backupDir, 0o750))

	// Create 5 fake archives with old dates
	fakeDates := []string{"2025-01-01", "2025-01-02", "2025-01-03", "2025-01-04", "2025-01-05"}
	for _, d := range fakeDates {
		path := filepath.Join(backupDir, backupPrefix+d+backupSuffix)
		require.NoError(t, os.WriteFile(path, []byte("fake"), 0o640))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{
		DBPath:         filepath.Join(tmpDir, "test.db"),
		BackupPath:     backupDir,
		BackupInterval: "24h",
		BackupMaxCount: 3,
	}
	task := newDatabaseBackupTask(logger, cfg)
	require.NoError(t, task.pruneBackups(backupDir))

	entries, err := os.ReadDir(backupDir)
	require.NoError(t, err)
	assert.Len(t, entries, 3, "should retain only the 3 newest backups")

	// Verify the 3 newest are kept
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	assert.Contains(t, names, backupPrefix+"2025-01-03"+backupSuffix)
	assert.Contains(t, names, backupPrefix+"2025-01-04"+backupSuffix)
	assert.Contains(t, names, backupPrefix+"2025-01-05"+backupSuffix)
}

// assertArchiveContains checks that all expected paths exist in the tar.gz archive.
func assertArchiveContains(t *testing.T, archivePath string, expectedPaths ...string) {
	t.Helper()

	f, err := os.Open(archivePath)
	require.NoError(t, err)
	defer f.Close() //nolint:errcheck

	gr, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gr.Close() //nolint:errcheck

	tr := tar.NewReader(gr)
	found := make(map[string]bool)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		found[hdr.Name] = true
	}

	for _, p := range expectedPaths {
		assert.True(t, found[p], "archive should contain %q; got %v", p, found)
	}
}

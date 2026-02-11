package database_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func TestImagesStorage(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	t.Run("Create and Get Image", func(t *testing.T) {
		data := []byte("fake image data")
		contentType := "image/png"

		img, err := st.CreateImage(data, contentType)
		if err != nil {
			t.Fatalf("failed to create image: %v", err)
		}

		if img.ID.String() == "" {
			t.Fatal("expected image ID to be set")
		}

		retrieved, err := st.GetImage(img.ID.String())
		if err != nil {
			t.Fatalf("failed to get image: %v", err)
		}

		if !bytes.Equal(retrieved.Data, data) {
			t.Errorf("expected data %v, got %v", data, retrieved.Data)
		}

		if retrieved.ContentType != contentType {
			t.Errorf("expected content type %s, got %s", contentType, retrieved.ContentType)
		}
	})

	t.Run("Get Non-existent Image", func(t *testing.T) {
		_, err := st.GetImage("00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Fatal("expected error for non-existent image, got nil")
		}

		if !errors.Is(err, database.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Delete Image", func(t *testing.T) {
		data := []byte("to be deleted")
		img, err := st.CreateImage(data, "image/jpeg")
		if err != nil {
			t.Fatalf("failed to create image: %v", err)
		}

		err = st.DeleteImage(img.ID.String())
		if err != nil {
			t.Fatalf("failed to delete image: %v", err)
		}

		_, err = st.GetImage(img.ID.String())
		if !errors.Is(err, database.ErrNotFound) {
			t.Errorf("expected ErrNotFound after deletion, got %v", err)
		}
	})
}

package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestNotificationStorage(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	t.Run("Create Notification with empty ID", func(t *testing.T) {
		n := &goserver.Notification{
			Date:        time.Now(),
			Type:        "error",
			Title:       "Test Error",
			Description: "Something went wrong",
		}

		created, err := st.CreateNotification(userID, n)
		if err != nil {
			t.Fatalf("failed to create notification: %v", err)
		}

		if created.Id == "" {
			t.Fatal("expected non-empty ID for created notification")
		}
		if created.Title != n.Title {
			t.Errorf("expected title %s, got %s", n.Title, created.Title)
		}
	})

	t.Run("Get Notifications", func(t *testing.T) {
		notifications, err := st.GetNotifications(userID)
		if err != nil {
			t.Fatalf("failed to get notifications: %v", err)
		}

		if len(notifications) == 0 {
			t.Fatal("expected at least one notification")
		}
	})

	t.Run("Delete Notification", func(t *testing.T) {
		notifications, _ := st.GetNotifications(userID)
		id := notifications[0].Id

		err := st.DeleteNotification(userID, id)
		if err != nil {
			t.Fatalf("failed to delete notification: %v", err)
		}

		// Verify it's gone
		remaining, _ := st.GetNotifications(userID)
		for _, n := range remaining {
			if n.Id == id {
				t.Fatalf("notification %s still exists after deletion", id)
			}
		}
	})
}

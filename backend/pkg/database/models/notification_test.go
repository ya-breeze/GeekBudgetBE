package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestNotificationToDB(t *testing.T) {
	t.Run("should parse valid UUID", func(t *testing.T) {
		id := uuid.New()
		gn := &goserver.Notification{
			Id:          id.String(),
			Date:        time.Now(),
			Type:        "info",
			Title:       "Test",
			Description: "Test Description",
		}

		mn, err := NotificationToDB(gn, "user1")
		assert.NoError(t, err)
		assert.Equal(t, id, mn.ID)
		assert.Equal(t, "user1", mn.UserID)
		assert.Equal(t, NotificationTypeInfo, mn.Type)
	})

	t.Run("should handle empty ID for creation", func(t *testing.T) {
		gn := &goserver.Notification{
			Id:          "",
			Date:        time.Now(),
			Type:        "error",
			Title:       "Error",
			Description: "Description",
		}

		mn, err := NotificationToDB(gn, "user1")
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, mn.ID)
		assert.Equal(t, NotificationTypeError, mn.Type)
	})

	t.Run("should return error for invalid UUID", func(t *testing.T) {
		gn := &goserver.Notification{
			Id: "invalid-uuid",
		}

		_, err := NotificationToDB(gn, "user1")
		assert.Error(t, err)
	})
}

func TestNotification_FromDB(t *testing.T) {
	id := uuid.New()
	mn := &Notification{
		ID:          id,
		UserID:      "user1",
		Date:        time.Now(),
		Type:        NotificationTypeError,
		Title:       "Title",
		Description: "Desc",
	}

	gn := mn.FromDB()
	assert.Equal(t, id.String(), gn.Id)
	assert.Equal(t, "error", gn.Type)
	assert.Equal(t, "Title", gn.Title)
}

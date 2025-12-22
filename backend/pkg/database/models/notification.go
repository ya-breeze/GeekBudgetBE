package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeOther              NotificationType = "other"
	NotificationTypeBalanceMatch       NotificationType = "balanceMatch"
	NotificationTypeBalanceDoesntMatch NotificationType = "balanceDoesntMatch"
	NotificationTypeError              NotificationType = "error"
	NotificationTypeInfo               NotificationType = "info"
)

type Notification struct {
	gorm.Model

	Date        time.Time
	Type        NotificationType
	URL         string
	Title       string
	Description string

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *Notification) FromDB() goserver.Notification {
	return goserver.Notification{
		Id:          t.ID.String(),
		Date:        t.Date,
		Type:        string(t.Type),
		Url:         t.URL,
		Title:       t.Title,
		Description: t.Description,
	}
}

func NotificationToDB(m *goserver.Notification, userID string) (*Notification, error) {
	var id uuid.UUID
	var err error
	if m.Id != "" {
		id, err = uuid.Parse(m.Id)
		if err != nil {
			return nil, err
		}
	}

	return &Notification{
		ID:     id,
		UserID: userID,

		Date:        m.Date,
		Type:        NotificationType(m.Type),
		URL:         m.Url,
		Title:       m.Title,
		Description: m.Description,
	}, nil
}

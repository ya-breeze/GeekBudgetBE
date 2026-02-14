package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func (s *storage) CreateNotification(userID string, notification *goserver.Notification) (goserver.Notification, error) {
	n, err := models.NotificationToDB(notification, userID)
	if err != nil {
		return goserver.Notification{}, fmt.Errorf(StorageError, err)
	}

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	if err := s.db.Create(n).Error; err != nil {
		return goserver.Notification{}, fmt.Errorf(StorageError, err)
	}

	return n.FromDB(), nil
}

func (s *storage) GetNotifications(userID string) ([]goserver.Notification, error) {
	result, err := s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Order("date DESC").Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	notifications := make([]goserver.Notification, 0)
	for result.Next() {
		var n models.Notification
		if err := s.db.ScanRows(result, &n); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		notifications = append(notifications, n.FromDB())
	}

	return notifications, nil
}

func (s *storage) DeleteNotification(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

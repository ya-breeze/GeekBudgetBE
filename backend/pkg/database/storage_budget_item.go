package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) CreateBudgetItem(userID string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error) {
	data := models.BudgetItemToDB(budgetItem, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "BudgetItem", data.ID.String(), "CREATED", nil, data); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	s.log.Info("BudgetItem created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) GetBudgetItems(userID string) ([]goserver.BudgetItem, error) {
	result, err := s.db.Model(&models.BudgetItem{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	items := make([]goserver.BudgetItem, 0)
	for result.Next() {
		var item models.BudgetItem
		if err := s.db.ScanRows(result, &item); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		items = append(items, item.FromDB())
	}

	return items, nil
}

func (s *storage) GetBudgetItem(userID string, id string) (goserver.BudgetItem, error) {
	var data models.BudgetItem
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BudgetItem{}, ErrNotFound
		}

		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) UpdateBudgetItem(
	userID string, id string, budgetItem *goserver.BudgetItemNoId,
) (goserver.BudgetItem, error) {
	return performUpdate[models.BudgetItem, goserver.BudgetItemNoIdInterface, goserver.BudgetItem](s, userID, "BudgetItem", id, budgetItem,
		models.BudgetItemToDB,
		func(m *models.BudgetItem) goserver.BudgetItem { return m.FromDB() },
		func(m *models.BudgetItem, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteBudgetItem(userID string, id string) error {
	var data models.BudgetItem
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err == nil {
		if err := s.recordAuditLog(s.db, userID, "BudgetItem", id, "DELETED", &data, nil); err != nil {
			s.log.Error("Failed to record audit log", "error", err)
		}
	}

	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.BudgetItem{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

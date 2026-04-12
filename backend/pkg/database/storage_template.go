package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) CreateTemplate(familyID uuid.UUID, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error) {
	tpl := models.TemplateToDB(t, familyID)
	tpl.ID = uuid.New()
	if err := s.db.Create(&tpl).Error; err != nil {
		return goserver.TransactionTemplate{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, familyID, "TransactionTemplate", tpl.ID.String(), "CREATED", nil, &tpl); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return tpl.FromDB(), nil
}

func (s *storage) GetTemplates(familyID uuid.UUID, accountID *string) ([]goserver.TransactionTemplate, error) {
	var records []models.TransactionTemplate
	if err := s.db.Where("family_id = ?", familyID).Order("name").Find(&records).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.TransactionTemplate, 0, len(records))
	for _, r := range records {
		if accountID != nil {
			matched := false
			for _, m := range r.Movements {
				if m.AccountId == *accountID {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		result = append(result, r.FromDB())
	}

	return result, nil
}

func (s *storage) UpdateTemplate(
	familyID uuid.UUID,
	id string,
	t *goserver.TransactionTemplateNoId,
) (goserver.TransactionTemplate, error) {
	return performUpdate[models.TransactionTemplate, *goserver.TransactionTemplateNoId, goserver.TransactionTemplate](
		s, familyID, "TransactionTemplate", id, t,
		func(t *goserver.TransactionTemplateNoId, familyID uuid.UUID) *models.TransactionTemplate {
			return models.TemplateToDB(t, familyID)
		},
		func(m *models.TransactionTemplate) goserver.TransactionTemplate { return m.FromDB() },
		func(m *models.TransactionTemplate, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteTemplate(familyID uuid.UUID, id string) error {
	var tpl models.TransactionTemplate
	if err := s.db.Where("id = ? AND family_id = ?", id, familyID).First(&tpl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, familyID, "TransactionTemplate", id, "DELETED", &tpl, nil); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	if err := s.db.Delete(&tpl).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

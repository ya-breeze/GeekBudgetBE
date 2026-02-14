package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"gorm.io/gorm"
)

func (s *storage) recordAuditLog(tx *gorm.DB, userID string, entityType string, entityID string, action string, snapshot interface{}) error {
	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal entity for audit log: %w", err)
	}

	changeSource := constants.ChangeSourceSystem
	if s.ctx != nil {
		if val, ok := s.ctx.Value(constants.ChangeSourceKey).(constants.ChangeSource); ok {
			changeSource = val
		}
	}

	auditLog := models.AuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		EntityType:   entityType,
		EntityID:     entityID,
		Action:       action,
		ChangeSource: string(changeSource),
		Snapshot:     string(jsonData),
		CreatedAt:    time.Now(),
	}

	return tx.Create(&auditLog).Error
}

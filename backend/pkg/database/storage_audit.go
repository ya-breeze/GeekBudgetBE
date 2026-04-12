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

func (s *storage) buildAuditLog(familyID uuid.UUID, entityType string, entityID string, action string, before interface{}, after interface{}) (models.AuditLog, error) {
	var beforeJSON, afterJSON *string

	if before != nil {
		b, err := json.Marshal(before)
		if err != nil {
			return models.AuditLog{}, fmt.Errorf("failed to marshal 'before' entity for audit log: %w", err)
		}
		str := string(b)
		beforeJSON = &str
	}

	if after != nil {
		b, err := json.Marshal(after)
		if err != nil {
			return models.AuditLog{}, fmt.Errorf("failed to marshal 'after' entity for audit log: %w", err)
		}
		str := string(b)
		afterJSON = &str
	}

	changeSource := constants.ChangeSourceSystem
	if s.ctx != nil {
		if val, ok := s.ctx.Value(constants.ChangeSourceKey).(constants.ChangeSource); ok {
			changeSource = val
		}
	}

	return models.AuditLog{
		ID:           uuid.New(),
		FamilyID:     familyID,
		EntityType:   entityType,
		EntityID:     entityID,
		Action:       action,
		ChangeSource: string(changeSource),
		Before:       beforeJSON,
		After:        afterJSON,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *storage) recordAuditLog(tx *gorm.DB, familyID uuid.UUID, entityType string, entityID string, action string, before interface{}, after interface{}) error {
	auditLog, err := s.buildAuditLog(familyID, entityType, entityID, action, before, after)
	if err != nil {
		return err
	}
	return tx.Create(&auditLog).Error
}

func (s *storage) GetAuditLogs(familyID uuid.UUID, filter AuditLogFilter) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := s.db.Where("family_id = ?", familyID)

	if filter.EntityType != nil {
		query = query.Where("entity_type = ?", *filter.EntityType)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at < ?", *filter.DateTo)
	}

	err := query.Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

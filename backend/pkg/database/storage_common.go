package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func performUpdate[M any, I any, O any](
	s *storage,
	userID string,
	entityType string,
	id string,
	input I,
	toDB func(I, string) *M,
	fromDB func(*M) O,
	setID func(*M, uuid.UUID),
) (O, error) {
	var empty O
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return empty, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var data *M
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return empty, ErrNotFound
		}

		return empty, fmt.Errorf(StorageError, err)
	}

	// Record audit log BEFORE update (with old state)
	if err := s.recordAuditLog(s.db, userID, entityType, id, "UPDATED", data); err != nil {
		s.log.Error("Failed to record audit log", "error", err, "entityType", entityType, "id", id)
	}

	data = toDB(input, userID)
	setID(data, idUUID)
	if err := s.db.Save(data).Error; err != nil {
		return empty, fmt.Errorf(StorageError, err)
	}

	return fromDB(data), nil
}

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type AuditLogsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAuditLogsAPIService(logger *slog.Logger, db database.Storage) goserver.AuditLogsAPIService {
	return &AuditLogsAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

func (s *AuditLogsAPIServiceImpl) GetAuditLogs(ctx context.Context, entityType string, entityId string, userId string, dateFrom time.Time, dateTo time.Time, limit int32, offset int32) (goserver.ImplResponse, error) {
	requestUserID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filter := database.AuditLogFilter{
		Limit:  int(limit),
		Offset: int(offset),
	}

	if entityType != "" {
		filter.EntityType = &entityType
	}
	if entityId != "" {
		filter.EntityID = &entityId
	}
	if !dateFrom.IsZero() {
		filter.DateFrom = &dateFrom
	}
	if !dateTo.IsZero() {
		filter.DateTo = &dateTo
	}

	logs, err := s.db.GetAuditLogs(requestUserID, filter)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get audit logs")
		return goserver.Response(500, nil), nil
	}

	apiLogs := make([]goserver.AuditLog, len(logs))
	for i, log := range logs {
		apiLogs[i] = goserver.AuditLog{
			Id:           log.ID.String(),
			UserId:       log.UserID,
			EntityType:   log.EntityType,
			EntityId:     log.EntityID,
			Action:       log.Action,
			ChangeSource: log.ChangeSource,
			Snapshot:     log.Snapshot,
			CreatedAt:    log.CreatedAt,
		}
	}

	return goserver.Response(200, apiLogs), nil
}

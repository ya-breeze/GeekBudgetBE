package api

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type NotificationsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewNotificationsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.NotificationsAPIServicer {
	return &NotificationsAPIServiceImpl{logger: logger, db: db}
}

func (s *NotificationsAPIServiceImpl) DeleteNotification(ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("DeleteNotification: UserID missing from context")
		return goserver.Response(500, nil), nil
	}

	err := s.db.DeleteNotification(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to delete notification")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *NotificationsAPIServiceImpl) GetNotifications(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("GetNotifications: UserID missing from context")
		return goserver.Response(500, nil), nil
	}

	notifications, err := s.db.GetNotifications(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get notifications")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, notifications), nil
}

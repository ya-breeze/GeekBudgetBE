package api

import (
	"context"
	"log/slog"

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

func (s *NotificationsAPIServiceImpl) DeleteNotification(context.Context, string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *NotificationsAPIServiceImpl) GetNotifications(context.Context) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

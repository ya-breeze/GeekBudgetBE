package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type TemplatesAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewTemplatesAPIServiceImpl(logger *slog.Logger, db database.Storage) *TemplatesAPIServiceImpl {
	return &TemplatesAPIServiceImpl{logger: logger, db: db}
}

func (s *TemplatesAPIServiceImpl) GetTemplates(ctx context.Context, accountId string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	var accountIDPtr *string
	if accountId != "" {
		accountIDPtr = &accountId
	}

	templates, err := s.db.GetTemplates(userID, accountIDPtr)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get templates")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	return goserver.Response(http.StatusOK, templates), nil
}

func (s *TemplatesAPIServiceImpl) CreateTemplate(
	ctx context.Context, t goserver.TransactionTemplateNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	if len(t.Movements) == 0 {
		return goserver.Response(http.StatusBadRequest, "movements must not be empty"), nil
	}

	result, err := s.db.CreateTemplate(userID, &t)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create template")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	return goserver.Response(http.StatusOK, result), nil
}

func (s *TemplatesAPIServiceImpl) UpdateTemplate(
	ctx context.Context, id string, t goserver.TransactionTemplateNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	result, err := s.db.UpdateTemplate(userID, id, &t)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update template")
		return mapErrorToResponse(err), nil
	}

	return goserver.Response(http.StatusOK, result), nil
}

func (s *TemplatesAPIServiceImpl) DeleteTemplate(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	if err := s.db.DeleteTemplate(userID, id); err != nil {
		s.logger.With("error", err).Error("Failed to delete template")
		return mapErrorToResponse(err), nil
	}

	return goserver.Response(http.StatusNoContent, nil), nil
}

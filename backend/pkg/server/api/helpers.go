package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func mapErrorToResponse(err error) goserver.ImplResponse {
	if err == database.ErrNotFound {
		return goserver.Response(http.StatusNotFound, nil)
	}
	return goserver.Response(http.StatusInternalServerError, nil)
}

func updateEntity[I any, O any](
	ctx context.Context,
	logger *slog.Logger,
	entityName string,
	id string,
	input I,
	updateFunc func(userID string, id string, input I) (O, error),
) (O, string, error) {
	var empty O
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		logger.Error("UserID not found in context")
		return empty, "", fmt.Errorf("UserID not found in context")
	}

	res, err := updateFunc(userID, id, input)
	if err != nil {
		if err == database.ErrNotFound {
			logger.With("error", err, "id", id, "userID", userID).Warn("Entity not found for update", "entity", entityName)
		} else {
			logger.With("error", err, "id", id, "userID", userID).Error("Failed to update entity", "entity", entityName)
		}
		return empty, userID, err
	}

	return res, userID, nil
}

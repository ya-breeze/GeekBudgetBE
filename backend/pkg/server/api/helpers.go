package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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
	updateFunc func(familyID uuid.UUID, id string, input I) (O, error),
) (O, uuid.UUID, error) {
	var empty O
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		logger.Error("FamilyID not found in context")
		return empty, uuid.UUID{}, fmt.Errorf("FamilyID not found in context")
	}

	res, err := updateFunc(familyID, id, input)
	if err != nil {
		if err == database.ErrNotFound {
			logger.With("error", err, "id", id, "familyID", familyID).Warn("Entity not found for update", "entity", entityName)
		} else {
			logger.With("error", err, "id", id, "familyID", familyID).Error("Failed to update entity", "entity", entityName)
		}
		return empty, familyID, err
	}

	return res, familyID, nil
}

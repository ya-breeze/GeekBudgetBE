package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type UserAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewUserAPIService(logger *slog.Logger, db database.Storage) goserver.UserAPIService {
	return &UserAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

// GetUser - return user object
func (s *UserAPIServiceImpl) GetUser(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	user, err := s.db.GetUser(userID)
	if err != nil && errors.Is(err, database.ErrNotFound) {
		return goserver.Response(500, nil), nil
	}
	if user == nil {
		return goserver.Response(404, nil), nil
	}

	return goserver.Response(200, user.FromDB()), nil
}

// UpdateUserFavoriteCurrency - update user's favorite currency
func (s *UserAPIServiceImpl) UpdateUserFavoriteCurrency(
	ctx context.Context, body goserver.UserPatchBody,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	user, err := s.db.GetUser(userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return goserver.Response(404, nil), nil
		}

		s.logger.With("error", err).Error("Failed to get user")
		return goserver.Response(500, nil), nil
	}

	// Standard PATCH logic: take the value from the request body as-is.
	// Sending a non-empty favoriteCurrencyId sets the favorite currency;
	// sending an empty string clears it.
	user.FavoriteCurrencyID = body.FavoriteCurrencyId

	if err := s.db.PutUser(user); err != nil {
		s.logger.With("error", err).Error("Failed to update user")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, user.FromDB()), nil
}

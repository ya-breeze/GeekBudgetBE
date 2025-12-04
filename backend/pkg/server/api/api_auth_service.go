package api

import (
	"context"
	"encoding/base64"
	"errors"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type AuthAPIService struct {
	logger    *slog.Logger
	db        database.Storage
	issuer    string
	jwtSecret string
}

func NewAuthAPIService(logger *slog.Logger, db database.Storage, issuer, jwtSecret string) goserver.AuthAPIService {
	return &AuthAPIService{
		logger:    logger,
		db:        db,
		issuer:    issuer,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthAPIService) Authorize(_ context.Context, authData goserver.AuthData) (goserver.ImplResponse, error) {
	userID, err := s.db.GetUserID(authData.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		s.logger.Warn("failed to get user ID", "username", authData.Email)
		return goserver.Response(500, nil), nil // TODO internal error
	}
	user, err := s.db.GetUser(userID)
	if err != nil {
		s.logger.Warn("failed to get user", "ID", userID)
		return goserver.Response(500, nil), nil // TODO internal error
	}
	if user == nil {
		s.logger.Warn("user not found", "ID", userID)
		return goserver.Response(401, nil), nil
	}

	hashed, err := base64.StdEncoding.DecodeString(user.HashedPassword)
	if err != nil {
		return goserver.Response(500, nil), nil // TODO internal error
	}
	if !auth.CheckPasswordHash([]byte(authData.Password), hashed) {
		return goserver.Response(401, nil), nil
	}

	token, err := auth.CreateJWT(userID, s.issuer, s.jwtSecret)
	if err != nil {
		return goserver.Response(500, nil), nil // TODO internal error
	}
	return goserver.Response(200, goserver.Authorize200Response{Token: token}), nil
}

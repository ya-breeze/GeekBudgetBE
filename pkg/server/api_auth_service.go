package server

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
	jwtSecret string
}

func NewAuthAPIService(logger *slog.Logger, db database.Storage, jwtSecret string) goserver.AuthAPIService {
	return &AuthAPIService{
		logger:    logger,
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthAPIService) Authorize(_ context.Context, authData goserver.AuthData) (goserver.ImplResponse, error) {
	user, err := s.db.GetUser(authData.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return goserver.Response(500, nil), nil // TODO internal error
	}
	if user == nil {
		return goserver.Response(401, nil), nil
	}

	hashed, err := base64.StdEncoding.DecodeString(user.HashedPassword)
	if err != nil {
		return goserver.Response(500, nil), nil // TODO internal error
	}
	if !auth.CheckPasswordHash([]byte(authData.Password), hashed) {
		return goserver.Response(401, nil), nil
	}

	token, err := auth.CreateJWT(user.Login, s.jwtSecret)
	if err != nil {
		return goserver.Response(500, nil), nil // TODO internal error
	}
	return goserver.Response(200, goserver.Authorize200Response{Token: token}), nil
}

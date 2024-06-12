package server

import (
	"context"
	"encoding/base64"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type AuthAPIService struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAuthAPIService(logger *slog.Logger, db database.Storage) goserver.AuthAPIService {
	return &AuthAPIService{
		logger: logger,
		db:     db,
	}
}

func (s *AuthAPIService) Authorize(_ context.Context, authData goserver.AuthData) (goserver.ImplResponse, error) {
	user, err := s.db.GetUser(authData.Email)
	if err != nil {
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

	return goserver.Response(200, goserver.Authorize200Response{Token: "token"}), nil
}

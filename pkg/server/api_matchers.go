package server

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type MatchersAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewMatchersAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.MatchersAPIServicer {
	return &MatchersAPIServiceImpl{logger: logger, db: db}
}

func (s *MatchersAPIServiceImpl) CheckMatcher(ctx context.Context, r goserver.CheckMatcherRequest,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *MatchersAPIServiceImpl) GetMatchers(ctx context.Context) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *MatchersAPIServiceImpl) CreateMatcher(context.Context, goserver.MatcherNoId) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *MatchersAPIServiceImpl) DeleteMatcher(context.Context, string) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *MatchersAPIServiceImpl) UpdateMatcher(context.Context, string, goserver.MatcherNoId,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

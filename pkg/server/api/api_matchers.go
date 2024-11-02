package api

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
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
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.GetMatchers(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

func (s *MatchersAPIServiceImpl) CreateMatcher(ctx context.Context, m goserver.MatcherNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.CreateMatcher(userID, &m)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create matcher")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

func (s *MatchersAPIServiceImpl) DeleteMatcher(context.Context, string) (goserver.ImplResponse, error) {
	return goserver.Response(500, nil), nil
}

func (s *MatchersAPIServiceImpl) UpdateMatcher(ctx context.Context, id string, m goserver.MatcherNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.UpdateMatcher(userID, id, &m)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update matcher")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

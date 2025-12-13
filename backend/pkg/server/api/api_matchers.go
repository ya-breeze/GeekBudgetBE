package api

import (
	"context"
	"log/slog"
	"regexp"

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
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Create MatcherRuntime from the request matcher data
	matcherRuntime, err := s.db.CreateMatcherRuntimeFromNoId(&r.Matcher)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create matcher runtime")
		return goserver.Response(400, nil), nil
	}

	// Convert TransactionNoId to Transaction for matching
	transaction := goserver.Transaction{
		Date:               r.Transaction.GetDate(),
		Description:        r.Transaction.GetDescription(),
		Place:              r.Transaction.GetPlace(),
		Tags:               r.Transaction.GetTags(),
		PartnerName:        r.Transaction.GetPartnerName(),
		PartnerAccount:     r.Transaction.GetPartnerAccount(),
		PartnerInternalId:  r.Transaction.GetPartnerInternalId(),
		Extra:              r.Transaction.GetExtra(),
		UnprocessedSources: r.Transaction.GetUnprocessedSources(),
		ExternalIds:        r.Transaction.GetExternalIds(),
		Movements:          r.Transaction.GetMovements(),
	}

	// Perform the match with detailed results
	matchDetails := common.MatchWithDetails(&matcherRuntime, &transaction)

	// Return the result
	response := goserver.CheckMatcher200Response{
		Result: matchDetails.Matched,
		Reason: matchDetails.FailureReason,
	}

	s.logger.With("userID", userID).With("matched", matchDetails.Matched).
		With("reason", matchDetails.FailureReason).Info("CheckMatcher result")

	return goserver.Response(200, response), nil
}

func (s *MatchersAPIServiceImpl) CheckRegex(ctx context.Context, r goserver.CheckRegexRequest,
) (goserver.ImplResponse, error) {
	regexStr := r.GetRegex()
	testStr := r.GetTestString()

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return goserver.Response(200, goserver.CheckRegex200Response{
			IsValid: false,
			IsMatch: false,
			Error:   err.Error(),
		}), nil
	}

	match := re.MatchString(testStr)

	return goserver.Response(200, goserver.CheckRegex200Response{
		IsValid: true,
		IsMatch: match,
	}), nil
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

func (s *MatchersAPIServiceImpl) GetMatcher(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.GetMatcher(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matcher")
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

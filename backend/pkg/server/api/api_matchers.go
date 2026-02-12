package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type MatchersAPIServiceImpl struct {
	logger             *slog.Logger
	db                 database.Storage
	cfg                *config.Config
	unprocessedService *UnprocessedTransactionsAPIServiceImpl
}

func NewMatchersAPIServiceImpl(
	logger *slog.Logger,
	db database.Storage,
	cfg *config.Config,
	unprocessedService *UnprocessedTransactionsAPIServiceImpl,
) goserver.MatchersAPIServicer {
	return &MatchersAPIServiceImpl{
		logger:             logger,
		db:                 db,
		cfg:                cfg,
		unprocessedService: unprocessedService,
	}
}

func (s *MatchersAPIServiceImpl) CheckMatcher(ctx context.Context, r goserver.CheckMatcherRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	// Create MatcherRuntime from the request matcher data
	matcherRuntime, err := s.db.CreateMatcherRuntimeFromNoId(&r.Matcher)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create matcher runtime")
		return goserver.Response(400, "failed to create matcher runtime"), nil
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
		s.logger.Error("UserID not found in context")
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
		s.logger.Error("UserID not found in context")
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
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.CreateMatcher(userID, &m)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create matcher")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

func (s *MatchersAPIServiceImpl) DeleteMatcher(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	if err := s.db.DeleteMatcher(userID, id); err != nil {
		s.logger.With("error", err, "matcherID", id).Error("Failed to delete matcher")
		return goserver.Response(500, nil), nil
	}

	s.logger.With("matcherID", id).Info("Matcher deleted")
	return goserver.Response(204, nil), nil
}

func (s *MatchersAPIServiceImpl) UpdateMatcher(ctx context.Context, id string, m goserver.MatcherNoId,
) (goserver.ImplResponse, error) {
	res, userID, err := updateEntity[goserver.MatcherNoIdInterface, goserver.Matcher](ctx, s.logger, "matcher", id, &m, s.db.UpdateMatcher)
	if err != nil {
		return mapErrorToResponse(err), nil
	}

	var autoProcessedIds []string
	// Check if this matcher is now "perfect" and can auto-process existing transactions
	autoIds, err := s.unprocessedService.ProcessUnprocessedTransactionsAgainstMatcher(ctx, userID, id, "")
	if err != nil {
		// Log but don't fail the request
		s.logger.With("error", err).Error("Failed to process unprocessed transactions against updated matcher")
	} else {
		autoProcessedIds = autoIds
	}

	return goserver.Response(http.StatusOK, goserver.UpdateMatcher200Response{
		Matcher:          res,
		AutoProcessedIds: autoProcessedIds,
	}), nil
}

func (s *MatchersAPIServiceImpl) UploadMatcherImage(
	ctx context.Context, matcherID string, file *os.File,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	matcher, err := s.db.GetMatcher(userID, matcherID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matcher")
		return goserver.Response(404, nil), nil
	}

	// The generated code created a temp file and closed it. We need to re-open it.
	// We also need to remove it after we are done.
	defer os.Remove(file.Name())

	f, err := os.Open(file.Name())
	if err != nil {
		s.logger.With("error", err).Error("Failed to open temp file")
		return goserver.Response(500, nil), nil
	}
	defer f.Close()

	// Read file content
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		s.logger.With("error", err).Error("Failed to read file content")
		return goserver.Response(500, nil), nil
	}

	// Detect content type
	contentType := http.DetectContentType(fileBytes)

	// Create image in DB
	image, err := s.db.CreateImage(fileBytes, contentType)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create image in DB")
		return goserver.Response(500, nil), nil
	}

	// Delete old image if exists
	if matcher.Image != "" {
		if err := s.db.DeleteImage(matcher.Image); err != nil {
			s.logger.With("error", err, "imageID", matcher.Image).Warn("Failed to delete old image")
		}
	}

	// Update DB
	mNoID := models.MatcherWithoutID(&matcher)
	mNoID.Image = image.ID.String()
	updatedMatcher, err := s.db.UpdateMatcher(userID, matcherID, mNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update matcher with image")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, updatedMatcher), nil
}

func (s *MatchersAPIServiceImpl) DeleteMatcherImage(
	ctx context.Context, matcherID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	matcher, err := s.db.GetMatcher(userID, matcherID)
	if err != nil {
		return goserver.Response(404, nil), nil
	}

	if matcher.Image != "" {
		if err := s.db.DeleteImage(matcher.Image); err != nil {
			s.logger.With("error", err).Error("Failed to delete image from DB")
			return goserver.Response(500, nil), nil
		}

		mNoID := models.MatcherWithoutID(&matcher)
		mNoID.Image = ""
		updatedMatcher, err := s.db.UpdateMatcher(userID, matcherID, mNoID)
		if err != nil {
			s.logger.With("error", err).Error("Failed to update matcher (remove image)")
			return goserver.Response(500, nil), nil
		}
		return goserver.Response(200, updatedMatcher), nil
	}

	return goserver.Response(200, matcher), nil
}

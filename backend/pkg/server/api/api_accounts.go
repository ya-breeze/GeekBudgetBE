package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type AccountsAPIServicerImpl struct {
	logger *slog.Logger
	db     database.Storage
	cfg    *config.Config
}

func NewAccountsAPIService(logger *slog.Logger, db database.Storage, cfg *config.Config) goserver.AccountsAPIServicer {
	return &AccountsAPIServicerImpl{
		logger: logger,
		db:     db,
		cfg:    cfg,
	}
}

func (s *AccountsAPIServicerImpl) CreateAccount(
	ctx context.Context, acc goserver.AccountNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.CreateAccount(userID, &acc)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}

func (s *AccountsAPIServicerImpl) GetAccounts(
	ctx context.Context,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, accounts), nil
}

func (s *AccountsAPIServicerImpl) UpdateAccount(
	ctx context.Context, accountID string, acc goserver.AccountNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.UpdateAccount(userID, accountID, &acc)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}

func (s *AccountsAPIServicerImpl) DeleteAccount(
	ctx context.Context, accountID string, replaceWithAccountId string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Delete image if exists
	account, err := s.db.GetAccount(userID, accountID)
	if err == nil && account.Image != "" {
		if err := s.db.DeleteImage(account.Image); err != nil {
			s.logger.With("error", err, "imageID", account.Image).Warn("Failed to delete account image")
		}
	}

	if replaceWithAccountId != "" {
		if _, err := s.db.GetAccount(userID, replaceWithAccountId); err != nil {
			s.logger.With("error", err, "replaceWithAccountId", replaceWithAccountId).Warn("Replacement account not found")
			return goserver.Response(400, nil), nil
		}
	}

	if err := s.db.DeleteAccount(userID, accountID, &replaceWithAccountId); err != nil {
		if errors.Is(err, database.ErrAccountInUse) {
			s.logger.With("error", err).Warn("Cannot delete account in use without replacement")
			return goserver.Response(400, nil), nil
		}
		s.logger.With("error", err).Error("Failed to delete account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *AccountsAPIServicerImpl) GetAccountHistory(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	history, err := s.db.GetAccountHistory(userID, accountID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get account history")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, history), nil
}

func (s *AccountsAPIServicerImpl) GetAccount(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.GetAccount(userID, accountID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}

func (s *AccountsAPIServicerImpl) UploadAccountImage(
	ctx context.Context, accountID string, file *os.File,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Validation: Check if account exists and belongs to user
	account, err := s.db.GetAccount(userID, accountID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get account")
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
	if account.Image != "" {
		if err := s.db.DeleteImage(account.Image); err != nil {
			s.logger.With("error", err, "imageID", account.Image).Warn("Failed to delete old image")
			// Continue execution, not critical
		}
	}

	// Update DB
	accNoID := models.AccountWithoutID(&account)
	accNoID.Image = image.ID.String()
	updatedAccount, err := s.db.UpdateAccount(userID, accountID, accNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update account with image")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, updatedAccount), nil
}

func (s *AccountsAPIServicerImpl) DeleteAccountImage(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.GetAccount(userID, accountID)
	if err != nil {
		return goserver.Response(404, nil), nil
	}

	if account.Image != "" {
		if err := s.db.DeleteImage(account.Image); err != nil {
			s.logger.With("error", err).Error("Failed to delete image from DB")
			return goserver.Response(500, nil), nil
		}

		accNoID := models.AccountWithoutID(&account)
		accNoID.Image = ""
		updatedAccount, err := s.db.UpdateAccount(userID, accountID, accNoID)
		if err != nil {
			s.logger.With("error", err).Error("Failed to update account (remove image)")
			return goserver.Response(500, nil), nil
		}
		return goserver.Response(200, updatedAccount), nil
	}

	return goserver.Response(200, account), nil
}

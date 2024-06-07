package database

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

var ErrStorageError = errors.New("storage error")

type Storage interface {
	Open() error
	Close() error

	CreateAccount(userId string, account *goserver.AccountNoId) (goserver.Account, error)
	GetAccounts(userId string) ([]goserver.Account, error)
}

type storage struct {
	log *slog.Logger
	db  *gorm.DB
}

func NewStorage(logger *slog.Logger, cfg *config.Config) Storage {
	return &storage{log: logger}
}

func (s *storage) Open() error {
	var err error
	s.db, err = openSqlite()
	if err != nil {
		panic("failed to connect database")
	}
	if err := migrate(s.db); err != nil {
		panic("failed to migrate database")
	}

	return nil
}

func (s *storage) Close() error {
	// return s.db.Close()
	return nil
}

func (s *storage) GetAccounts(userId string) ([]goserver.Account, error) {
	result, err := s.db.Model(&models.Account{}).Where("user_id = ?", userId).Rows()
	if err != nil {
		return nil, ErrStorageError
	}
	defer result.Close()

	var accounts []goserver.Account
	for result.Next() {
		var acc models.Account
		if err := s.db.ScanRows(result, &acc); err != nil {
			return nil, ErrStorageError
		}

		accounts = append(accounts, acc.FromDb())
	}

	return accounts, nil
}

func (s *storage) CreateAccount(userId string, account *goserver.AccountNoId) (goserver.Account, error) {
	acc := models.Account{
		ID:          uuid.New(),
		UserId:      userId,
		AccountNoId: *account,
	}
	if err := s.db.Create(&acc).Error; err != nil {
		return goserver.Account{}, ErrStorageError
	}

	return acc.FromDb(), nil
}

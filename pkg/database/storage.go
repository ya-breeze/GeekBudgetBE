package database

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

const ErrStorageError = "storage error: %w"

type Storage interface {
	Open() error
	Close() error

	CreateUser(username, password string) error

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
		return nil, fmt.Errorf(ErrStorageError, err)
	}
	defer result.Close()

	var accounts []goserver.Account
	for result.Next() {
		var acc models.Account
		if err := s.db.ScanRows(result, &acc); err != nil {
			return nil, fmt.Errorf(ErrStorageError, err)
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
		return goserver.Account{}, fmt.Errorf(ErrStorageError, err)
	}

	return acc.FromDb(), nil
}

func (s *storage) CreateUser(username, hashedPassword string) error {
	user := models.User{
		ID:             uuid.New(),
		Email:          username,
		HashedPassword: hashedPassword,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return fmt.Errorf(ErrStorageError, err)
	}

	return nil
}

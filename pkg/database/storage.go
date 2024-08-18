package database

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

const StorageError = "storage error: %w"

var ErrNotFound = errors.New("not found")

//nolint:interfacebloat
type Storage interface {
	Open() error
	Close() error

	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) error
	PutUser(user *models.User) error

	CreateAccount(userID string, account *goserver.AccountNoId) (goserver.Account, error)
	GetAccounts(userID string) ([]goserver.Account, error)
	GetAccount(userID string, id string) (goserver.Account, error)
	UpdateAccount(userID string, id string, account *goserver.AccountNoId) (goserver.Account, error)
	DeleteAccount(userID string, id string) error
	GetAccountHistory(userID string, accountID string) ([]goserver.Transaction, error)

	CreateCurrency(userID string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	GetCurrencies(userID string) ([]goserver.Currency, error)
	GetCurrency(userID string, id string) (goserver.Currency, error)
	UpdateCurrency(userID string, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	DeleteCurrency(userID string, id string) error

	GetTransactions(userID string) ([]goserver.Transaction, error)
	CreateTransaction(userID string, transaction *goserver.TransactionNoId) (goserver.Transaction, error)
	UpdateTransaction(userID string, id string, transaction *goserver.TransactionNoId) (goserver.Transaction, error)
	DeleteTransaction(userID string, id string) error
	GetTransaction(userID string, id string) (goserver.Transaction, error)
}

type storage struct {
	log *slog.Logger
	cfg *config.Config
	db  *gorm.DB
}

func NewStorage(logger *slog.Logger, cfg *config.Config) Storage {
	return &storage{log: logger, db: nil, cfg: cfg}
}

func (s *storage) Open() error {
	var err error
	s.db, err = openSqlite(s.log, s.cfg.DBPath, s.cfg.Verbose)
	if err != nil {
		s.log.Error("failed to connect database", "error", err)
		panic("failed to connect database")
	}
	if err := migrate(s.db); err != nil {
		s.log.Error("failed to migrate database", "error", err)
		panic("failed to migrate database")
	}

	return nil
}

func (s *storage) Close() error {
	// return s.db.Close()
	return nil
}

func (s *storage) GetAccounts(userID string) ([]goserver.Account, error) {
	result, err := s.db.Model(&models.Account{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	accounts := make([]goserver.Account, 0)
	for result.Next() {
		var acc models.Account
		if err := s.db.ScanRows(result, &acc); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		accounts = append(accounts, acc.FromDB())
	}

	return accounts, nil
}

func (s *storage) CreateAccount(userID string, account *goserver.AccountNoId) (goserver.Account, error) {
	acc := models.Account{
		ID:          uuid.New(),
		UserID:      userID,
		AccountNoId: *account,
	}
	if err := s.db.Create(&acc).Error; err != nil {
		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	return acc.FromDB(), nil
}

func (s *storage) CreateUser(username, hashedPassword string) error {
	user := models.User{
		ID:             uuid.New(),
		Login:          username,
		HashedPassword: hashedPassword,
		StartDate:      time.Now(),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetUser(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf(StorageError, err)
	}

	return &user, nil
}

func (s *storage) PutUser(user *models.User) error {
	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetUserID(username string) (string, error) {
	var user models.User
	if err := s.db.Where("login = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf(StorageError, err)
	}

	return user.ID.String(), nil
}

func (s *storage) UpdateAccount(userID string, id string, account *goserver.AccountNoId) (goserver.Account, error) {
	var acc models.Account
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Account{}, ErrNotFound
		}

		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	acc.AccountNoId = *account
	if err := s.db.Save(&acc).Error; err != nil {
		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	return acc.FromDB(), nil
}

func (s *storage) DeleteAccount(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetAccountHistory(userID string, accountID string) ([]goserver.Transaction, error) {
	// result, err := s.db.Model(&models.Transaction{}).Where("user_id = ? AND account_id = ?", userID, accountID).Rows()
	// if err != nil {
	// 	return nil, fmt.Errorf(StorageError, err)
	// }
	// defer result.Close()

	var transactions []goserver.Transaction
	// for result.Next() {
	// 	var tr models.Transaction
	// 	if err := s.db.ScanRows(result, &tr); err != nil {
	// 		return nil, fmt.Errorf(StorageError, err)
	// 	}

	// 	transactions = append(transactions, tr.FromDB())
	// }

	return transactions, nil
}

func (s *storage) GetAccount(userID string, id string) (goserver.Account, error) {
	var acc models.Account
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Account{}, ErrNotFound
		}

		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	return acc.FromDB(), nil
}

func (s *storage) CreateCurrency(userID string, currency *goserver.CurrencyNoId) (goserver.Currency, error) {
	cur := models.Currency{
		ID:           uuid.New(),
		UserID:       userID,
		CurrencyNoId: *currency,
	}
	if err := s.db.Create(&cur).Error; err != nil {
		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	return cur.FromDB(), nil
}

func (s *storage) GetCurrencies(userID string) ([]goserver.Currency, error) {
	result, err := s.db.Model(&models.Currency{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	currencies := make([]goserver.Currency, 0)
	for result.Next() {
		var cur models.Currency
		if err := s.db.ScanRows(result, &cur); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		currencies = append(currencies, cur.FromDB())
	}

	return currencies, nil
}

func (s *storage) GetCurrency(userID string, id string) (goserver.Currency, error) {
	var cur models.Currency
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&cur).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Currency{}, ErrNotFound
		}

		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	return cur.FromDB(), nil
}

func (s *storage) UpdateCurrency(userID string, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error) {
	var cur models.Currency
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&cur).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Currency{}, ErrNotFound
		}

		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	cur.CurrencyNoId = *currency
	if err := s.db.Save(&cur).Error; err != nil {
		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	return cur.FromDB(), nil
}

func (s *storage) DeleteCurrency(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Currency{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetTransactions(userID string) ([]goserver.Transaction, error) {
	result, err := s.db.Model(&models.Transaction{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	transactions := make([]goserver.Transaction, 0)
	for result.Next() {
		var tr models.Transaction
		if err := s.db.ScanRows(result, &tr); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		transactions = append(transactions, tr.FromDB())
	}

	return transactions, nil
}

func (s *storage) CreateTransaction(userID string, input *goserver.TransactionNoId,
) (goserver.Transaction, error) {
	t := models.TransactionToDB(input, userID)
	t.ID = uuid.New()
	if err := s.db.Create(t).Error; err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}
	s.log.Info("Transaction created", "id", t.ID)

	return t.FromDB(), nil
}

func (s *storage) UpdateTransaction(userID string, id string, input *goserver.TransactionNoId,
) (goserver.Transaction, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var t *models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Transaction{}, ErrNotFound
		}

		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	t = models.TransactionToDB(input, userID)
	t.ID = idUUID
	if err := s.db.Save(&t).Error; err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	return t.FromDB(), nil
}

func (s *storage) DeleteTransaction(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetTransaction(userID string, id string) (goserver.Transaction, error) {
	var transaction models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Transaction{}, ErrNotFound
		}

		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	return transaction.FromDB(), nil
}

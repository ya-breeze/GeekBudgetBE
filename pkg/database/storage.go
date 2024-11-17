package database

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

const StorageError = "storage error: %w"

var ErrNotFound = errors.New("not found")

type ImportInfo struct {
	UserID           string
	BankImporterID   string
	BankImporterType string
}

//nolint:interfacebloat
type Storage interface {
	Open() error
	Close() error

	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
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

	GetTransactions(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error)
	CreateTransaction(userID string, transaction *goserver.TransactionNoId) (goserver.Transaction, error)
	UpdateTransaction(
		userID string, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	DeleteTransaction(userID string, id string) error
	GetTransaction(userID string, id string) (goserver.Transaction, error)

	GetBankImporters(userID string) ([]goserver.BankImporter, error)
	CreateBankImporter(userID string, bankImporter *goserver.BankImporterNoId) (goserver.BankImporter, error)
	UpdateBankImporter(userID string, id string, bankImporter goserver.BankImporterNoIdInterface,
	) (goserver.BankImporter, error)
	DeleteBankImporter(userID string, id string) error
	GetBankImporter(userID string, id string) (goserver.BankImporter, error)
	GetAllBankImporters() ([]ImportInfo, error)

	GetMatchers(userID string) ([]goserver.Matcher, error)
	GetMatcher(userID string, id string) (goserver.Matcher, error)
	GetMatchersRuntime(userID string) ([]MatcherRuntime, error)
	CreateMatcher(userID string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	UpdateMatcher(userID string, id string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	DeleteMatcher(userID string, id string) error
}

type MatcherRuntime struct {
	Matcher           *goserver.Matcher
	DescriptionRegexp *regexp.Regexp
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
	s.log.Info("Opening database", "path", s.cfg.DBPath)
	var err error
	s.db, err = openSqlite(s.log, s.cfg.DBPath, s.cfg.Verbose)
	if err != nil {
		s.log.Error("failed to connect database", "error", err)
		panic("failed to connect database")
	}
	if err := autoMigrateModels(s.db); err != nil {
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
	result, err := s.db.Model(&models.Account{}).Where("user_id = ?", userID).Order("type, name").Rows()
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
	acc := models.AccountToDB(account, userID)
	acc.ID = uuid.New()
	if err := s.db.Create(&acc).Error; err != nil {
		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	return acc.FromDB(), nil
}

func (s *storage) CreateUser(username, hashedPassword string) (*models.User, error) {
	user := models.User{
		ID:             uuid.New(),
		Login:          username,
		HashedPassword: hashedPassword,
		StartDate:      time.Now(),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	return &user, nil
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

	accID := acc.ID
	acc = *models.AccountToDB(account, userID)
	acc.ID = accID
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

func (s *storage) GetTransactions(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error) {
	req := s.db.Model(&models.Transaction{}).Where("user_id = ?", userID)
	if !dateFrom.IsZero() {
		req = req.Where("date >= ?", dateFrom)
	}
	if !dateTo.IsZero() {
		req = req.Where("date < ?", dateTo)
	}
	req = req.Order("date")

	result, err := req.Rows()
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

//nolint:dupl // TODO: refactor
func (s *storage) UpdateTransaction(userID string, id string, input goserver.TransactionNoIdInterface,
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

// #region BankImporters
func (s *storage) GetBankImporters(userID string) ([]goserver.BankImporter, error) {
	result, err := s.db.Model(&models.BankImporter{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	importers := make([]goserver.BankImporter, 0)
	for result.Next() {
		var imp models.BankImporter
		if err := s.db.ScanRows(result, &imp); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		importers = append(importers, imp.FromDB())
	}

	return importers, nil
}

func (s *storage) CreateBankImporter(userID string, bankImporter *goserver.BankImporterNoId,
) (goserver.BankImporter, error) {
	data := models.BankImporterToDB(bankImporter, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}
	s.log.Info("BankImporter created", "id", data.ID)

	return data.FromDB(), nil
}

//nolint:dupl // TODO: refactor
func (s *storage) UpdateBankImporter(userID string, id string, bankImporter goserver.BankImporterNoIdInterface,
) (goserver.BankImporter, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return goserver.BankImporter{}, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var data *models.BankImporter
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BankImporter{}, ErrNotFound
		}

		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}

	data = models.BankImporterToDB(bankImporter, userID)
	data.ID = idUUID
	if err := s.db.Save(&data).Error; err != nil {
		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) DeleteBankImporter(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.BankImporter{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) GetBankImporter(userID string, id string) (goserver.BankImporter, error) {
	var data models.BankImporter
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BankImporter{}, ErrNotFound
		}

		return goserver.BankImporter{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) GetAllBankImporters() ([]ImportInfo, error) {
	result, err := s.db.Model(&models.BankImporter{}).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	importers := make([]ImportInfo, 0)
	for result.Next() {
		var imp models.BankImporter
		if err := s.db.ScanRows(result, &imp); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		importers = append(importers, ImportInfo{
			UserID:           imp.UserID,
			BankImporterID:   imp.ID.String(),
			BankImporterType: imp.Type,
		})
	}

	return importers, nil
}

// #endregion BankImporters

// #region Matchers
func (s *storage) GetMatchers(userID string) ([]goserver.Matcher, error) {
	result, err := s.db.Model(&models.Matcher{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	matchers := make([]goserver.Matcher, 0)
	for result.Next() {
		var m models.Matcher
		if err := s.db.ScanRows(result, &m); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		matchers = append(matchers, m.FromDB())
	}

	return matchers, nil
}

func (s *storage) CreateMatcher(userID string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error) {
	data := models.MatcherToDB(matcher, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}
	s.log.Info("Matcher created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) GetMatchersRuntime(userID string) ([]MatcherRuntime, error) {
	matchers, err := s.GetMatchers(userID)
	if err != nil {
		return nil, err
	}

	res := make([]MatcherRuntime, 0, len(matchers))
	for _, m := range matchers {
		runtime := MatcherRuntime{Matcher: &m}
		if m.DescriptionRegExp != "" {
			r, err := regexp.Compile(m.DescriptionRegExp)
			if err != nil {
				return nil, fmt.Errorf("failed to compile regexp: %w", err)
			}
			runtime.DescriptionRegexp = r
		}

		res = append(res, runtime)
	}

	return res, nil
}

func (s *storage) UpdateMatcher(userID string, id string, matcher goserver.MatcherNoIdInterface,
) (goserver.Matcher, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return goserver.Matcher{}, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var data *models.Matcher
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Matcher{}, ErrNotFound
		}

		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}

	data = models.MatcherToDB(matcher, userID)
	data.ID = idUUID
	if err := s.db.Save(&data).Error; err != nil {
		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) GetMatcher(userID string, id string) (goserver.Matcher, error) {
	var data models.Matcher
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Matcher{}, ErrNotFound
		}

		return goserver.Matcher{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) DeleteMatcher(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Matcher{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

//#endregion Matchers

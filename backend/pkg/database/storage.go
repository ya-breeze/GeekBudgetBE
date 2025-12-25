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

//go:generate go tool github.com/golang/mock/mockgen -destination=mocks/mock_storage.go -package=mocks github.com/ya-breeze/geekbudgetbe/pkg/database Storage

const StorageError = "storage error: %w"

var (
	ErrNotFound                           = errors.New("not found")
	ErrAccountInUse                       = errors.New("account is in use")
	ErrImportedTransactionCannotBeDeleted = errors.New("imported transaction cannot be deleted")
	ErrCurrencyInUse                      = errors.New("currency is in use")
)

type ImportInfo struct {
	UserID           string
	BankImporterID   string
	BankImporterType string
	FetchAll         bool
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
	DeleteAccount(userID string, id string, replaceWithAccountID *string) error
	GetAccountHistory(userID string, accountID string) ([]goserver.Transaction, error)

	CreateCurrency(userID string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	GetCurrencies(userID string) ([]goserver.Currency, error)
	GetCurrency(userID string, id string) (goserver.Currency, error)
	UpdateCurrency(userID string, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	DeleteCurrency(userID string, id string, replaceWithCurrencyID *string) error

	GetTransactions(userID string, dateFrom, dateTo time.Time, onlySuspicious bool) ([]goserver.Transaction, error)
	CreateTransaction(userID string, transaction goserver.TransactionNoIdInterface) (goserver.Transaction, error)
	UpdateTransaction(
		userID string, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	DeleteTransaction(userID string, id string) error
	DeleteDuplicateTransaction(userID string, id, duplicateID string) error
	GetTransaction(userID string, id string) (goserver.Transaction, error)
	GetTransactionsIncludingDeleted(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error)
	GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error)
	UnmergeTransaction(userID, id string) error

	GetBankImporters(userID string) ([]goserver.BankImporter, error)
	CreateBankImporter(userID string, bankImporter *goserver.BankImporterNoId) (goserver.BankImporter, error)
	UpdateBankImporter(
		userID string, id string, bankImporter goserver.BankImporterNoIdInterface,
	) (goserver.BankImporter, error)
	DeleteBankImporter(userID string, id string) error
	GetBankImporter(userID string, id string) (goserver.BankImporter, error)
	GetAllBankImporters() ([]ImportInfo, error)

	GetMatchers(userID string) ([]goserver.Matcher, error)
	GetMatcher(userID string, id string) (goserver.Matcher, error)
	// Add a single confirmation (true = confirmed, false = rejected) to a matcher
	// This operation is performed atomically and enforces the configured
	// confirmation history maximum length.
	AddMatcherConfirmation(userID string, id string, confirmed bool) error
	GetMatchersRuntime(userID string) ([]MatcherRuntime, error)
	GetMatcherRuntime(userID, id string) (MatcherRuntime, error)
	CreateMatcherRuntimeFromNoId(m goserver.MatcherNoIdInterface) (MatcherRuntime, error)
	CreateMatcher(userID string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	UpdateMatcher(userID string, id string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	DeleteMatcher(userID string, id string) error

	SaveCNBRates(rates map[string]float64, day time.Time) error
	GetCNBRates(day time.Time) (map[string]float64, error)

	CreateBudgetItem(userID string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	GetBudgetItems(userID string) ([]goserver.BudgetItem, error)
	GetBudgetItem(userID string, id string) (goserver.BudgetItem, error)
	UpdateBudgetItem(userID string, id string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	DeleteBudgetItem(userID string, id string) error

	CreateImage(data []byte, contentType string) (models.Image, error)
	GetImage(id string) (models.Image, error)
	DeleteImage(id string) error

	CreateNotification(userID string, notification *goserver.Notification) (goserver.Notification, error)
	GetNotifications(userID string) ([]goserver.Notification, error)
	DeleteNotification(userID string, id string) error
}

type MatcherRuntime struct {
	Matcher              *goserver.Matcher
	DescriptionRegexp    *regexp.Regexp
	PartnerAccountRegexp *regexp.Regexp
	PartnerNameRegexp    *regexp.Regexp
	CurrencyRegexp       *regexp.Regexp
	PlaceRegexp          *regexp.Regexp
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

func (s *storage) DeleteAccount(userID string, id string, replaceWithAccountID *string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if replaceWithAccountID != nil && *replaceWithAccountID != "" {
			newAccountID := *replaceWithAccountID

			// 1. Reassign BankImporters
			if err := tx.Model(&models.BankImporter{}).Where("account_id = ? AND user_id = ?", id, userID).
				Update("account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign bank importers: %w", err)
			}

			// 2. Reassign Matchers
			if err := tx.Model(&models.Matcher{}).Where("output_account_id = ? AND user_id = ?", id, userID).
				Update("output_account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign matchers: %w", err)
			}

			// 3. Reassign BudgetItems
			if err := tx.Model(&models.BudgetItem{}).Where("account_id = ? AND user_id = ?", id, userID).
				Update("account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign budget items: %w", err)
			}

			// 4. Reassign Transactions (Movements)
			// Since movements are stored as JSON, we need to fetch, modify, and save.
			// Ideally we should process in batches, but for simplicity/MVP we do it here.
			// We find all transactions that have a movement with this account ID.
			// Note: simpler approach is to iterate over all transactions for this user.
			// Optimization: Use a LIKE query on the JSON column if possible, but GORM/SQLite support varies.
			// Assuming low volume or acceptable performance for now.
			// A safer approach for SQLite/Postgres JSON: .Where("movements LIKE ?", "%"+id+"%")

			var transactions []models.Transaction
			// Use LIKE to pre-filter transactions that *might* have the account ID in movements
			// "movements" column contains the JSON.
			if err := tx.Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Find(&transactions).Error; err != nil {
				return fmt.Errorf("failed to find transactions for reassignment: %w", err)
			}

			for _, t := range transactions {
				updated := false
				newMovements := make([]goserver.Movement, len(t.Movements))
				for i, m := range t.Movements {
					if m.AccountId == id {
						m.AccountId = newAccountID
						updated = true
					}
					newMovements[i] = m
				}

				if updated {
					t.Movements = newMovements
					if err := tx.Save(&t).Error; err != nil {
						return fmt.Errorf("failed to save reassigned transaction %s: %w", t.ID, err)
					}
				}
			}
		} else {
			// User chose NOT to reassign.
			// Check if account is in use by any entity
			var count int64

			// Check BankImporters
			if err := tx.Model(&models.BankImporter{}).Where("account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check bank importers: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check Matchers
			if err := tx.Model(&models.Matcher{}).Where("output_account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check matchers: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check BudgetItems
			if err := tx.Model(&models.BudgetItem{}).Where("account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check budget items: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check Transactions (Movements)
			// Using the same LIKE query as before
			if err := tx.Model(&models.Transaction{}).Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions: %w", err)
			}
			if count > 0 {
				// We need to double check because LIKE is loose
				var transactions []models.Transaction
				if err := tx.Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Find(&transactions).Error; err != nil {
					return fmt.Errorf("failed to fetch transactions for verification: %w", err)
				}
				for _, t := range transactions {
					for _, m := range t.Movements {
						if m.AccountId == id {
							return ErrAccountInUse
						}
					}
				}
			}
		}

		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}
		return nil
	})
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

func (s *storage) DeleteCurrency(userID string, id string, replaceWithCurrencyID *string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if replaceWithCurrencyID != nil && *replaceWithCurrencyID != "" {
			newCurrencyID := *replaceWithCurrencyID

			// 1. Reassign in User (favorite currency)
			if err := tx.Model(&models.User{}).Where("id = ? AND favorite_currency_id = ?", userID, id).
				Update("favorite_currency_id", newCurrencyID).Error; err != nil {
				return fmt.Errorf("failed to reassign user favorite currency: %w", err)
			}

			// 2. Reassign in Accounts (BankInfo)
			var accounts []models.Account
			if err := tx.Where("user_id = ? AND bank_info LIKE ?", userID, "%"+id+"%").Find(&accounts).Error; err != nil {
				return fmt.Errorf("failed to find accounts for currency reassignment: %w", err)
			}
			for _, acc := range accounts {
				updated := false
				newBalances := make([]goserver.BankAccountInfoBalancesInner, len(acc.BankInfo.Balances))

				for i, b := range acc.BankInfo.Balances {
					if b.CurrencyId == id {
						b.CurrencyId = newCurrencyID
						updated = true
					}
					newBalances[i] = b
				}

				if updated {
					acc.BankInfo.Balances = newBalances
					if err := tx.Save(&acc).Error; err != nil {
						return fmt.Errorf("failed to save reassigned account %s: %w", acc.ID, err)
					}
				}
			}

			// 3. Reassign in Transactions (Movements)
			var transactions []models.Transaction
			if err := tx.Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Find(&transactions).Error; err != nil {
				return fmt.Errorf("failed to find transactions for currency reassignment: %w", err)
			}

			for _, t := range transactions {
				updated := false
				newMovements := make([]goserver.Movement, len(t.Movements))
				for i, m := range t.Movements {
					if m.CurrencyId == id {
						m.CurrencyId = newCurrencyID
						updated = true
					}
					newMovements[i] = m
				}

				if updated {
					t.Movements = newMovements
					if err := tx.Save(&t).Error; err != nil {
						return fmt.Errorf("failed to save reassigned transaction %s: %w", t.ID, err)
					}
				}
			}
		} else {
			// Check if currency is in use
			var count int64

			// Check User favorite currency
			if err := tx.Model(&models.User{}).Where("id = ? AND favorite_currency_id = ?", userID, id).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check user favorite currency: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
			}

			// Check Accounts
			if err := tx.Model(&models.Account{}).Where("user_id = ? AND bank_info LIKE ?", userID, "%"+id+"%").Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check accounts for currency usage: %w", err)
			}
			if count > 0 {
				var accounts []models.Account
				if err := tx.Where("user_id = ? AND bank_info LIKE ?", userID, "%"+id+"%").Find(&accounts).Error; err != nil {
					return fmt.Errorf("failed to fetch accounts for currency usage verification: %w", err)
				}
				for _, acc := range accounts {
					for _, b := range acc.BankInfo.Balances {
						if b.CurrencyId == id {
							return ErrCurrencyInUse
						}
					}
				}
			}

			// Check Transactions
			if err := tx.Model(&models.Transaction{}).Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions for currency usage: %w", err)
			}
			if count > 0 {
				var transactions []models.Transaction
				if err := tx.Where("user_id = ? AND movements LIKE ?", userID, "%"+id+"%").Find(&transactions).Error; err != nil {
					return fmt.Errorf("failed to fetch transactions for currency usage verification: %w", err)
				}
				for _, t := range transactions {
					for _, m := range t.Movements {
						if m.CurrencyId == id {
							return ErrCurrencyInUse
						}
					}
				}
			}
		}

		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Currency{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}
		return nil
	})
}

func (s *storage) GetTransactions(userID string, dateFrom, dateTo time.Time, onlySuspicious bool) ([]goserver.Transaction, error) {
	req := s.db.Model(&models.Transaction{}).Where("user_id = ? AND merged_into_id IS NULL", userID)
	if onlySuspicious {
		// Filter transactions where suspicious_reasons is not null and not an empty JSON array
		req = req.Where("suspicious_reasons IS NOT NULL AND suspicious_reasons != '[]' AND suspicious_reasons != ''")
	}
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

func (s *storage) GetTransactionsIncludingDeleted(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error) {
	req := s.db.Model(&models.Transaction{}).Unscoped().Where("user_id = ?", userID)
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

func (s *storage) CreateTransaction(userID string, input goserver.TransactionNoIdInterface,
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
	var t models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf(StorageError, err)
	}

	if len(t.ExternalIDs) > 0 || t.UnprocessedSources != "" {
		return ErrImportedTransactionCannotBeDeleted
	}

	if err := s.db.Delete(&t).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

func (s *storage) DeleteDuplicateTransaction(userID string, id, duplicateID string) error {
	s.log.Info("Deleting duplicate transaction", "id", id, "duplicate_id", duplicateID)
	return s.db.Transaction(func(tx *gorm.DB) error {
		var t, duplicate models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
			s.log.Warn("Failed to find transaction", "id", id, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		if err := tx.Where("id = ? AND user_id = ?", duplicateID, userID).First(&duplicate).Error; err != nil {
			s.log.Warn("Failed to find duplicate transaction", "id", duplicateID, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		duplicate.ExternalIDs = append(duplicate.ExternalIDs, t.ExternalIDs...)
		if err := tx.Save(&duplicate).Error; err != nil {
			s.log.Warn("Failed to update duplicate transaction", "id", duplicateID, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		// Mark as merged instead of deleting
		now := time.Now()
		if err := tx.Model(&models.Transaction{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(map[string]interface{}{
				"merged_into_id": duplicateID,
				"merged_at":      &now,
			}).Error; err != nil {
			s.log.Warn("Failed to mark transaction as merged", "id", id, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		return nil
	})
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

func (s *storage) GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error) {
	var mergedModels []models.Transaction
	if err := s.db.Where("user_id = ? AND merged_into_id IS NOT NULL", userID).Order("merged_at DESC").Find(&mergedModels).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.MergedTransaction, 0, len(mergedModels))
	for _, m := range mergedModels {
		var kept models.Transaction
		if err := s.db.Where("id = ? AND user_id = ?", m.MergedIntoID, userID).First(&kept).Error; err != nil {
			s.log.Warn("Failed to find kept transaction for merged transaction", "merged_id", m.ID, "kept_id", m.MergedIntoID)
			continue
		}

		result = append(result, goserver.MergedTransaction{
			Transaction: m.FromDB(),
			MergedInto:  kept.FromDB(),
			MergedAt:    *m.MergedAt,
		})
	}

	return result, nil
}

func (s *storage) UnmergeTransaction(userID, id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var t models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}

		if t.MergedIntoID == nil {
			return fmt.Errorf("transaction %s is not merged", id)
		}

		keptID := t.MergedIntoID.String()

		// 1. Remove external IDs from kept transaction
		var kept models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", keptID, userID).First(&kept).Error; err == nil {
			newExternalIDs := make([]string, 0)
			for _, extID := range kept.ExternalIDs {
				found := false
				for _, mergedExtID := range t.ExternalIDs {
					if extID == mergedExtID {
						found = true
						break
					}
				}
				if !found {
					newExternalIDs = append(newExternalIDs, extID)
				}
			}
			kept.ExternalIDs = newExternalIDs
			if err := tx.Save(&kept).Error; err != nil {
				return fmt.Errorf("failed to update kept transaction: %w", err)
			}
		}

		// 2. Clear merge fields
		if err := tx.Model(&models.Transaction{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(map[string]interface{}{
				"merged_into_id": gorm.Expr("NULL"),
				"merged_at":      gorm.Expr("NULL"),
			}).Error; err != nil {
			return fmt.Errorf("failed to unmerge transaction: %w", err)
		}

		return nil
	})
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
			FetchAll:         imp.FetchAll,
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

func (s *storage) createMatcherRuntime(m goserver.Matcher) (MatcherRuntime, error) {
	runtime := MatcherRuntime{Matcher: &m}
	if m.DescriptionRegExp != "" {
		r, err := regexp.Compile(m.DescriptionRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile description regexp: %w", err)
		}
		runtime.DescriptionRegexp = r
	}

	if m.PartnerAccountNumberRegExp != "" {
		r, err := regexp.Compile(m.PartnerAccountNumberRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile partner account regexp: %w", err)
		}
		runtime.PartnerAccountRegexp = r
	}

	if m.PartnerNameRegExp != "" {
		r, err := regexp.Compile(m.PartnerNameRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile partner name regexp: %w", err)
		}
		runtime.PartnerNameRegexp = r
	}

	if m.CurrencyRegExp != "" {
		r, err := regexp.Compile(m.CurrencyRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile currency regexp: %w", err)
		}
		runtime.CurrencyRegexp = r
	}

	if m.PlaceRegExp != "" {
		r, err := regexp.Compile(m.PlaceRegExp)
		if err != nil {
			return MatcherRuntime{}, fmt.Errorf("failed to compile place regexp: %w", err)
		}
		runtime.PlaceRegexp = r
	}

	return runtime, nil
}

// CreateMatcherRuntimeFromNoId creates a MatcherRuntime from a MatcherNoId (without needing to save to DB first).
// This is useful for testing matchers before they are saved.
//
//nolint:stylecheck
func (s *storage) CreateMatcherRuntimeFromNoId(m goserver.MatcherNoIdInterface) (MatcherRuntime, error) {
	// Convert MatcherNoId to Matcher by creating a temporary matcher with empty ID
	matcher := goserver.Matcher{
		OutputDescription:          m.GetOutputDescription(),
		OutputAccountId:            m.GetOutputAccountId(),
		OutputTags:                 m.GetOutputTags(),
		CurrencyRegExp:             m.GetCurrencyRegExp(),
		PartnerNameRegExp:          m.GetPartnerNameRegExp(),
		PartnerAccountNumberRegExp: m.GetPartnerAccountNumberRegExp(),
		DescriptionRegExp:          m.GetDescriptionRegExp(),
		ExtraRegExp:                m.GetExtraRegExp(),
		PlaceRegExp:                m.GetPlaceRegExp(),
		ConfirmationHistory:        m.GetConfirmationHistory(),
	}

	return s.createMatcherRuntime(matcher)
}

func (s *storage) GetMatcherRuntime(userID, id string) (MatcherRuntime, error) {
	m, err := s.GetMatcher(userID, id)
	if err != nil {
		return MatcherRuntime{}, err
	}

	return s.createMatcherRuntime(m)
}

func (s *storage) GetMatchersRuntime(userID string) ([]MatcherRuntime, error) {
	matchers, err := s.GetMatchers(userID)
	if err != nil {
		return nil, err
	}

	res := make([]MatcherRuntime, 0, len(matchers))
	for _, m := range matchers {
		runtime, err := s.createMatcherRuntime(m)
		if err != nil {
			return nil, err
		}

		res = append(res, runtime)
	}

	return res, nil
}

//nolint:dupl
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

// AddMatcherConfirmation atomically appends a confirmation boolean to the matcher's
// confirmation history and trims it to the configured maximum length.
func (s *storage) AddMatcherConfirmation(userID string, id string, confirmed bool) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var m models.Matcher
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&m).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.log.Warn("Matcher not found when adding confirmation", "userID", userID, "matcherID", id)
				return ErrNotFound
			}
			s.log.Error("DB error when loading matcher for confirmation", "error", err)
			return fmt.Errorf(StorageError, err)
		}

		// Use the model helper to add confirmation and respect config max length
		m.AddConfirmation(confirmed, s.cfg.MatcherConfirmationHistoryMax)

		if err := tx.Save(&m).Error; err != nil {
			s.log.Error("DB error when saving matcher after adding confirmation", "error", err)
			return fmt.Errorf(StorageError, err)
		}

		return nil
	})
}

//#endregion Matchers

// #region CNB rates
func (s *storage) SaveCNBRates(rates map[string]float64, date time.Time) error {
	// Use a transaction to ensure all rates are saved together
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First delete all existing rates for this date to avoid duplicates
		if err := tx.Where("rate_date = ?", date).Delete(&models.CNBCurrencyRate{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}

		// Create new rates
		for currencyCode, rate := range rates {
			currencyRate := models.CNBCurrencyRate{
				CurrencyCode: currencyCode,
				RateToCZK:    rate,
				RateDate:     date,
			}

			if err := tx.Create(&currencyRate).Error; err != nil {
				return fmt.Errorf(StorageError, err)
			}
		}

		return nil
	})
}

func (s *storage) GetCNBRates(date time.Time) (map[string]float64, error) {
	var rates []models.CNBCurrencyRate
	query := s.db.Model(&models.CNBCurrencyRate{})

	// If a specific date is provided, use it
	if !date.IsZero() {
		query = query.Where("rate_date = ?", date)
	} else {
		// Otherwise get the most recent rates
		var latestDate time.Time
		if err := s.db.Model(&models.CNBCurrencyRate{}).
			Select("MAX(rate_date)").
			Scan(&latestDate).Error; err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		query = query.Where("rate_date = ?", latestDate)
	}

	if err := query.Find(&rates).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	// Convert to map
	result := make(map[string]float64, len(rates))
	for _, rate := range rates {
		result[rate.CurrencyCode] = rate.RateToCZK
	}

	return result, nil
}

// #endregion CNB rates

// #region Images

func (s *storage) CreateImage(data []byte, contentType string) (models.Image, error) {
	image := models.Image{
		Data:        data,
		ContentType: contentType,
	}

	if err := s.db.Create(&image).Error; err != nil {
		return models.Image{}, fmt.Errorf(StorageError, err)
	}

	return image, nil
}

func (s *storage) GetImage(id string) (models.Image, error) {
	var image models.Image
	if err := s.db.Where("id = ?", id).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Image{}, ErrNotFound
		}
		return models.Image{}, fmt.Errorf(StorageError, err)
	}

	return image, nil
}

func (s *storage) DeleteImage(id string) error {
	if err := s.db.Where("id = ?", id).Delete(&models.Image{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}
	return nil
}

// #endregion Images
//#endregion CNB rates

// #region BudgetItems
func (s *storage) CreateBudgetItem(userID string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error) {
	data := models.BudgetItemToDB(budgetItem, userID)
	data.ID = uuid.New()
	if err := s.db.Create(data).Error; err != nil {
		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}
	s.log.Info("BudgetItem created", "id", data.ID)

	return data.FromDB(), nil
}

func (s *storage) GetBudgetItems(userID string) ([]goserver.BudgetItem, error) {
	result, err := s.db.Model(&models.BudgetItem{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	items := make([]goserver.BudgetItem, 0)
	for result.Next() {
		var item models.BudgetItem
		if err := s.db.ScanRows(result, &item); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		items = append(items, item.FromDB())
	}

	return items, nil
}

func (s *storage) GetBudgetItem(userID string, id string) (goserver.BudgetItem, error) {
	var data models.BudgetItem
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BudgetItem{}, ErrNotFound
		}

		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) UpdateBudgetItem(
	userID string, id string, budgetItem *goserver.BudgetItemNoId,
) (goserver.BudgetItem, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return goserver.BudgetItem{}, fmt.Errorf(StorageError+"; id is not UUID", err)
	}

	var data *models.BudgetItem
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.BudgetItem{}, ErrNotFound
		}

		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}

	data = models.BudgetItemToDB(budgetItem, userID)
	data.ID = idUUID
	if err := s.db.Save(&data).Error; err != nil {
		return goserver.BudgetItem{}, fmt.Errorf(StorageError, err)
	}

	return data.FromDB(), nil
}

func (s *storage) DeleteBudgetItem(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.BudgetItem{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

//#endregion BudgetItems

// #region Notifications
func (s *storage) CreateNotification(userID string, notification *goserver.Notification) (goserver.Notification, error) {
	n, err := models.NotificationToDB(notification, userID)
	if err != nil {
		return goserver.Notification{}, fmt.Errorf(StorageError, err)
	}

	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}

	if err := s.db.Create(n).Error; err != nil {
		return goserver.Notification{}, fmt.Errorf(StorageError, err)
	}

	return n.FromDB(), nil
}

func (s *storage) GetNotifications(userID string) ([]goserver.Notification, error) {
	result, err := s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Order("date DESC").Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	notifications := make([]goserver.Notification, 0)
	for result.Next() {
		var n models.Notification
		if err := s.db.ScanRows(result, &n); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		notifications = append(notifications, n.FromDB())
	}

	return notifications, nil
}

func (s *storage) DeleteNotification(userID string, id string) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	return nil
}

// #endregion Notifications

package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
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
	GetAllUserIDs() ([]string, error)

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
	MergeTransactions(userID, keepID, mergeID string) (goserver.Transaction, error)
	DeleteDuplicateTransaction(userID string, id, duplicateID string) error
	GetTransaction(userID string, id string) (goserver.Transaction, error)
	GetTransactionsIncludingDeleted(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error)
	GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error)
	GetMergedTransaction(userID, originalTransactionID string) (goserver.MergedTransaction, error)
	UnmergeTransaction(userID, id string) error
	GetDuplicateTransactionIDs(userID, transactionID string) ([]string, error)
	AddDuplicateRelationship(userID, transactionID1, transactionID2 string) error
	RemoveDuplicateRelationship(userID, transactionID1, transactionID2 string) error
	ClearDuplicateRelationships(userID, transactionID string) error

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

	SaveCNBRates(rates map[string]decimal.Decimal, day time.Time) error
	GetCNBRates(day time.Time) (map[string]decimal.Decimal, error)

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

	GetAccountBalance(userID, accountID, currencyID string) (decimal.Decimal, error)
	CountUnprocessedTransactionsForAccount(userID, accountID string, ignoreUnprocessedBefore time.Time) (int, error)
	HasTransactionsAfterDate(userID, accountID string, date time.Time) (bool, error)

	// Reconciliation methods
	GetLatestReconciliation(userID, accountID, currencyID string) (*goserver.Reconciliation, error)
	GetReconciliationsForAccount(userID, accountID string) ([]goserver.Reconciliation, error)
	CreateReconciliation(userID string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error)
	InvalidateReconciliation(userID, accountID, currencyID string) error
	GetBulkReconciliationData(userID string) (*BulkReconciliationData, error)
	Backup(destination string) error
}

type MatcherRuntime struct {
	Matcher              *goserver.Matcher
	DescriptionRegexp    *regexp.Regexp
	PartnerAccountRegexp *regexp.Regexp
	PartnerNameRegexp    *regexp.Regexp
	CurrencyRegexp       *regexp.Regexp
	PlaceRegexp          *regexp.Regexp
	Keywords             []string
	KeywordOutputs       []string
	KeywordRegexps       []*regexp.Regexp
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

func (s *storage) GetAllUserIDs() ([]string, error) {
	var ids []string
	if err := s.db.Model(&models.User{}).Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	return ids, nil
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
			// Optimization: Use SQLite's JSON functions for accurate querying.
			var transactions []models.Transaction
			if err := tx.Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.accountId') = ?", userID, id).
				Group("transactions.id").
				Find(&transactions).Error; err != nil {
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
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("transactions").
				Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.accountId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
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
			if err := tx.Joins("CROSS JOIN json_each(accounts.bank_info, '$.balances')").
				Where("accounts.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Group("accounts.id").
				Find(&accounts).Error; err != nil {
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
			if err := tx.Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Group("transactions.id").
				Find(&transactions).Error; err != nil {
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
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("accounts").
				Joins("CROSS JOIN json_each(accounts.bank_info, '$.balances')").
				Where("accounts.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check accounts for currency usage: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
			}

			// Check Transactions
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("transactions").
				Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions for currency usage: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
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

	// Populate DuplicateTransactionIds in batch
	if len(transactions) > 0 {
		ids := make([]uuid.UUID, len(transactions))
		for i, t := range transactions {
			id, _ := uuid.Parse(t.Id)
			ids[i] = id
		}

		var relationships []models.TransactionDuplicate
		if err := s.db.Where("user_id = ? AND transaction_id1 IN ?", userID, ids).Find(&relationships).Error; err == nil {
			relMap := make(map[string][]string)
			for _, r := range relationships {
				t1 := r.TransactionID1.String()
				t2 := r.TransactionID2.String()
				relMap[t1] = append(relMap[t1], t2)
			}

			for i := range transactions {
				if dups, ok := relMap[transactions[i].Id]; ok {
					transactions[i].DuplicateTransactionIds = dups
				}
			}
		}

		// Populate MergedTransactionIds in batch (find transactions that were merged into these from archive)
		var mergedRecords []struct {
			KeptTransactionID     uuid.UUID
			OriginalTransactionID uuid.UUID
		}
		if err := s.db.Model(&models.MergedTransaction{}).
			Select("kept_transaction_id, original_transaction_id").
			Where("user_id = ? AND kept_transaction_id IN ?", userID, ids).
			Find(&mergedRecords).Error; err == nil {

			mergedMap := make(map[string][]string)
			for _, r := range mergedRecords {
				key := r.KeptTransactionID.String()
				mergedMap[key] = append(mergedMap[key], r.OriginalTransactionID.String())
			}
			for i := range transactions {
				if merged, ok := mergedMap[transactions[i].Id]; ok {
					transactions[i].MergedTransactionIds = merged
				}
			}
		}
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

	if err := s.recordTransactionHistory(s.db, userID, t, "CREATED"); err != nil {
		s.log.Error("Failed to record transaction history", "error", err)
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

	if err := s.recordTransactionHistory(s.db, userID, t, "UPDATED"); err != nil {
		s.log.Error("Failed to record transaction history", "error", err)
	}

	// Get old movements for smart invalidation
	oldMovements := models.MovementsToAPI(t.Movements)

	t = models.TransactionToDB(input, userID)
	t.ID = idUUID
	if err := s.db.Save(&t).Error; err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	// If dismissed, clear relationships
	if t.DuplicateDismissed {
		if err := s.ClearDuplicateRelationships(userID, id); err != nil {
			s.log.Error("Failed to clear duplicate relationships on dismissal", "error", err, "id", id)
		}
	} else {
		// Revalidate duplicate links in case date/amount changed
		if err := s.RevalidateDuplicateRelationships(userID, id); err != nil {
			s.log.Error("Failed to revalidate duplicate relationships", "error", err, "id", id)
		}
	}

	// Smart invalidation: only if amounts or currencies changed
	s.invalidateReconciliationIfAmountsChanged(userID, oldMovements, models.MovementsToAPI(t.Movements), t.Date)

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

	if err := s.recordTransactionHistory(s.db, userID, &t, "DELETED"); err != nil {
		s.log.Error("Failed to record transaction history", "error", err)
	}

	if err := s.db.Delete(&t).Error; err != nil {
		return fmt.Errorf(StorageError, err)
	}

	// Clear duplicate relationships if any
	if err := s.ClearDuplicateRelationships(userID, id); err != nil {
		s.log.Error("Failed to clear duplicate relationships on deletion", "error", err, "id", id)
	}

	// Invalidate reconciliation for deleted movements
	s.invalidateReconciliationIfAmountsChanged(userID, models.MovementsToAPI(t.Movements), []goserver.Movement{}, t.Date)

	return nil
}

func (s *storage) MergeTransactions(userID, keepID, mergeID string) (goserver.Transaction, error) {
	kID, err := uuid.Parse(keepID)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf("invalid keep ID: %w", err)
	}
	mID, err := uuid.Parse(mergeID)
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf("invalid merge ID: %w", err)
	}

	var keepT models.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var mergeT models.Transaction
		if err := tx.Where("user_id = ? AND id = ?", userID, kID).First(&keepT).Error; err != nil {
			return fmt.Errorf("failed to find keep transaction: %w", err)
		}
		if err := tx.Where("user_id = ? AND id = ?", userID, mID).First(&mergeT).Error; err != nil {
			return fmt.Errorf("failed to find merge transaction: %w", err)
		}

		// 2. Transfer external IDs
		existingIDs := make(map[string]bool)
		for _, id := range keepT.ExternalIDs {
			existingIDs[id] = true
		}
		for _, id := range mergeT.ExternalIDs {
			if !existingIDs[id] {
				keepT.ExternalIDs = append(keepT.ExternalIDs, id)
			}
		}

		// 3. Update keepT: remove suspicious reason, inherit external IDs
		newReasons := make([]string, 0)
		for _, r := range keepT.SuspiciousReasons {
			if r != "Potential duplicate from different importer" {
				newReasons = append(newReasons, r)
			}
		}
		keepT.SuspiciousReasons = newReasons
		if err := tx.Save(&keepT).Error; err != nil {
			return fmt.Errorf("failed to update keep transaction: %w", err)
		}

		// 4. Archive and hard-delete mergeT (archive is now source of truth)
		now := time.Now()

		if err := s.archiveMergedTransaction(tx, userID, &mergeT, kID, now); err != nil {
			return err
		}

		// Hard-delete the merged transaction
		if err := tx.Unscoped().Delete(&mergeT).Error; err != nil {
			return fmt.Errorf("failed to hard-delete merge transaction: %w", err)
		}

		// 5. Clear duplicate relationships for both (they are resolved now)
		// We use tx so it's atomic within our transaction
		if err := s.clearDuplicateRelationshipsWithTx(tx, userID, mID.String()); err != nil {
			return fmt.Errorf("failed to clear relationships for merged transaction: %w", err)
		}

		// Clear specific relationship between keepT and mergeT (the one where keepT was the primary)
		if err := s.clearDuplicateRelationshipsWithTx(tx, userID, kID.String()); err != nil {
			return fmt.Errorf("failed to clear relationships for keep transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return goserver.Transaction{}, fmt.Errorf(StorageError, err)
	}

	return s.GetTransaction(userID, keepID)
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

		// Archive and hard-delete the transaction
		now := time.Now()
		duplicateIDUUID, _ := uuid.Parse(duplicateID)
		if err := s.archiveMergedTransaction(tx, userID, &t, duplicateIDUUID, now); err != nil {
			return err
		}

		// Hard-delete the transaction (archive is source of truth)
		if err := tx.Unscoped().Delete(&t).Error; err != nil {
			s.log.Warn("Failed to hard-delete merged transaction", "id", id, "error", err)
			return fmt.Errorf(StorageError, err)
		}

		if err := s.recordTransactionHistory(tx, userID, &t, "MERGED"); err != nil {
			s.log.Error("Failed to record transaction history", "error", err)
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

	apiTransaction := transaction.FromDB()
	duplicateIds, err := s.GetDuplicateTransactionIDs(userID, id)
	if err == nil {
		apiTransaction.DuplicateTransactionIds = duplicateIds
	}

	// Populate MergedTransactionIds for single transaction from archive
	var mergedIds []string
	if err := s.db.Model(&models.MergedTransaction{}).
		Where("user_id = ? AND kept_transaction_id = ?", userID, id).
		Pluck("original_transaction_id", &mergedIds).Error; err == nil {
		apiTransaction.MergedTransactionIds = mergedIds
	}

	return apiTransaction, nil
}

func (s *storage) GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error) {
	var mergedModels []models.MergedTransaction
	if err := s.db.Where("user_id = ?", userID).Order("merged_at DESC").Find(&mergedModels).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	result := make([]goserver.MergedTransaction, 0, len(mergedModels))
	for _, m := range mergedModels {
		var kept models.Transaction
		if err := s.db.Where("id = ? AND user_id = ?", m.KeptTransactionID, userID).First(&kept).Error; err != nil {
			s.log.Warn("Failed to find kept transaction for merged transaction", "merged_id", m.OriginalTransactionID, "kept_id", m.KeptTransactionID)
			continue
		}

		// Create a temporary transaction structure to use FromDB
		tr := models.Transaction{
			ID:                 m.OriginalTransactionID,
			UserID:             m.UserID,
			Date:               m.Date,
			Description:        m.Description,
			Place:              m.Place,
			Tags:               m.Tags,
			PartnerName:        m.PartnerName,
			PartnerAccount:     m.PartnerAccount,
			PartnerInternalID:  m.PartnerInternalID,
			Extra:              m.Extra,
			UnprocessedSources: m.UnprocessedSources,
			ExternalIDs:        m.ExternalIDs,
			Movements:          m.Movements,
			MatcherID:          m.MatcherID,
			IsAuto:             m.IsAuto,
			SuspiciousReasons:  m.SuspiciousReasons,
		}

		result = append(result, goserver.MergedTransaction{
			Transaction: tr.FromDB(),
			MergedInto:  kept.FromDB(),
			MergedAt:    m.MergedAt,
		})
	}

	return result, nil
}

func (s *storage) GetMergedTransaction(userID, originalTransactionID string) (goserver.MergedTransaction, error) {
	// 1. Find the archived transaction
	var archived models.MergedTransaction
	if err := s.db.Where("user_id = ? AND original_transaction_id = ?", userID, originalTransactionID).First(&archived).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.MergedTransaction{}, ErrNotFound
		}
		return goserver.MergedTransaction{}, fmt.Errorf(StorageError, err)
	}

	// 2. Find the kept transaction (mergedInto)
	var kept models.Transaction
	if err := s.db.Where("id = ? AND user_id = ?", archived.KeptTransactionID, userID).First(&kept).Error; err != nil {
		s.log.Warn("Failed to find kept transaction for merged transaction", "merged_id", archived.OriginalTransactionID, "kept_id", archived.KeptTransactionID)
		// We still return the merged transaction, just without the full kept transaction details if not found
	}

	// 3. Construct the response
	// Create a temporary transaction structure to use FromDB for the archived transaction
	tr := models.Transaction{
		ID:                 archived.OriginalTransactionID,
		UserID:             archived.UserID,
		Date:               archived.Date,
		Description:        archived.Description,
		Place:              archived.Place,
		Tags:               archived.Tags,
		PartnerName:        archived.PartnerName,
		PartnerAccount:     archived.PartnerAccount,
		PartnerInternalID:  archived.PartnerInternalID,
		Extra:              archived.Extra,
		UnprocessedSources: archived.UnprocessedSources,
		ExternalIDs:        archived.ExternalIDs,
		Movements:          archived.Movements,
		MatcherID:          archived.MatcherID,
		IsAuto:             archived.IsAuto,
		SuspiciousReasons:  archived.SuspiciousReasons,
	}

	return goserver.MergedTransaction{
		Transaction: tr.FromDB(),
		MergedInto:  kept.FromDB(),
		MergedAt:    archived.MergedAt,
	}, nil
}

func (s *storage) UnmergeTransaction(userID, id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Find the archived transaction
		var archived models.MergedTransaction
		if err := tx.Where("user_id = ? AND original_transaction_id = ?", userID, id).First(&archived).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("transaction %s is not merged or archive not found", id)
			}
			return fmt.Errorf(StorageError, err)
		}

		keptID := archived.KeptTransactionID.String()

		// 2. Remove external IDs from kept transaction
		var kept models.Transaction
		if err := tx.Where("id = ? AND user_id = ?", keptID, userID).First(&kept).Error; err == nil {
			newExternalIDs := make([]string, 0)
			for _, extID := range kept.ExternalIDs {
				found := false
				for _, mergedExtID := range archived.ExternalIDs {
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

		// 3. Recreate the transaction from archive data (don't rely on soft-deleted record)
		restoredTransaction := models.Transaction{
			ID:                 archived.OriginalTransactionID,
			UserID:             archived.UserID,
			Date:               archived.Date,
			Description:        archived.Description,
			Place:              archived.Place,
			Tags:               archived.Tags,
			PartnerName:        archived.PartnerName,
			PartnerAccount:     archived.PartnerAccount,
			PartnerInternalID:  archived.PartnerInternalID,
			Extra:              archived.Extra,
			UnprocessedSources: archived.UnprocessedSources,
			ExternalIDs:        archived.ExternalIDs,
			Movements:          archived.Movements,
			MatcherID:          archived.MatcherID,
			IsAuto:             archived.IsAuto,
			SuspiciousReasons:  archived.SuspiciousReasons,
			// MergedIntoID and MergedAt are left nil (transaction is no longer merged)
		}

		// Create the restored transaction
		if err := tx.Create(&restoredTransaction).Error; err != nil {
			return fmt.Errorf("failed to recreate transaction: %w", err)
		}

		// 4. Delete from archive
		if err := tx.Delete(&archived).Error; err != nil {
			return fmt.Errorf("failed to delete from archive during unmerge: %w", err)
		}

		// 5. Record history
		if err := s.recordTransactionHistory(tx, userID, &restoredTransaction, "UNMERGED"); err != nil {
			s.log.Error("Failed to record transaction history", "error", err)
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

	if m.Simplified {
		runtime.Keywords = make([]string, len(m.Keywords))
		runtime.KeywordOutputs = make([]string, len(m.Keywords))
		runtime.KeywordRegexps = make([]*regexp.Regexp, len(m.Keywords))
		for i, k := range m.Keywords {
			matcherPart := k
			outputPart := k

			if idx := strings.Index(k, "|"); idx != -1 {
				matcherPart = k[:idx]
				outputPart = k[idx+1:]
			}

			runtime.Keywords[i] = matcherPart
			runtime.KeywordOutputs[i] = outputPart

			// Case-insensitive, whole-word matching
			// We wrap the keyword in \b (word boundary)
			r, err := regexp.Compile(`(?i)\b` + regexp.QuoteMeta(matcherPart) + `\b`)
			if err != nil {
				return MatcherRuntime{}, fmt.Errorf("failed to compile keyword regexp %q: %w", matcherPart, err)
			}
			runtime.KeywordRegexps[i] = r
		}
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
		Simplified:                 m.GetSimplified(),
		Keywords:                   m.GetKeywords(),
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
func (s *storage) SaveCNBRates(rates map[string]decimal.Decimal, date time.Time) error {
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

func (s *storage) GetCNBRates(date time.Time) (map[string]decimal.Decimal, error) {
	var rates []models.CNBCurrencyRate
	query := s.db.Model(&models.CNBCurrencyRate{})

	// If a specific date is provided, use it
	if !date.IsZero() {
		query = query.Where("rate_date = ?", date)
	} else {
		// Otherwise get the most recent rates
		var latestRate models.CNBCurrencyRate
		if err := s.db.Model(&models.CNBCurrencyRate{}).
			Order("rate_date DESC").
			First(&latestRate).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return make(map[string]decimal.Decimal), nil
			}
			return nil, fmt.Errorf(StorageError, err)
		}

		query = query.Where("rate_date = ?", latestRate.RateDate)
	}

	if err := query.Find(&rates).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	// Convert to map
	result := make(map[string]decimal.Decimal, len(rates))
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

func (s *storage) GetAccountBalance(userID, accountID, currencyID string) (decimal.Decimal, error) {
	acc, err := s.GetAccount(userID, accountID)
	if err != nil {
		return decimal.Zero, err
	}

	var total decimal.Decimal
	for _, b := range acc.BankInfo.Balances {
		if b.CurrencyId == currencyID {
			total = total.Add(b.OpeningBalance)
			break
		}
	}

	// Sum all movements for this account and currency
	// We use raw SQL to iterate over the movements JSON column in SQLite
	// Since movements is a JSON array of objects, we need to parse it.
	// For simplicity and to avoid complex SQLite JSON path expressions that might vary,
	// we fetch transactions and sum in Go, but filter at DB level if possible.
	var transactions []models.Transaction
	err = s.db.Where("user_id = ? AND movements LIKE ? AND merged_into_id IS NULL", userID, "%"+accountID+"%").Find(&transactions).Error
	if err != nil {
		return decimal.Zero, fmt.Errorf(StorageError, err)
	}

	for _, t := range transactions {
		for _, m := range t.Movements {
			if m.AccountId == accountID && m.CurrencyId == currencyID {
				total = total.Add(m.Amount)
			}
		}
	}

	return total, nil
}

func (s *storage) CountUnprocessedTransactionsForAccount(userID, accountID string, ignoreUnprocessedBefore time.Time) (int, error) {
	var count int
	// An unprocessed transaction is one that has at least one movement with an empty AccountId.
	// We also filter by accountID being present in at least one movement.
	var transactions []models.Transaction
	query := s.db.Where("user_id = ? AND movements LIKE ? AND merged_into_id IS NULL", userID, "%"+accountID+"%")
	if !ignoreUnprocessedBefore.IsZero() {
		query = query.Where("date >= ?", ignoreUnprocessedBefore)
	}
	err := query.Find(&transactions).Error
	if err != nil {
		return 0, fmt.Errorf(StorageError, err)
	}

	s.log.Debug("CountUnprocessedTransactionsForAccount query result",
		"accountId", accountID, "ignoreUnprocessedBefore", ignoreUnprocessedBefore, "totalTransactions", len(transactions))

	for _, t := range transactions {
		hasEmpty := false
		hasAccount := false
		for _, m := range t.Movements {
			// If a movement has 0 amount, it doesn't represent a financial impact
			// and shouldn't block reconciliation even if its AccountId is empty.
			if m.Amount.IsZero() {
				continue
			}
			if m.AccountId == "" {
				hasEmpty = true
			}
			if m.AccountId == accountID {
				hasAccount = true
			}
		}
		if hasEmpty && hasAccount {
			count++
		}
	}

	return count, nil
}

func (s *storage) HasTransactionsAfterDate(userID, accountID string, date time.Time) (bool, error) {
	var count int64
	err := s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND movements LIKE ? AND date > ? AND merged_into_id IS NULL", userID, "%"+accountID+"%", date).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf(StorageError, err)
	}
	return count > 0, nil
}

func (s *storage) recordTransactionHistory(tx *gorm.DB, userID string, transaction *models.Transaction, action string) error {
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction for history: %w", err)
	}

	history := models.TransactionHistory{
		ID:            uuid.New(),
		TransactionID: transaction.ID,
		UserID:        userID,
		Action:        action,
		Snapshot:      string(jsonData),
		CreatedAt:     time.Now(),
	}

	return tx.Create(&history).Error
}

func (s *storage) GetLatestReconciliation(userID, accountID, currencyID string) (*goserver.Reconciliation, error) {
	var rec models.Reconciliation
	result := s.db.Where("user_id = ? AND account_id = ? AND currency_id = ?", userID, accountID, currencyID).
		Order("reconciled_at DESC").First(&rec)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No reconciliation yet
		}
		return nil, fmt.Errorf("failed to get latest reconciliation: %w", result.Error)
	}
	return models.ReconciliationToAPI(&rec), nil
}

func (s *storage) GetBulkReconciliationData(userID string) (*BulkReconciliationData, error) {
	data := &BulkReconciliationData{
		Balances:              make(map[string]map[string]decimal.Decimal),
		LatestReconciliations: make(map[string]map[string]*goserver.Reconciliation),
		UnprocessedCounts:     make(map[string]int),
		MaxTransactionDates:   make(map[string]map[string]time.Time),
	}

	// 1. Get Accounts for opening balances and ignore dates
	accounts, err := s.GetAccounts(userID)
	if err != nil {
		return nil, err
	}
	ignoreMap := make(map[string]time.Time)
	for _, acc := range accounts {
		data.Balances[acc.Id] = make(map[string]decimal.Decimal)
		data.MaxTransactionDates[acc.Id] = make(map[string]time.Time)
		for _, b := range acc.BankInfo.Balances {
			data.Balances[acc.Id][b.CurrencyId] = b.OpeningBalance
		}
		if !acc.IgnoreUnprocessedBefore.IsZero() {
			ignoreMap[acc.Id] = acc.IgnoreUnprocessedBefore
		}
	}

	// 2. Get latest reconciliations
	var recs []models.Reconciliation
	err = s.db.Where("user_id = ?", userID).Order("reconciled_at DESC").Find(&recs).Error
	if err != nil {
		return nil, err
	}
	for i := range recs {
		r := &recs[i]
		apiRec := models.ReconciliationToAPI(r)
		if _, ok := data.LatestReconciliations[apiRec.AccountId]; !ok {
			data.LatestReconciliations[apiRec.AccountId] = make(map[string]*goserver.Reconciliation)
		}
		if _, ok := data.LatestReconciliations[apiRec.AccountId][apiRec.CurrencyId]; !ok {
			data.LatestReconciliations[apiRec.AccountId][apiRec.CurrencyId] = apiRec
		}
	}

	// 3. Get all transactions (merged_into_id IS NULL)
	var transactions []models.Transaction
	err = s.db.Where("user_id = ? AND merged_into_id IS NULL", userID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	for _, t := range transactions {
		hasEmpty := false
		involvedAccounts := make(map[string]bool)

		for _, m := range t.Movements {
			// Update balances
			if m.AccountId != "" {
				if _, ok := data.Balances[m.AccountId]; ok {
					data.Balances[m.AccountId][m.CurrencyId] = data.Balances[m.AccountId][m.CurrencyId].Add(m.Amount)
				}
				involvedAccounts[m.AccountId] = true

				// Update MaxTransactionDates
				if _, ok := data.MaxTransactionDates[m.AccountId]; !ok {
					data.MaxTransactionDates[m.AccountId] = make(map[string]time.Time)
				}
				if t.Date.After(data.MaxTransactionDates[m.AccountId][m.CurrencyId]) {
					data.MaxTransactionDates[m.AccountId][m.CurrencyId] = t.Date
				}
			}

			if m.Amount.IsZero() {
				continue
			}
			if m.AccountId == "" {
				hasEmpty = true
			}
		}

		// Update unprocessed counts
		if hasEmpty {
			for accID := range involvedAccounts {
				ignoreDate := ignoreMap[accID]
				if ignoreDate.IsZero() || !t.Date.Before(ignoreDate) {
					data.UnprocessedCounts[accID]++
				}
			}
		}
	}

	return data, nil
}

func (s *storage) GetReconciliationsForAccount(userID, accountID string) ([]goserver.Reconciliation, error) {
	var recs []models.Reconciliation
	err := s.db.Where("user_id = ? AND account_id = ?", userID, accountID).
		Order("reconciled_at DESC").Find(&recs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reconciliations for account: %w", err)
	}

	result := make([]goserver.Reconciliation, len(recs))
	for i := range recs {
		result[i] = *models.ReconciliationToAPI(&recs[i])
	}
	return result, nil
}

func (s *storage) CreateReconciliation(userID string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
	model := models.ReconciliationFromAPI(userID, rec)
	model.ID = uuid.New()
	model.ReconciledAt = time.Now()

	if err := s.db.Create(&model).Error; err != nil {
		return goserver.Reconciliation{}, fmt.Errorf("failed to create reconciliation: %w", err)
	}
	return *models.ReconciliationToAPI(model), nil
}

func (s *storage) InvalidateReconciliation(userID, accountID, currencyID string) error {
	err := s.db.Where("user_id = ? AND account_id = ? AND currency_id = ?",
		userID, accountID, currencyID).Delete(&models.Reconciliation{}).Error
	if err != nil {
		return fmt.Errorf("failed to invalidate reconciliation: %w", err)
	}
	return nil
}

func (s *storage) invalidateReconciliationIfAmountsChanged(
	userID string,
	oldMovements, newMovements []goserver.Movement,
	txDate time.Time,
) {
	// Build lookup for old movements
	oldByKey := make(map[string]goserver.Movement)
	for _, m := range oldMovements {
		key := m.AccountId + "|" + m.CurrencyId
		oldByKey[key] = m
	}

	// Check if any financial data changed
	affectedAccounts := make(map[string]string) // accountId -> currencyId

	for _, newM := range newMovements {
		if newM.AccountId == "" {
			continue // Unprocessed movements don't affect reconciliation directly
		}
		key := newM.AccountId + "|" + newM.CurrencyId
		oldM, exists := oldByKey[key]

		// New movement or amount changed
		if !exists || oldM.Amount != newM.Amount {
			affectedAccounts[newM.AccountId] = newM.CurrencyId
		}
		delete(oldByKey, key)
	}

	// Remaining old movements were removed
	for _, oldM := range oldByKey {
		if oldM.AccountId != "" {
			affectedAccounts[oldM.AccountId] = oldM.CurrencyId
		}
	}

	// Only invalidate if there were actual financial changes
	for accountId, currencyId := range affectedAccounts {
		lastRec, err := s.GetLatestReconciliation(userID, accountId, currencyId)
		if err != nil || lastRec == nil {
			continue
		}
		if txDate.Before(lastRec.ReconciledAt) {
			s.log.Info("Invalidating reconciliation due to financial change",
				"accountId", accountId, "currencyId", currencyId, "txDate", txDate, "recAt", lastRec.ReconciledAt)

			if err := s.InvalidateReconciliation(userID, accountId, currencyId); err != nil {
				s.log.Error("Failed to invalidate reconciliation", "error", err)
				continue
			}

			accountName := accountId
			if acc, err := s.GetAccount(userID, accountId); err == nil {
				accountName = acc.Name
			}

			_, _ = s.CreateNotification(userID, &goserver.Notification{
				Date:  time.Now(),
				Type:  string(models.NotificationTypeInfo),
				Title: "Reconciliation Invalidated",
				Description: fmt.Sprintf("Financial change to transaction before checkpoint invalidated reconciliation for account %q",
					accountName),
			})
		}
	}
}

func (s *storage) GetDuplicateTransactionIDs(userID, transactionID string) ([]string, error) {
	var duplicates []models.TransactionDuplicate
	err := s.db.Where("user_id = ? AND transaction_id1 = ?", userID, transactionID).Find(&duplicates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get duplicate relationships: %w", err)
	}

	ids := make([]string, len(duplicates))
	for i, d := range duplicates {
		ids[i] = d.TransactionID2.String()
	}
	return ids, nil
}

func (s *storage) AddDuplicateRelationship(userID, transactionID1, transactionID2 string) error {
	id1, err := uuid.Parse(transactionID1)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 1: %w", err)
	}
	id2, err := uuid.Parse(transactionID2)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 2: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Add T1 -> T2
		var d1 models.TransactionDuplicate
		err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).First(&d1).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				d1 = models.TransactionDuplicate{
					UserID:         userID,
					TransactionID1: id1,
					TransactionID2: id2,
				}
				if err := tx.Create(&d1).Error; err != nil {
					return fmt.Errorf("failed to create link T1->T2: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check link T1->T2: %w", err)
			}
		}

		// Add T2 -> T1
		var d2 models.TransactionDuplicate
		err = tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).First(&d2).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				d2 = models.TransactionDuplicate{
					UserID:         userID,
					TransactionID1: id2,
					TransactionID2: id1,
				}
				if err := tx.Create(&d2).Error; err != nil {
					return fmt.Errorf("failed to create link T2->T1: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check link T2->T1: %w", err)
			}
		}

		return nil
	})
}

func (s *storage) RemoveDuplicateRelationship(userID, transactionID1, transactionID2 string) error {
	id1, err := uuid.Parse(transactionID1)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 1: %w", err)
	}
	id2, err := uuid.Parse(transactionID2)
	if err != nil {
		return fmt.Errorf("invalid transaction ID 2: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete link T1->T2: %w", err)
		}
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete link T2->T1: %w", err)
		}
		return nil
	})
}

func (s *storage) ClearDuplicateRelationships(userID, transactionID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.clearDuplicateRelationshipsWithTx(tx, userID, transactionID)
	})
}

func (s *storage) clearDuplicateRelationshipsWithTx(tx *gorm.DB, userID, transactionID string) error {
	id, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	// 1. Find all duplicates linked to this transaction to update them later
	var duplicates []models.TransactionDuplicate
	if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Find(&duplicates).Error; err != nil {
		return fmt.Errorf("failed to find duplicate links: %w", err)
	}

	// 2. Delete bidirectional links
	for _, d := range duplicates {
		if err := tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, d.TransactionID2, id).Delete(&models.TransactionDuplicate{}).Error; err != nil {
			return fmt.Errorf("failed to delete inverse link: %w", err)
		}
	}

	if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Delete(&models.TransactionDuplicate{}).Error; err != nil {
		return fmt.Errorf("failed to delete primary links: %w", err)
	}

	// 3. Sync suspicious reasons for all affected transactions
	// This ensures that if they no longer have duplicates, the flag is removed.
	affectedIDs := []uuid.UUID{id}
	for _, d := range duplicates {
		affectedIDs = append(affectedIDs, d.TransactionID2)
	}

	for _, affectedID := range affectedIDs {
		if err := s.syncDuplicateSuspiciousReason(tx, userID, affectedID); err != nil {
			s.log.Error("Failed to sync suspicious reason", "error", err, "id", affectedID)
		}
	}

	return nil
}

func (s *storage) syncDuplicateSuspiciousReason(tx *gorm.DB, userID string, transactionID uuid.UUID) error {
	// Check if any duplicate links remain for this transaction
	var count int64
	if err := tx.Model(&models.TransactionDuplicate{}).Where("user_id = ? AND transaction_id1 = ?", userID, transactionID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Still has duplicates, keep the reason
	}

	// No more duplicates, remove the reason if present
	var t models.Transaction
	if err := tx.Where("user_id = ? AND id = ?", userID, transactionID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Transaction might have been deleted (e.g. in Merge)
		}
		return err
	}

	newReasons := make([]string, 0)
	reasonsChanged := false
	for _, r := range t.SuspiciousReasons {
		if r == models.DuplicateReason {
			reasonsChanged = true
			continue
		}
		newReasons = append(newReasons, r)
	}

	if reasonsChanged {
		if err := tx.Model(&t).Update("suspicious_reasons", newReasons).Error; err != nil {
			return fmt.Errorf("failed to update suspicious reasons: %w", err)
		}
	}

	return nil
}

// RevalidateDuplicateRelationships re-checks all duplicate links for a transaction.
// If a linked transaction no longer passes IsDuplicate, the link is removed and
// suspicious reasons are synchronized for both transactions.
func (s *storage) RevalidateDuplicateRelationships(userID, transactionID string) error {
	id, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get the transaction
		var t models.Transaction
		if err := tx.Where("user_id = ? AND id = ?", userID, id).First(&t).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil // Transaction doesn't exist, nothing to revalidate
			}
			return err
		}

		// 2. Get all linked duplicates
		var duplicates []models.TransactionDuplicate
		if err := tx.Where("user_id = ? AND transaction_id1 = ?", userID, id).Find(&duplicates).Error; err != nil {
			return err
		}

		// 3. For each link, re-check IsDuplicate
		for _, d := range duplicates {
			var linkedT models.Transaction
			if err := tx.Where("user_id = ? AND id = ?", userID, d.TransactionID2).First(&linkedT).Error; err != nil {
				// Linked transaction doesn't exist, clean up the link
				s.removeDuplicateLinkWithTx(tx, userID, id, d.TransactionID2)
				continue
			}

			if !utils.IsDuplicate(t.Date, t.Movements, linkedT.Date, linkedT.Movements) {
				// No longer duplicates, remove bidirectional link
				s.removeDuplicateLinkWithTx(tx, userID, id, d.TransactionID2)
				// Sync suspicious reasons for both
				s.syncDuplicateSuspiciousReason(tx, userID, id)
				s.syncDuplicateSuspiciousReason(tx, userID, d.TransactionID2)
			}
		}

		return nil
	})
}

func (s *storage) removeDuplicateLinkWithTx(tx *gorm.DB, userID string, id1, id2 uuid.UUID) {
	tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id1, id2).Delete(&models.TransactionDuplicate{})
	tx.Where("user_id = ? AND transaction_id1 = ? AND transaction_id2 = ?", userID, id2, id1).Delete(&models.TransactionDuplicate{})
}

func (s *storage) archiveMergedTransaction(tx *gorm.DB, userID string,
	merged *models.Transaction, keptID uuid.UUID, mergedAt time.Time,
) error {
	archive := models.MergedTransaction{
		ID:                    uuid.New(),
		UserID:                userID,
		KeptTransactionID:     keptID,
		OriginalTransactionID: merged.ID,
		Date:                  merged.Date,
		Description:           merged.Description,
		Place:                 merged.Place,
		Tags:                  merged.Tags,
		PartnerName:           merged.PartnerName,
		PartnerAccount:        merged.PartnerAccount,
		PartnerInternalID:     merged.PartnerInternalID,
		Extra:                 merged.Extra,
		UnprocessedSources:    merged.UnprocessedSources,
		ExternalIDs:           merged.ExternalIDs,
		Movements:             merged.Movements,
		MatcherID:             merged.MatcherID,
		IsAuto:                merged.IsAuto,
		SuspiciousReasons:     merged.SuspiciousReasons,
		MergedAt:              mergedAt,
	}

	if err := tx.Create(&archive).Error; err != nil {
		return fmt.Errorf("failed to create merged transaction archive: %w", err)
	}

	return nil
}

func (s *storage) Backup(destination string) error {
	return s.db.Exec("VACUUM INTO ?", destination).Error
}

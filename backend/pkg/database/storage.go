package database

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/shopspring/decimal"
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

type AuditLogFilter struct {
	EntityType *string
	EntityID   *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Limit      int
	Offset     int
}

type UserStorage interface {
	GetUserID(username string) (string, error)
	GetUser(userID string) (*models.User, error)
	CreateUser(username, password string) (*models.User, error)
	PutUser(user *models.User) error
	GetAllUserIDs() ([]string, error)
}

type AccountStorage interface {
	CreateAccount(userID string, account *goserver.AccountNoId) (goserver.Account, error)
	GetAccounts(userID string) ([]goserver.Account, error)
	GetAccount(userID string, id string) (goserver.Account, error)
	UpdateAccount(userID string, id string, account *goserver.AccountNoId) (goserver.Account, error)
	DeleteAccount(userID string, id string, replaceWithAccountID *string) error
	GetAccountHistory(userID string, accountID string) ([]goserver.Transaction, error)
	GetAccountBalance(userID, accountID, currencyID string) (decimal.Decimal, error)
}

type CurrencyStorage interface {
	CreateCurrency(userID string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	GetCurrencies(userID string) ([]goserver.Currency, error)
	GetCurrency(userID string, id string) (goserver.Currency, error)
	UpdateCurrency(userID string, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	DeleteCurrency(userID string, id string, replaceWithCurrencyID *string) error
}

type TransactionStorage interface {
	GetTransactions(userID string, dateFrom, dateTo time.Time, onlySuspicious bool) ([]goserver.Transaction, error)
	CreateTransaction(userID string, transaction goserver.TransactionNoIdInterface) (goserver.Transaction, error)
	// CreateTransactionsBatch atomically creates multiple transactions in a single database transaction.
	// If any transaction fails to be created, the entire batch is rolled back.
	CreateTransactionsBatch(userID string, transactions []goserver.TransactionNoIdInterface) ([]goserver.Transaction, error)
	UpdateTransaction(
		userID string, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	// UpdateTransactionInternal updates all fields provided in transaction, without preservation logic.
	// Use only for internal operations like auto-matching or unprocessing.
	UpdateTransactionInternal(
		userID string, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	DeleteTransaction(userID string, id string) error
	MergeTransactions(userID, keepID, mergeID string) (goserver.Transaction, error)
	GetTransaction(userID string, id string) (goserver.Transaction, error)
	GetTransactionsIncludingDeleted(userID string, dateFrom, dateTo time.Time) ([]goserver.Transaction, error)
	GetMergedTransactions(userID string) ([]goserver.MergedTransaction, error)
	GetMergedTransaction(userID, originalTransactionID string) (goserver.MergedTransaction, error)
	UnmergeTransaction(userID, id string) error
	GetDuplicateTransactionIDs(userID, transactionID string) ([]string, error)
	AddDuplicateRelationship(userID, transactionID1, transactionID2 string) error
	RemoveDuplicateRelationship(userID, transactionID1, transactionID2 string) error
	ClearDuplicateRelationships(userID, transactionID string) error
	CountUnprocessedTransactionsForAccount(userID, accountID string, ignoreUnprocessedBefore time.Time) (int, error)
	HasTransactionsAfterDate(userID, accountID string, date time.Time) (bool, error)
}

type BankImporterStorage interface {
	GetBankImporters(userID string) ([]goserver.BankImporter, error)
	CreateBankImporter(userID string, bankImporter *goserver.BankImporterNoId) (goserver.BankImporter, error)
	UpdateBankImporter(
		userID string, id string, bankImporter goserver.BankImporterNoIdInterface,
	) (goserver.BankImporter, error)
	DeleteBankImporter(userID string, id string) error
	GetBankImporter(userID string, id string) (goserver.BankImporter, error)
	GetAllBankImporters() ([]ImportInfo, error)
	GetBankImporterFiles(userID string) ([]goserver.BankImporterFile, error)
	GetBankImporterFile(userID string, id string) (models.BankImporterFile, error)
	CreateBankImporterFile(userID string, file *models.BankImporterFile) (goserver.BankImporterFile, error)
	DeleteBankImporterFile(userID string, id string) error
}

type MatcherStorage interface {
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
}

type RateStorage interface {
	SaveCNBRates(rates map[string]decimal.Decimal, day time.Time) error
	GetCNBRates(day time.Time) (map[string]decimal.Decimal, error)
}

type BudgetItemStorage interface {
	CreateBudgetItem(userID string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	GetBudgetItems(userID string) ([]goserver.BudgetItem, error)
	GetBudgetItem(userID string, id string) (goserver.BudgetItem, error)
	UpdateBudgetItem(userID string, id string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	DeleteBudgetItem(userID string, id string) error
}

type ImageStorage interface {
	CreateImage(data []byte, contentType string) (models.Image, error)
	GetImage(id string) (models.Image, error)
	DeleteImage(id string) error
}

type NotificationStorage interface {
	CreateNotification(userID string, notification *goserver.Notification) (goserver.Notification, error)
	GetNotifications(userID string) ([]goserver.Notification, error)
	DeleteNotification(userID string, id string) error
}

type ReconciliationStorage interface {
	GetLatestReconciliation(userID, accountID, currencyID string) (*goserver.Reconciliation, error)
	GetReconciliationsForAccount(userID, accountID string) ([]goserver.Reconciliation, error)
	GetReconciliationsForAccountAndCurrency(userID, accountID, currencyID string) ([]goserver.Reconciliation, error)
	CreateReconciliation(userID string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error)
	InvalidateReconciliation(userID, accountID, currencyID string, fromDate time.Time) error
	GetBulkReconciliationData(userID string) (*BulkReconciliationData, error)
}

type AuditLogStorage interface {
	GetAuditLogs(userID string, filter AuditLogFilter) ([]models.AuditLog, error)
}

type SystemStorage interface {
	Open() error
	Close() error
	Backup(destination string) error
}

//nolint:interfacebloat
type Storage interface {
	UserStorage
	AccountStorage
	CurrencyStorage
	TransactionStorage
	BankImporterStorage
	MatcherStorage
	RateStorage
	BudgetItemStorage
	ImageStorage
	NotificationStorage
	ReconciliationStorage
	AuditLogStorage
	SystemStorage
	WithContext(ctx context.Context) Storage
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
	ctx context.Context
}

func NewStorage(logger *slog.Logger, cfg *config.Config) Storage {
	return &storage{log: logger, db: nil, cfg: cfg, ctx: context.Background()}
}

func (s *storage) WithContext(ctx context.Context) Storage {
	return &storage{
		log: s.log,
		cfg: s.cfg,
		db:  s.db,
		ctx: ctx,
	}
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

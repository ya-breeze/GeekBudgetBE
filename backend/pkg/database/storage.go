package database

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/google/uuid"
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
	FamilyID         uuid.UUID
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
	GetUserByUsername(username string) (*models.User, error)
	GetUser(userID uuid.UUID) (*models.User, error)
	CreateUser(username, passwordHash string, familyID uuid.UUID) (*models.User, error)
	PutUser(user *models.User) error
	GetFamilyByName(name string) (*models.Family, error)
	CreateFamily(name string) (*models.Family, error)
	GetAllFamilyIDs() ([]uuid.UUID, error)
}

type AccountStorage interface {
	CreateAccount(familyID uuid.UUID, account *goserver.AccountNoId) (goserver.Account, error)
	GetAccounts(familyID uuid.UUID) ([]goserver.Account, error)
	GetAccount(familyID uuid.UUID, id string) (goserver.Account, error)
	UpdateAccount(familyID uuid.UUID, id string, account *goserver.AccountNoId) (goserver.Account, error)
	DeleteAccount(familyID uuid.UUID, id string, replaceWithAccountID *string) error
	GetAccountHistory(familyID uuid.UUID, accountID string) ([]goserver.Transaction, error)
	GetAccountBalance(familyID uuid.UUID, accountID, currencyID string) (decimal.Decimal, error)
}

type CurrencyStorage interface {
	CreateCurrency(familyID uuid.UUID, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	GetCurrencies(familyID uuid.UUID) ([]goserver.Currency, error)
	GetCurrency(familyID uuid.UUID, id string) (goserver.Currency, error)
	UpdateCurrency(familyID uuid.UUID, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error)
	DeleteCurrency(familyID uuid.UUID, id string, replaceWithCurrencyID *string) error
}

type TransactionStorage interface {
	GetTransactions(familyID uuid.UUID, dateFrom, dateTo time.Time, onlySuspicious bool) ([]goserver.Transaction, error)
	CreateTransaction(familyID uuid.UUID, transaction goserver.TransactionNoIdInterface) (goserver.Transaction, error)
	// CreateTransactionsBatch atomically creates multiple transactions in a single database transaction.
	// If any transaction fails to be created, the entire batch is rolled back.
	CreateTransactionsBatch(familyID uuid.UUID, transactions []goserver.TransactionNoIdInterface) ([]goserver.Transaction, error)
	UpdateTransaction(
		familyID uuid.UUID, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	// UpdateTransactionInternal updates all fields provided in transaction, without preservation logic.
	// Use only for internal operations like auto-matching or unprocessing.
	UpdateTransactionInternal(
		familyID uuid.UUID, id string, transaction goserver.TransactionNoIdInterface,
	) (goserver.Transaction, error)
	DeleteTransaction(familyID uuid.UUID, id string) error
	MergeTransactions(familyID uuid.UUID, keepID, mergeID string) (goserver.Transaction, error)
	GetTransaction(familyID uuid.UUID, id string) (goserver.Transaction, error)
	GetTransactionsIncludingDeleted(familyID uuid.UUID, dateFrom, dateTo time.Time) ([]goserver.Transaction, error)
	GetMergedTransactions(familyID uuid.UUID) ([]goserver.MergedTransaction, error)
	GetMergedTransaction(familyID uuid.UUID, originalTransactionID string) (goserver.MergedTransaction, error)
	UnmergeTransaction(familyID uuid.UUID, id string) error
	GetDuplicateTransactionIDs(familyID uuid.UUID, transactionID string) ([]string, error)
	AddDuplicateRelationship(familyID uuid.UUID, transactionID1, transactionID2 string) error
	RemoveDuplicateRelationship(familyID uuid.UUID, transactionID1, transactionID2 string) error
	ClearDuplicateRelationships(familyID uuid.UUID, transactionID string) error
	CountUnprocessedTransactionsForAccount(familyID uuid.UUID, accountID string, ignoreUnprocessedBefore time.Time) (int, error)
	HasTransactionsAfterDate(familyID uuid.UUID, accountID string, date time.Time) (bool, error)
}

type BankImporterStorage interface {
	GetBankImporters(familyID uuid.UUID) ([]goserver.BankImporter, error)
	CreateBankImporter(familyID uuid.UUID, bankImporter *goserver.BankImporterNoId) (goserver.BankImporter, error)
	UpdateBankImporter(
		familyID uuid.UUID, id string, bankImporter goserver.BankImporterNoIdInterface,
	) (goserver.BankImporter, error)
	DeleteBankImporter(familyID uuid.UUID, id string) error
	GetBankImporter(familyID uuid.UUID, id string) (goserver.BankImporter, error)
	GetAllBankImporters() ([]ImportInfo, error)
	GetBankImporterFiles(familyID uuid.UUID) ([]goserver.BankImporterFile, error)
	GetBankImporterFile(familyID uuid.UUID, id string) (models.BankImporterFile, error)
	CreateBankImporterFile(familyID uuid.UUID, file *models.BankImporterFile) (goserver.BankImporterFile, error)
	DeleteBankImporterFile(familyID uuid.UUID, id string) error
}

type MatcherStorage interface {
	GetMatchers(familyID uuid.UUID) ([]goserver.Matcher, error)
	GetMatcher(familyID uuid.UUID, id string) (goserver.Matcher, error)
	// Add a single confirmation (true = confirmed, false = rejected) to a matcher
	// This operation is performed atomically and enforces the configured
	// confirmation history maximum length.
	AddMatcherConfirmation(familyID uuid.UUID, id string, confirmed bool) error
	GetMatchersRuntime(familyID uuid.UUID) ([]MatcherRuntime, error)
	GetMatcherRuntime(familyID uuid.UUID, id string) (MatcherRuntime, error)
	CreateMatcherRuntimeFromNoId(m goserver.MatcherNoIdInterface) (MatcherRuntime, error)
	CreateMatcher(familyID uuid.UUID, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	UpdateMatcher(familyID uuid.UUID, id string, matcher goserver.MatcherNoIdInterface) (goserver.Matcher, error)
	DeleteMatcher(familyID uuid.UUID, id string) error
}

type TemplateStorage interface {
	CreateTemplate(familyID uuid.UUID, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
	GetTemplates(familyID uuid.UUID, accountID *string) ([]goserver.TransactionTemplate, error)
	UpdateTemplate(familyID uuid.UUID, id string, t *goserver.TransactionTemplateNoId) (goserver.TransactionTemplate, error)
	DeleteTemplate(familyID uuid.UUID, id string) error
}

type RateStorage interface {
	SaveCNBRates(rates map[string]decimal.Decimal, day time.Time) error
	GetCNBRates(day time.Time) (map[string]decimal.Decimal, error)
}

type BudgetItemStorage interface {
	CreateBudgetItem(familyID uuid.UUID, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	GetBudgetItems(familyID uuid.UUID) ([]goserver.BudgetItem, error)
	GetBudgetItem(familyID uuid.UUID, id string) (goserver.BudgetItem, error)
	UpdateBudgetItem(familyID uuid.UUID, id string, budgetItem *goserver.BudgetItemNoId) (goserver.BudgetItem, error)
	DeleteBudgetItem(familyID uuid.UUID, id string) error
}

type ImageStorage interface {
	CreateImage(data []byte, contentType string) (models.Image, error)
	GetImage(id string) (models.Image, error)
	DeleteImage(id string) error
}

type NotificationStorage interface {
	CreateNotification(familyID uuid.UUID, notification *goserver.Notification) (goserver.Notification, error)
	GetNotifications(familyID uuid.UUID) ([]goserver.Notification, error)
	DeleteNotification(familyID uuid.UUID, id string) error
}

type ReconciliationStorage interface {
	GetLatestReconciliation(familyID uuid.UUID, accountID, currencyID string) (*goserver.Reconciliation, error)
	GetReconciliationsForAccount(familyID uuid.UUID, accountID string) ([]goserver.Reconciliation, error)
	GetReconciliationsForAccountAndCurrency(familyID uuid.UUID, accountID, currencyID string) ([]goserver.Reconciliation, error)
	CreateReconciliation(familyID uuid.UUID, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error)
	InvalidateReconciliation(familyID uuid.UUID, accountID, currencyID string, fromDate time.Time) error
	GetBulkReconciliationData(familyID uuid.UUID) (*BulkReconciliationData, error)
}

type AuditLogStorage interface {
	GetAuditLogs(familyID uuid.UUID, filter AuditLogFilter) ([]models.AuditLog, error)
}

type SystemStorage interface {
	Open() error
	Close() error
	Backup(destination string) error
	GetDB() *gorm.DB
}

//nolint:interfacebloat
type Storage interface {
	UserStorage
	AccountStorage
	CurrencyStorage
	TransactionStorage
	BankImporterStorage
	MatcherStorage
	TemplateStorage
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

func (s *storage) GetDB() *gorm.DB {
	return s.db
}

func (s *storage) Open() error {
	s.log.Info("Opening database", "path", s.cfg.DBPath)
	var err error
	s.db, err = openSqlite(s.log, s.cfg.DBPath, s.cfg.Verbose)
	if err != nil {
		s.log.Error("failed to connect database", "error", err)
		panic("failed to connect database")
	}
	if err := runMigrationIfNeeded(s.log, s.db); err != nil {
		s.log.Error("failed to run kin-core migration", "error", err)
		panic("failed to run kin-core migration")
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

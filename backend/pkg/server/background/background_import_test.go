package background_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

// TestFixture holds commonly created entities for tests
type TestFixture struct {
	Storage             database.Storage
	Logger              *slog.Logger
	UserID              string
	Currency            goserver.Currency
	Account             goserver.Account
	UnprocessedService  *api.UnprocessedTransactionsAPIServiceImpl
	OriginalTransaction *goserver.TransactionNoId
	CreatedTransaction  goserver.Transaction
}

// mustOpenTestStorage creates and opens an in-memory storage for tests and fails
// the test on error.
func mustOpenTestStorage(t *testing.T, logger *slog.Logger) database.Storage {
	t.Helper()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false, MatcherConfirmationHistoryMax: 10}
	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	return storage
}

// setupTestFixture creates a test fixture with common entities (user, currency, account, etc.)
func setupTestFixture(t *testing.T, logger *slog.Logger, username string) *TestFixture {
	t.Helper()
	storage := mustOpenTestStorage(t, logger)

	// Create a user
	createdUser, err := storage.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	userID := createdUser.ID.String()

	// Create a currency
	currency := &goserver.CurrencyNoId{Name: "USD"}
	createdCurrency, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}

	// Create an account
	account := &goserver.AccountNoId{
		Name:        "Test Account",
		Description: "Test account",
	}
	createdAccount, err := storage.CreateAccount(userID, account)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Create unprocessed service
	unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(logger, storage)

	return &TestFixture{
		Storage:            storage,
		Logger:             logger,
		UserID:             userID,
		Currency:           createdCurrency,
		Account:            createdAccount,
		UnprocessedService: unprocessedService,
	}
}

// createMatcher creates a matcher with the given configuration
func (f *TestFixture) createMatcher(t *testing.T, description string, tags []string, history []bool) goserver.Matcher {
	t.Helper()
	matcher := &goserver.MatcherNoId{
		OutputDescription:   description,
		OutputAccountId:     f.Account.Id,
		OutputTags:          tags,
		DescriptionRegExp:   `(?i)test`,
		ConfirmationHistory: history,
	}
	createdMatcher, err := f.Storage.CreateMatcher(f.UserID, matcher)
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}
	return createdMatcher
}

// createTransaction creates an unprocessed transaction
func (f *TestFixture) createTransaction(t *testing.T, description string) goserver.Transaction {
	t.Helper()
	transaction := &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: description,
		Tags:        []string{"unprocessed"},
		Movements: []goserver.Movement{
			{AccountId: "", CurrencyId: f.Currency.Id, Amount: decimal.NewFromInt(-100)},
			{AccountId: f.Account.Id, CurrencyId: f.Currency.Id, Amount: decimal.NewFromInt(100)},
		},
	}
	createdTransaction, err := f.Storage.CreateTransaction(f.UserID, transaction)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	f.OriginalTransaction = transaction
	f.CreatedTransaction = createdTransaction
	return createdTransaction
}

// createBankImporter creates a bank importer to ensure the user is included in getAllUsers
func (f *TestFixture) createBankImporter(t *testing.T) {
	t.Helper()
	bankImporter := &goserver.BankImporterNoId{
		Name:        "Test Importer",
		Description: "Test bank importer",
		Type:        "fio",
	}
	_, err := f.Storage.CreateBankImporter(f.UserID, bankImporter)
	if err != nil {
		t.Fatalf("failed to create bank importer: %v", err)
	}
}

// getUnprocessedTransactions retrieves unprocessed transactions for the user
func (f *TestFixture) getUnprocessedTransactions(t *testing.T) []goserver.UnprocessedTransaction {
	t.Helper()
	unprocessed, _, err := f.UnprocessedService.PrepareUnprocessedTransactions(
		t.Context(), f.UserID, false, "",
	)
	if err != nil {
		t.Fatalf("failed to get unprocessed transactions: %v", err)
	}
	return unprocessed
}

//nolint:cyclop,funlen
func TestProcessUnprocessedTransactionsForAutoConversion(t *testing.T) {
	// Test scenario: Single perfect matcher (≥10 confirmations) → Transaction is auto-converted
	logger := slog.Default()
	fixture := setupTestFixture(t, logger, "testuser")
	defer fixture.Storage.Close()

	// Create matchers with specific confirmation histories
	perfectMatcher := fixture.createMatcher(t,
		"Auto-converted transaction",
		[]string{"auto-converted"},
		[]bool{true, true, true, true, true, true, true, true, true, true}, // 10 perfect confirmations
	)

	fixture.createMatcher(t,
		"Mixed transaction",
		[]string{"mixed"},
		[]bool{true, false, true}, // Mixed history - insufficient
	)

	// Create transaction and bank importer
	createdTransaction := fixture.createTransaction(t, "Test transaction for auto-conversion")
	fixture.createBankImporter(t)

	// Verify initial state: transaction is unprocessed with 2 matchers
	unprocessedBefore := fixture.getUnprocessedTransactions(t)
	if len(unprocessedBefore) != 1 {
		t.Fatalf("expected 1 unprocessed transaction before auto-conversion, got %d", len(unprocessedBefore))
	}
	if unprocessedBefore[0].Transaction.Id != createdTransaction.Id {
		t.Fatalf("unexpected transaction ID before auto-conversion")
	}
	if len(unprocessedBefore[0].Matched) != 2 {
		t.Fatalf("expected 2 matched matchers, got %d", len(unprocessedBefore[0].Matched))
	}

	// Run the auto-conversion process
	background.ProcessUnprocessedTransactionsForAutoConversion(t.Context(), logger, fixture.Storage)

	// Verify transaction was auto-converted
	unprocessedAfter := fixture.getUnprocessedTransactions(t)
	if len(unprocessedAfter) != 0 {
		t.Fatalf("expected 0 unprocessed transactions after auto-conversion, got %d", len(unprocessedAfter))
	}

	// Verify transaction properties were updated correctly
	convertedTransaction, err := fixture.Storage.GetTransaction(fixture.UserID, createdTransaction.Id)
	if err != nil {
		t.Fatalf("failed to get converted transaction: %v", err)
	}

	if convertedTransaction.Description != perfectMatcher.OutputDescription {
		t.Fatalf("expected description '%s', got '%s'",
			perfectMatcher.OutputDescription, convertedTransaction.Description)
	}

	if convertedTransaction.Movements[0].AccountId != perfectMatcher.OutputAccountId {
		t.Fatalf("expected account ID '%s', got '%s'",
			perfectMatcher.OutputAccountId, convertedTransaction.Movements[0].AccountId)
	}

	found := false
	for _, tag := range convertedTransaction.Tags {
		if tag == "auto-converted" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected tag 'auto-converted' to be added to converted transaction")
	}

	// Verify matcher's confirmation history was updated
	updatedMatcher, err := fixture.Storage.GetMatcher(fixture.UserID, perfectMatcher.Id)
	if err != nil {
		t.Fatalf("failed to get updated matcher: %v", err)
	}

	history := updatedMatcher.GetConfirmationHistory()
	if len(history) != 10 { // History is capped at 10 (MatcherConfirmationHistoryMax)
		t.Fatalf("expected 10 confirmations in history, got %d", len(history))
	}
	if !history[len(history)-1] {
		t.Fatalf("expected last confirmation to be true (successful auto-conversion)")
	}
}

//nolint:cyclop,funlen
func TestProcessUnprocessedTransactionsMultiplePerfectMatchers(t *testing.T) {
	// Test scenario: Multiple perfect matchers → Transaction remains unprocessed (ambiguous)
	logger := slog.Default()
	fixture := setupTestFixture(t, logger, "testuser2")
	defer fixture.Storage.Close()

	// Create two matchers with perfect success history
	fixture.createMatcher(t,
		"Auto-converted by matcher 1",
		[]string{"matcher1"},
		[]bool{true, true, true, true, true, true, true, true, true, true}, // 10 perfect confirmations
	)

	fixture.createMatcher(t,
		"Auto-converted by matcher 2",
		[]string{"matcher2"},
		[]bool{true, true, true, true, true, true, true, true, true, true}, // 10 perfect confirmations
	)

	// Create transaction and bank importer
	createdTransaction := fixture.createTransaction(t, "Test transaction with multiple perfect matchers")
	fixture.createBankImporter(t)

	// Verify initial state: transaction is unprocessed
	unprocessedBefore := fixture.getUnprocessedTransactions(t)
	if len(unprocessedBefore) != 1 {
		t.Fatalf("expected 1 unprocessed transaction before auto-conversion, got %d", len(unprocessedBefore))
	}

	// Run the auto-conversion process
	background.ProcessUnprocessedTransactionsForAutoConversion(t.Context(), logger, fixture.Storage)

	// Verify transaction is still unprocessed (multiple perfect matchers = ambiguous)
	unprocessedAfter := fixture.getUnprocessedTransactions(t)
	if len(unprocessedAfter) != 1 {
		t.Fatalf(
			"expected 1 unprocessed transaction after auto-conversion (multiple perfect matchers), got %d",
			len(unprocessedAfter),
		)
	}

	// Verify transaction was not modified
	unchangedTransaction, err := fixture.Storage.GetTransaction(fixture.UserID, createdTransaction.Id)
	if err != nil {
		t.Fatalf("failed to get unchanged transaction: %v", err)
	}

	if unchangedTransaction.Description != fixture.OriginalTransaction.Description {
		t.Fatalf("transaction description was unexpectedly changed")
	}
}

//nolint:cyclop,funlen
func TestProcessUnprocessedTransactionsInsufficientConfirmationHistory(t *testing.T) {
	// Test scenario: Matcher with <10 confirmations → Transaction remains unprocessed (insufficient history)
	logger := slog.Default()
	fixture := setupTestFixture(t, logger, "testuser3")
	defer fixture.Storage.Close()

	// Create a matcher with perfect success history but insufficient confirmations (< 10)
	insufficientMatcher := fixture.createMatcher(t,
		"Should not auto-convert",
		[]string{"insufficient"},
		[]bool{true, true, true, true, true, true, true, true, true}, // Only 9 confirmations
	)

	// Create transaction and bank importer
	createdTransaction := fixture.createTransaction(t, "Test transaction with insufficient matcher history")
	fixture.createBankImporter(t)

	// Verify initial state: transaction is unprocessed
	unprocessedBefore := fixture.getUnprocessedTransactions(t)
	if len(unprocessedBefore) != 1 {
		t.Fatalf("expected 1 unprocessed transaction before auto-conversion, got %d", len(unprocessedBefore))
	}

	// Run the auto-conversion process
	background.ProcessUnprocessedTransactionsForAutoConversion(t.Context(), logger, fixture.Storage)

	// Verify transaction is still unprocessed (insufficient confirmation history)
	unprocessedAfter := fixture.getUnprocessedTransactions(t)
	if len(unprocessedAfter) != 1 {
		t.Fatalf(
			"expected 1 unprocessed transaction after auto-conversion (insufficient history), got %d",
			len(unprocessedAfter),
		)
	}

	// Verify transaction was not modified
	unchangedTransaction, err := fixture.Storage.GetTransaction(fixture.UserID, createdTransaction.Id)
	if err != nil {
		t.Fatalf("failed to get unchanged transaction: %v", err)
	}

	if unchangedTransaction.Description != fixture.OriginalTransaction.Description {
		t.Fatalf("transaction description was unexpectedly changed")
	}

	// Verify matcher's confirmation history was NOT updated (transaction not auto-converted)
	updatedMatcher, err := fixture.Storage.GetMatcher(fixture.UserID, insufficientMatcher.Id)
	if err != nil {
		t.Fatalf("failed to get updated matcher: %v", err)
	}

	history := updatedMatcher.GetConfirmationHistory()
	if len(history) != 9 { // Should remain at 9 confirmations
		t.Fatalf("expected 9 confirmations in history (unchanged), got %d", len(history))
	}
}

// MockStorage wraps database.Storage to intercept calls for testing
type MockStorage struct {
	database.Storage
	OnGetAllBankImporters func() ([]database.ImportInfo, error)
}

func (m *MockStorage) GetAllBankImporters() ([]database.ImportInfo, error) {
	if m.OnGetAllBankImporters != nil {
		return m.OnGetAllBankImporters()
	}
	return m.Storage.GetAllBankImporters()
}

func TestStartBankImporters_RunsOnStartup(t *testing.T) {
	logger := slog.Default()
	fixture := setupTestFixture(t, logger, "testuser_import")
	defer fixture.Storage.Close()

	// Channel to signal that GetAllBankImporters was called
	called := make(chan struct{}, 1)

	// Create mock storage
	mock := &MockStorage{
		Storage: fixture.Storage,
		OnGetAllBankImporters: func() ([]database.ImportInfo, error) {
			select {
			case called <- struct{}{}:
			default:
			}
			return fixture.Storage.GetAllBankImporters()
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	forcedImports := make(chan common.ForcedImport)

	// Start the background task
	cfg := &config.Config{}
	done := background.StartBankImporters(ctx, logger, mock, cfg, forcedImports)

	// Verify it runs immediately (within reasonable time)
	select {
	case <-called:
		// Success: called immediately on startup
	case <-time.After(5 * time.Second):
		t.Fatal("StartBankImporters did not run GetAllBankImporters on startup")
	}

	// Clean up
	cancel()
	select {
	case <-done:
		// Success: stopped correctly
	case <-time.After(5 * time.Second):
		t.Fatal("StartBankImporters did not stop after context cancellation")
	}
}

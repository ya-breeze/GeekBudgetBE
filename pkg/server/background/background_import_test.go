package background

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

func TestProcessUnprocessedTransactionsForAutoConversion(t *testing.T) {
	// Create in-memory database for testing
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false, MatcherConfirmationHistoryMax: 10}
	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer storage.Close()

	// Create a user
	createdUser, err := storage.CreateUser("testuser", "password123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	userID := createdUser.ID.String()

	// Create a currency
	currency := &goserver.CurrencyNoId{
		Name: "USD",
	}
	createdCurrency, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}

	// Create an account
	account := &goserver.AccountNoId{
		Name:        "Test Account",
		Description: "Test account for auto-conversion",
	}
	createdAccount, err := storage.CreateAccount(userID, account)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Create a matcher with perfect success history (all true confirmations)
	matcher := &goserver.MatcherNoId{
		Name:                "Perfect Matcher",
		OutputDescription:   "Auto-converted transaction",
		OutputAccountId:     createdAccount.Id,
		OutputTags:          []string{"auto-converted"},
		DescriptionRegExp:   `(?i)test`,
		ConfirmationHistory: []bool{true, true, true}, // Perfect success history
	}
	createdMatcher, err := storage.CreateMatcher(userID, matcher)
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	// Create another matcher with mixed history
	mixedMatcher := &goserver.MatcherNoId{
		Name:                "Mixed Matcher",
		OutputDescription:   "Mixed transaction",
		OutputAccountId:     createdAccount.Id,
		OutputTags:          []string{"mixed"},
		DescriptionRegExp:   `(?i)test`,
		ConfirmationHistory: []bool{true, false, true}, // Mixed history
	}
	_, err = storage.CreateMatcher(userID, mixedMatcher)
	if err != nil {
		t.Fatalf("failed to create mixed matcher: %v", err)
	}

	// Create an unprocessed transaction that should match both matchers
	transaction := &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: "Test transaction for auto-conversion",
		Tags:        []string{"unprocessed"},
		Movements: []goserver.Movement{
			{AccountId: "", CurrencyId: createdCurrency.Id, Amount: -100.0}, // Empty account - should be filled by matcher
			{AccountId: createdAccount.Id, CurrencyId: createdCurrency.Id, Amount: 100.0},
		},
	}
	createdTransaction, err := storage.CreateTransaction(userID, transaction)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	// Create a bank importer to ensure the user is included in getAllUsers
	bankImporter := &goserver.BankImporterNoId{
		Name:        "Test Importer",
		Description: "Test bank importer",
		Type:        "fio",
	}
	_, err = storage.CreateBankImporter(userID, bankImporter)
	if err != nil {
		t.Fatalf("failed to create bank importer: %v", err)
	}

	// Before auto-conversion: verify transaction exists and is unprocessed
	unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(logger, storage)
	unprocessedBefore, _, err := unprocessedService.PrepareUnprocessedTransactions(
		context.Background(), userID, false, "",
	)
	if err != nil {
		t.Fatalf("failed to get unprocessed transactions before: %v", err)
	}
	if len(unprocessedBefore) != 1 {
		t.Fatalf("expected 1 unprocessed transaction before auto-conversion, got %d", len(unprocessedBefore))
	}
	if unprocessedBefore[0].Transaction.Id != createdTransaction.Id {
		t.Fatalf("unexpected transaction ID before auto-conversion")
	}

	// Verify that the unprocessed transaction has 2 matchers (perfect and mixed)
	if len(unprocessedBefore[0].Matched) != 2 {
		t.Fatalf("expected 2 matched matchers, got %d", len(unprocessedBefore[0].Matched))
	}

	// Run the auto-conversion process
	ctx := context.Background()
	processUnprocessedTransactionsForAutoConversion(ctx, logger, storage)

	// After auto-conversion: verify transaction was converted using the perfect matcher
	unprocessedAfter, _, err := unprocessedService.PrepareUnprocessedTransactions(
		ctx, userID, false, "",
	)
	if err != nil {
		t.Fatalf("failed to get unprocessed transactions after: %v", err)
	}

	// The transaction should be converted (no longer unprocessed)
	if len(unprocessedAfter) != 0 {
		t.Fatalf("expected 0 unprocessed transactions after auto-conversion, got %d", len(unprocessedAfter))
	}

	// Verify the transaction was actually converted with the correct properties
	convertedTransaction, err := storage.GetTransaction(userID, createdTransaction.Id)
	if err != nil {
		t.Fatalf("failed to get converted transaction: %v", err)
	}

	// Check that the transaction was updated with matcher's output
	if convertedTransaction.Description != createdMatcher.OutputDescription {
		t.Fatalf("expected description '%s', got '%s'", 
			createdMatcher.OutputDescription, convertedTransaction.Description)
	}

	// Check that the empty account was filled
	if convertedTransaction.Movements[0].AccountId != createdMatcher.OutputAccountId {
		t.Fatalf("expected account ID '%s', got '%s'", 
			createdMatcher.OutputAccountId, convertedTransaction.Movements[0].AccountId)
	}

	// Check that output tags were added
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

	// Verify that the matcher's confirmation history was updated
	updatedMatcher, err := storage.GetMatcher(userID, createdMatcher.Id)
	if err != nil {
		t.Fatalf("failed to get updated matcher: %v", err)
	}

	history := updatedMatcher.GetConfirmationHistory()
	if len(history) != 4 { // Original 3 + 1 new confirmation
		t.Fatalf("expected 4 confirmations in history, got %d", len(history))
	}
	if !history[len(history)-1] { // Last confirmation should be true
		t.Fatalf("expected last confirmation to be true (successful auto-conversion)")
	}
}

func TestProcessUnprocessedTransactionsMultiplePerfectMatchers(t *testing.T) {
	// Create in-memory database for testing
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false, MatcherConfirmationHistoryMax: 10}
	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer storage.Close()

	// Create a user
	createdUser, err := storage.CreateUser("testuser2", "password123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	userID := createdUser.ID.String()

	// Create a currency
	currency := &goserver.CurrencyNoId{
		Name: "USD",
	}
	createdCurrency, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}

	// Create an account
	account := &goserver.AccountNoId{
		Name:        "Test Account",
		Description: "Test account for multiple perfect matchers",
	}
	createdAccount, err := storage.CreateAccount(userID, account)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Create two matchers with perfect success history
	matcher1 := &goserver.MatcherNoId{
		Name:                "Perfect Matcher 1",
		OutputDescription:   "Auto-converted by matcher 1",
		OutputAccountId:     createdAccount.Id,
		OutputTags:          []string{"matcher1"},
		DescriptionRegExp:   `(?i)test`,
		ConfirmationHistory: []bool{true, true, true}, // Perfect success history
	}
	_, err = storage.CreateMatcher(userID, matcher1)
	if err != nil {
		t.Fatalf("failed to create matcher1: %v", err)
	}

	matcher2 := &goserver.MatcherNoId{
		Name:                "Perfect Matcher 2",
		OutputDescription:   "Auto-converted by matcher 2",
		OutputAccountId:     createdAccount.Id,
		OutputTags:          []string{"matcher2"},
		DescriptionRegExp:   `(?i)test`,
		ConfirmationHistory: []bool{true, true}, // Perfect success history
	}
	_, err = storage.CreateMatcher(userID, matcher2)
	if err != nil {
		t.Fatalf("failed to create matcher2: %v", err)
	}

	// Create an unprocessed transaction that should match both perfect matchers
	transaction := &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: "Test transaction with multiple perfect matchers",
		Tags:        []string{"unprocessed"},
		Movements: []goserver.Movement{
			{AccountId: "", CurrencyId: createdCurrency.Id, Amount: -100.0},
			{AccountId: createdAccount.Id, CurrencyId: createdCurrency.Id, Amount: 100.0},
		},
	}
	createdTransaction, err := storage.CreateTransaction(userID, transaction)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	// Create a bank importer to ensure the user is included in getAllUsers
	bankImporter := &goserver.BankImporterNoId{
		Name:        "Test Importer",
		Description: "Test bank importer",
		Type:        "fio",
	}
	_, err = storage.CreateBankImporter(userID, bankImporter)
	if err != nil {
		t.Fatalf("failed to create bank importer: %v", err)
	}

	// Before auto-conversion: verify transaction exists and is unprocessed
	unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(logger, storage)
	unprocessedBefore, _, err := unprocessedService.PrepareUnprocessedTransactions(
		context.Background(), userID, false, "",
	)
	if err != nil {
		t.Fatalf("failed to get unprocessed transactions before: %v", err)
	}
	if len(unprocessedBefore) != 1 {
		t.Fatalf("expected 1 unprocessed transaction before auto-conversion, got %d", len(unprocessedBefore))
	}

	// Run the auto-conversion process
	ctx := context.Background()
	processUnprocessedTransactionsForAutoConversion(ctx, logger, storage)

	// After auto-conversion: verify transaction is still unprocessed
	// (because there are multiple perfect matchers)
	unprocessedAfter, _, err := unprocessedService.PrepareUnprocessedTransactions(
		ctx, userID, false, "",
	)
	if err != nil {
		t.Fatalf("failed to get unprocessed transactions after: %v", err)
	}

	// The transaction should still be unprocessed due to multiple perfect matchers
	if len(unprocessedAfter) != 1 {
		t.Fatalf("expected 1 unprocessed transaction after auto-conversion (multiple perfect matchers), got %d", len(unprocessedAfter))
	}

	// Verify the transaction was not modified
	unchangedTransaction, err := storage.GetTransaction(userID, createdTransaction.Id)
	if err != nil {
		t.Fatalf("failed to get unchanged transaction: %v", err)
	}

	// Transaction should remain unchanged
	if unchangedTransaction.Description != transaction.Description {
		t.Fatalf("transaction description was unexpectedly changed")
	}
}

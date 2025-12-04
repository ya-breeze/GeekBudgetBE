package test

import (
	"log/slog"
	"time"

	"github.com/dusted-go/logging/prettylog"
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CreateTestLogger() *slog.Logger {
	return slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
}

func PrepareAccounts() []goserver.Account {
	return []goserver.Account{
		{Id: "0", Name: "Cash", Type: constants.AccountAsset},
		{Id: "1", Name: "Bank", Type: constants.AccountAsset},
		{Id: "2", Name: "Food", Type: constants.AccountExpense},
		{Id: "3", Name: "Salary", Type: constants.AccountIncome},
		{Id: "4", Name: "Groceries", Type: constants.AccountExpense},
	}
}

func PrepareCurrencies() []goserver.Currency {
	return []goserver.Currency{
		{Id: "0", Name: "USD"},
		{Id: "1", Name: "EUR"},
		{Id: "2", Name: "CZK"},
	}
}

func PrepareTransactions(accounts []goserver.Account, currencies []goserver.Currency) []goserver.Transaction {
	return []goserver.Transaction{
		PrepareTransaction(
			"salary 9",
			time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), 2000,
			currencies[0].Id, accounts[1].Id, accounts[3].Id),
		PrepareTransaction(
			"food 9",
			time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), 200,
			currencies[0].Id, accounts[2].Id, accounts[0].Id),
		PrepareTransaction(
			"food 9.2",
			time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC), 250,
			currencies[0].Id, accounts[2].Id, accounts[0].Id),
		PrepareTransaction(
			"groceries 9",
			time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), 300,
			currencies[0].Id, accounts[4].Id, accounts[1].Id),

		PrepareTransaction(
			"salary 10",
			time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), 1000,
			currencies[0].Id, accounts[1].Id, accounts[3].Id),
		PrepareTransaction(
			"food 10",
			time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC), 10,
			currencies[0].Id, accounts[2].Id, accounts[0].Id),
		PrepareTransaction(
			"groceries 10",
			time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC), 100,
			currencies[0].Id, accounts[4].Id, accounts[1].Id),
		PrepareTransaction(
			"groceries 10.2",
			time.Date(2024, 10, 18, 0, 0, 0, 0, time.UTC), 150,
			currencies[0].Id, accounts[4].Id, accounts[1].Id),
	}
}

func PrepareTransaction(
	description string, date time.Time, amount float64, currencyID, accountPlus, accountMinus string,
) goserver.Transaction {
	return goserver.Transaction{
		Id:          uuid.New().String(),
		Description: description,
		Date:        date,
		Movements: []goserver.Movement{
			{Amount: amount, CurrencyId: currencyID, AccountId: accountPlus},
			{Amount: -amount, CurrencyId: currencyID, AccountId: accountMinus},
		},
	}
}

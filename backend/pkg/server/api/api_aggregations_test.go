package api_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Aggregation API", func() {
	log := test.CreateTestLogger()
	ctx := context.Background()
	accounts := test.PrepareAccounts()
	currencies := test.PrepareCurrencies()
	transactions := test.PrepareTransactions(accounts, currencies)
	dateFrom := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)

	// Create currency map for tests (ID -> Name)
	currencyMap := make(map[string]string)
	for _, currency := range currencies {
		currencyMap[currency.Id] = currency.Name
	}

	It("aggregate expenses", func() {
		sut := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			"", nil, currencyMap,
			func(a goserver.Account) bool { return a.Type == "expense" },
			log)
		Expect(sut.From.UnixMilli()).To(Equal(time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).UnixMilli()))
		Expect(sut.To.UnixMilli()).To(Equal(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC).UnixMilli()))

		Expect(sut.Intervals).To(HaveLen(2))

		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id))
		Expect(sut.Currencies[0].Accounts).To(HaveLen(2))
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id))

		Expect(sut.Currencies[0].Accounts[0].AccountId).To(Equal(accounts[2].Id))
		Expect(sut.Currencies[0].Accounts[0].Amounts).To(HaveLen(2))
		Expect(sut.Currencies[0].Accounts[0].Amounts[0].Equal(decimal.NewFromFloat(450.0))).To(BeTrue())
		Expect(sut.Currencies[0].Accounts[0].Amounts[1].Equal(decimal.NewFromFloat(10.0))).To(BeTrue())

		Expect(sut.Currencies[0].Accounts[1].AccountId).To(Equal(accounts[4].Id))
		Expect(sut.Currencies[0].Accounts[1].Amounts).To(HaveLen(2))
		Expect(sut.Currencies[0].Accounts[1].Amounts[0].Equal(decimal.NewFromFloat(300.0))).To(BeTrue())
		Expect(sut.Currencies[0].Accounts[1].Amounts[1].Equal(decimal.NewFromFloat(250.0))).To(BeTrue())
	})

	It("handles edge cases for currency conversion", func() {
		// Test with nil currenciesRatesFetcher - should use original currencies
		sut := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[1].Id, nil, currencyMap, // outputCurrencyID="1" (EUR), but nil fetcher
			func(a goserver.Account) bool { return a.Type == "expense" },
			log)

		// Should still group by original currency (USD = "0") since conversion fails
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id)) // Still USD

		// Test with empty outputCurrencyID - should use original currencies
		sut2 := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			"", nil, currencyMap, // empty outputCurrencyID
			func(a goserver.Account) bool { return a.Type == "expense" },
			log)

		// Should group by original currency
		Expect(sut2.Currencies).To(HaveLen(1))
		Expect(sut2.Currencies[0].CurrencyId).To(Equal(currencies[0].Id)) // USD

		// Test with same currency but nil fetcher - should still warn about nil fetcher
		sut3 := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[0].Id, nil, currencyMap, // same currency as transactions (USD), but nil fetcher
			func(a goserver.Account) bool { return a.Type == "expense" },
			log)

		// Should group by original currency since fetcher is nil
		Expect(sut3.Currencies).To(HaveLen(1))
		Expect(sut3.Currencies[0].CurrencyId).To(Equal(currencies[0].Id)) // USD
	})

	It("successfully converts currencies with real exchange rates", func() {
		ctrl := gomock.NewController(GinkgoT())
		defer ctrl.Finish()
		mockStorage := mocks.NewMockStorage(ctrl)

		// Mock storage to return realistic CNB exchange rates (rates to CZK)
		// Use the actual currency names from test data
		mockRates := map[string]decimal.Decimal{
			currencies[0].Name: decimal.NewFromFloat(22.5), // 1 USD = 22.5 CZK
			currencies[1].Name: decimal.NewFromFloat(25.0), // 1 EUR = 25.0 CZK
			currencies[2].Name: decimal.NewFromFloat(1.0),  // 1 CZK = 1 CZK (base currency)
		}

		// Expect storage calls for currency conversion
		mockStorage.EXPECT().
			GetCNBRates(gomock.Any()).
			Return(mockRates, nil).
			AnyTimes()

		fetcher := common.NewCurrenciesRatesFetcher(log, mockStorage)

		// Create transactions in different currencies
		mixedTransactions := []goserver.Transaction{
			test.PrepareTransaction("USD expense", time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), decimal.NewFromInt(100),
				currencies[0].Id, accounts[2].Id, accounts[0].Id), // 100 USD
			test.PrepareTransaction("EUR expense", time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC), decimal.NewFromInt(50),
				currencies[1].Id, accounts[2].Id, accounts[0].Id), // 50 EUR
		}

		sut := api.Aggregate(
			ctx, accounts, mixedTransactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[2].Id, // Convert everything to CZK (use ID, not Name)
			fetcher, currencyMap,
			func(a goserver.Account) bool { return a.Type == "expense" },
			log)

		// Should have only one currency (CZK) after conversion
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[2].Id)) // CZK

		// Check that amounts were converted: 100 USD * 22.5 + 50 EUR * 25.0 = 2250 + 1250 = 3500 CZK
		totalAmount := decimal.Zero
		for _, account := range sut.Currencies[0].Accounts {
			for _, amount := range account.Amounts {
				totalAmount = totalAmount.Add(amount)
			}
		}
		expectedAmount := decimal.NewFromInt(3500)
		Expect(totalAmount.Sub(expectedAmount).Abs().LessThanOrEqual(decimal.NewFromFloat(1.0))).To(BeTrue())
	})

	It("aggregate balances (assets)", func() {
		sut := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			"", nil, currencyMap,
			func(a goserver.Account) bool { return a.Type == "asset" },
			log)

		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id))
		// We expect accounts "Cash" (0) and "Bank" (1)
		Expect(sut.Currencies[0].Accounts).To(HaveLen(2))

		// Find "Bank" account (Id="1")
		var bankAcc goserver.AccountAggregation
		for _, acc := range sut.Currencies[0].Accounts {
			if acc.AccountId == "1" {
				bankAcc = acc
			}
		}
		Expect(bankAcc.AccountId).To(Equal("1"))

		// History:
		// Sep: +2000 (salary), -300 (groceries). Net: +1700.
		// Oct: +1000 (salary), -100 (groceries), -150 (groceries). Net: +750.
		Expect(bankAcc.Amounts).To(HaveLen(2))
		Expect(bankAcc.Amounts[0].Equal(decimal.NewFromFloat(1700.0))).To(BeTrue())
		Expect(bankAcc.Amounts[1].Equal(decimal.NewFromFloat(750.0))).To(BeTrue())

		// Find "Cash" account (Id="0")
		// Sep: -200 (food), -250 (food). Net: -450.
		// Oct: -10 (food). Net: -10.
		var cashAcc goserver.AccountAggregation
		for _, acc := range sut.Currencies[0].Accounts {
			if acc.AccountId == "0" {
				cashAcc = acc
			}
		}
		Expect(cashAcc.AccountId).To(Equal("0"))
		Expect(cashAcc.Amounts[0].Equal(decimal.NewFromFloat(-450.0))).To(BeTrue())
		Expect(cashAcc.Amounts[1].Equal(decimal.NewFromFloat(-10.0))).To(BeTrue())
	})

	It("demonstrates compounding precision errors (baseline)", func() {
		// Create 1000 transactions of 0.1
		manyTransactions := make([]goserver.Transaction, 1000)
		for i := 0; i < 1000; i++ {
			manyTransactions[i] = test.PrepareTransaction("small", dateFrom.Add(time.Duration(i)*time.Minute), decimal.NewFromFloat(0.1), currencies[0].Id, accounts[0].Id, accounts[1].Id)
		}

		sut := api.Aggregate(
			ctx, accounts, manyTransactions, dateFrom, dateTo, utils.GranularityMonth,
			"", nil, currencyMap,
			func(a goserver.Account) bool { return a.Type == "asset" },
			log)

		total := decimal.Zero
		for _, acc := range sut.Currencies[0].Accounts {
			if acc.AccountId == accounts[0].Id {
				for _, amount := range acc.Amounts {
					total = total.Add(amount)
				}
			}
		}
		// 1000 * 0.1 should be 100.0 exactly
		expected := decimal.NewFromInt(100)
		Expect(total.Equal(expected)).To(BeTrue())
		if !total.Equal(expected) {
			log.Info("Compounded error detected", "expected", expected, "actual", total, "diff", total.Sub(expected))
		}
	})
})

package api_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		Expect(sut.Currencies[0].Accounts[0].Amounts[0]).To(Equal(450.0))
		Expect(sut.Currencies[0].Accounts[0].Amounts[1]).To(Equal(10.0))

		Expect(sut.Currencies[0].Accounts[1].AccountId).To(Equal(accounts[4].Id))
		Expect(sut.Currencies[0].Accounts[1].Amounts).To(HaveLen(2))
		Expect(sut.Currencies[0].Accounts[1].Amounts[0]).To(Equal(300.0))
		Expect(sut.Currencies[0].Accounts[1].Amounts[1]).To(Equal(250.0))
	})

	It("handles edge cases for currency conversion", func() {
		// Test with nil currenciesRatesFetcher - should use original currencies
		sut := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[1].Name, nil, currencyMap, // outputCurrencyName="EUR", but nil fetcher
			log)

		// Should still group by original currency (USD = "0") since conversion fails
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id)) // Still USD

		// Test with empty outputCurrencyName - should use original currencies
		sut2 := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			"", nil, currencyMap, // empty outputCurrencyName
			log)

		// Should group by original currency
		Expect(sut2.Currencies).To(HaveLen(1))
		Expect(sut2.Currencies[0].CurrencyId).To(Equal(currencies[0].Id)) // USD

		// Test with same currency but nil fetcher - should still warn about nil fetcher
		sut3 := api.Aggregate(
			ctx, accounts, transactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[0].Name, nil, currencyMap, // same currency as transactions (USD), but nil fetcher
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
		mockRates := map[string]float64{
			currencies[0].Name: 22.5, // 1 USD = 22.5 CZK
			currencies[1].Name: 25.0, // 1 EUR = 25.0 CZK
			currencies[2].Name: 1.0,  // 1 CZK = 1 CZK (base currency)
		}

		// Expect storage calls for currency conversion
		mockStorage.EXPECT().
			GetCNBRates(gomock.Any()).
			Return(mockRates, nil).
			AnyTimes()

		fetcher := common.NewCurrenciesRatesFetcher(log, mockStorage)

		// Create transactions in different currencies
		mixedTransactions := []goserver.Transaction{
			test.PrepareTransaction("USD expense", time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), 100,
				currencies[0].Id, accounts[2].Id, accounts[0].Id), // 100 USD
			test.PrepareTransaction("EUR expense", time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC), 50,
				currencies[1].Id, accounts[2].Id, accounts[0].Id), // 50 EUR
		}

		sut := api.Aggregate(
			ctx, accounts, mixedTransactions, dateFrom, dateTo, utils.GranularityMonth,
			currencies[2].Name, // Convert everything to CZK
			fetcher, currencyMap,
			log)

		// Should have only one currency (CZK) after conversion
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[2].Id)) // CZK

		// Check that amounts were converted: 100 USD * 22.5 + 50 EUR * 25.0 = 2250 + 1250 = 3500 CZK
		totalAmount := 0.0
		for _, account := range sut.Currencies[0].Accounts {
			for _, amount := range account.Amounts {
				totalAmount += amount
			}
		}
		Expect(totalAmount).To(BeNumerically("~", 3500.0, 1.0)) // Allow small floating point variance
	})
})

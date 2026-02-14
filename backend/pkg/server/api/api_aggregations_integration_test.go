package api_test

import (
	"context"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Soft Delete Aggregation Integration", func() {
	var (
		st     database.Storage
		sut    *api.AggregationsAPIServiceImpl
		ctx    context.Context
		log    = test.CreateTestLogger()
		userID = "test-user"
		usdID  = "usd-id"
	)

	BeforeEach(func() {
		cfg := &config.Config{DBPath: ":memory:", Verbose: false}
		st = database.NewStorage(log, cfg)
		Expect(st.Open()).To(Succeed())
		sut = api.NewAggregationsAPIServiceImpl(log, st)
		ctx = context.WithValue(context.Background(), constants.UserIDKey, userID)

		// Setup base data
		_, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "USD"})
		Expect(err).ToNot(HaveOccurred())
		curs, _ := st.GetCurrencies(userID)
		usdID = curs[0].Id
	})

	AfterEach(func() {
		st.Close()
	})

	It("ignores soft-deleted transactions in expenses and balances", func() {
		// 1. Create asset and expense accounts
		assetAcc, err := st.CreateAccount(userID, &goserver.AccountNoId{
			Name: "Bank",
			Type: "asset",
			BankInfo: goserver.BankAccountInfo{
				Balances: []goserver.BankAccountInfoBalancesInner{
					{CurrencyId: usdID, OpeningBalance: decimal.NewFromFloat(1000.0)},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		expenseAcc, err := st.CreateAccount(userID, &goserver.AccountNoId{
			Name: "Groceries",
			Type: "expense",
		})
		Expect(err).ToNot(HaveOccurred())

		date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
		dateFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

		// 2. Create a transaction (expense 100 USD)
		t1, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        date,
			Description: "Expense to keep",
			Movements: []goserver.Movement{
				{AccountId: assetAcc.Id, CurrencyId: usdID, Amount: decimal.NewFromFloat(-100.0)},
				{AccountId: expenseAcc.Id, CurrencyId: usdID, Amount: decimal.NewFromFloat(100.0)},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(t1.Id).ToNot(BeEmpty())

		// 3. Create another transaction (expense 500 USD) and soft-delete it
		t2, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        date,
			Description: "Expense to delete",
			Movements: []goserver.Movement{
				{AccountId: assetAcc.Id, CurrencyId: usdID, Amount: decimal.NewFromFloat(-500.0)},
				{AccountId: expenseAcc.Id, CurrencyId: usdID, Amount: decimal.NewFromFloat(500.0)},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(st.DeleteTransaction(userID, t2.Id)).To(Succeed())

		// 4. Verify Expenses - should be 100, not 600
		expResp, err := sut.GetExpenses(ctx, dateFrom, dateTo, usdID, "month", false)
		Expect(err).ToNot(HaveOccurred())
		expAgg := expResp.Body.(*goserver.Aggregation)
		Expect(expAgg.Currencies[0].Accounts[0].Total.Equal(decimal.NewFromFloat(100.0))).To(BeTrue())

		// 5. Verify Balances - should be 900 (1000 - 100), not 400 (1000 - 100 - 500)
		balResp, err := sut.GetBalances(ctx, dateFrom, dateTo, usdID, false)
		Expect(err).ToNot(HaveOccurred())
		balAgg := balResp.Body.(*goserver.Aggregation)

		// Find bank account in balances
		var bankTotal decimal.Decimal
		found := false
		for _, cur := range balAgg.Currencies {
			for _, acc := range cur.Accounts {
				if acc.AccountId == assetAcc.Id {
					bankTotal = acc.Total
					found = true
				}
			}
		}
		Expect(found).To(BeTrue())
		Expect(bankTotal.Equal(decimal.NewFromFloat(900.0))).To(BeTrue())
	})
})

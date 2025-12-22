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
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Balances Aggregation API", func() {
	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         *api.AggregationsAPIServiceImpl
		ctx         context.Context
		log         = test.CreateTestLogger()
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewAggregationsAPIServiceImpl(log, mockStorage)
		// Inject UserID into context
		ctx = context.WithValue(context.Background(), common.UserIDKey, "user-1")
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("calculates cumulative balances correctly", func() {
		userID := "user-1"
		usdID := "usd-id"
		accountID := "acc-1"
		dateFrom := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC) // 2 months: Sep, Oct

		// 1. Prepare Currencies
		currencies := []goserver.Currency{
			{Id: usdID, Name: "USD"},
		}
		mockStorage.EXPECT().GetCurrencies(userID).Return(currencies, nil).AnyTimes()

		// 2. Prepare Accounts
		// Account has Opening Balance of 100 USD
		accounts := []goserver.Account{
			{
				Id:   accountID,
				Name: "Bank",
				Type: "asset",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{CurrencyId: usdID, OpeningBalance: 100.0},
					},
				},
			},
		}

		mockStorage.EXPECT().GetAccounts(userID).Return(accounts, nil).AnyTimes()

		// 3. Prepare Transactions
		// T1: Past transaction (Jan 2024), +50 USD
		// T2: In-range transaction (Oct 2024), -20 USD
		t1Date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		t2Date := time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC)

		t1 := goserver.Transaction{
			Date: t1Date,
			Movements: []goserver.Movement{
				{AccountId: accountID, CurrencyId: usdID, Amount: 50.0},
			},
		}
		t2 := goserver.Transaction{
			Date: t2Date,
			Movements: []goserver.Movement{
				{AccountId: accountID, CurrencyId: usdID, Amount: -20.0},
			},
		}

		// Expect GetTransactions checks
		// 1. Range query (Sep-Nov) -> Returns T2
		mockStorage.EXPECT().
			GetTransactions(userID, dateFrom, dateTo).
			Return([]goserver.Transaction{t2}, nil)

		// 2. Past query (2000 - Sep) -> Returns T1
		beginningOfTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		mockStorage.EXPECT().
			GetTransactions(userID, beginningOfTime, dateFrom).
			Return([]goserver.Transaction{t1}, nil)

		// Call SUT
		agg, err := sut.GetAggregatedBalances(ctx, userID, dateFrom, dateTo, usdID, false)
		Expect(err).ToNot(HaveOccurred())

		// Verify Results
		// Currencies
		Expect(agg.Currencies).To(HaveLen(1))
		Expect(agg.Currencies[0].CurrencyId).To(Equal(usdID))

		// Accounts
		Expect(agg.Currencies[0].Accounts).To(HaveLen(1))
		accAgg := agg.Currencies[0].Accounts[0]
		Expect(accAgg.AccountId).To(Equal(accountID))

		// Intervals: Sep, Oct
		Expect(accAgg.Amounts).To(HaveLen(2))

		// Initial Balance = 100 (Opening) + 50 (Past) = 150
		// Sep: No transactions. Cumulative = 150.
		Expect(accAgg.Amounts[0]).To(Equal(150.0))

		// Oct: -20 transaction. Cumulative = 150 - 20 = 130.
		Expect(accAgg.Amounts[1]).To(Equal(130.0))
	})
})

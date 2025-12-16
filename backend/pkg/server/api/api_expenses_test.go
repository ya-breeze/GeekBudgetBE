package api_test

import (
	"context"
	"net/http"
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

var _ = Describe("Expenses Aggregation API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), common.UserIDKey, "user1")

	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         goserver.AggregationsAPIService
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewAggregationsAPIServiceImpl(log, mockStorage)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("calculates yearly expenses correctly", func() {
		// Time setup
		// Start of current year to ensure consistent filtering
		// Actually, let's use a fixed date range
		from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		accountA := goserver.Account{
			Id:   "acc-a",
			Name: "Account A",
			Type: "expense",
			BankInfo: goserver.BankAccountInfo{
				Balances: []goserver.BankAccountInfoBalancesInner{{CurrencyId: "USD"}},
			},
		}

		// Transactions in Jan, June, Dec
		transactions := []goserver.Transaction{
			{
				Date: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{AccountId: "acc-a", Amount: 100.0, CurrencyId: "USD"},
				},
			},
			{
				Date: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{AccountId: "acc-a", Amount: 200.0, CurrencyId: "USD"},
				},
			},
			{
				Date: time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{AccountId: "acc-a", Amount: 300.0, CurrencyId: "USD"},
				},
			},
		}

		mockStorage.EXPECT().GetAccounts("user1").Return([]goserver.Account{accountA}, nil)
		mockStorage.EXPECT().GetTransactions("user1", from, to).Return(transactions, nil)
		mockStorage.EXPECT().GetCurrencies("user1").Return([]goserver.Currency{{Id: "USD", Name: "USD"}}, nil)

		resp, err := sut.GetExpenses(ctx, from, to, "", "year")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Code).To(Equal(http.StatusOK))

		agg := resp.Body.(*goserver.Aggregation)

		// Expect 1 interval (the year 2023)
		// Intervals logic depends on utils.RoundToGranularity inside GetExpenses if from/to are zero,
		// OR inside Aggregate calls getIntervals.
		// getIntervals(from, to, year) -> [2023-01-01]
		Expect(agg.Intervals).To(HaveLen(1))
		Expect(agg.Intervals[0].Year()).To(Equal(2023))

		Expect(agg.Currencies).To(HaveLen(1))
		cur := agg.Currencies[0]
		Expect(cur.CurrencyId).To(Equal("USD"))

		Expect(cur.Accounts).To(HaveLen(1))
		acc := cur.Accounts[0]
		Expect(acc.AccountId).To(Equal("acc-a"))
		Expect(acc.Amounts).To(HaveLen(1))

		// Total: 100 + 200 + 300 = 600
		Expect(acc.Amounts[0]).To(Equal(600.0))
	})
})

package api_test

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("BudgetItems API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), common.UserIDKey, "user1")

	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         goserver.BudgetItemsAPIServicer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewBudgetItemsAPIService(log, mockStorage)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("calculates budget status with rollover correctly", func() {
		// Setup dates
		startOfMonth := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		// endOfMonth := startOfMonth.AddDate(0, 1, 0)
		budgetItems := []goserver.BudgetItem{
			{
				Date:      startOfMonth,
				AccountId: "account-a",
				Amount:    decimal.NewFromInt(100),
			},
			{
				Date:      startOfMonth.AddDate(0, 1, 0), // Feb
				AccountId: "account-a",
				Amount:    decimal.NewFromInt(100),
			},
		}

		// Setup Transactions (Expenses)
		// User spent 50 for Account A in Jan
		// User spent 120 for Account A in Feb
		transactions := []goserver.Transaction{
			{
				Date: startOfMonth.Add(time.Hour),
				Movements: []goserver.Movement{
					{AccountId: "account-a", Amount: decimal.NewFromInt(50)},
				},
			},
			{
				Date: startOfMonth.AddDate(0, 1, 0).Add(time.Hour),
				Movements: []goserver.Movement{
					{AccountId: "account-a", Amount: decimal.NewFromInt(120)},
				},
			},
		}

		// Mock Calls
		mockStorage.EXPECT().GetBudgetItems("user1").Return(budgetItems, nil)
		mockStorage.EXPECT().GetAccounts("user1").Return([]goserver.Account{
			{Id: "account-a", HideFromReports: false},
		}, nil)
		mockStorage.EXPECT().GetCurrencies("user1").Return([]goserver.Currency{}, nil)
		// It will fetch transactions from MinDate (Jan 1) to requested To date.
		mockStorage.EXPECT().GetTransactions("user1", gomock.Any(), gomock.Any(), false).Return(transactions, nil)

		// Call SUT for Jan and Feb status
		resp, err := sut.GetBudgetStatus(ctx, startOfMonth, startOfMonth.AddDate(0, 2, 0), "", false)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Code).To(Equal(http.StatusOK))

		body := resp.Body.([]goserver.BudgetStatus)
		Expect(body).To(HaveLen(2))

		// Check Jan
		// Budget: 100, Spent: 50, Rollover: 0, Available: 50
		jan := body[0]
		Expect(jan.Budgeted.Equal(decimal.NewFromInt(100))).To(BeTrue())
		Expect(jan.Spent.Equal(decimal.NewFromInt(50))).To(BeTrue())
		Expect(jan.Rollover.Equal(decimal.Zero)).To(BeTrue())
		Expect(jan.Available.Equal(decimal.NewFromInt(50))).To(BeTrue()) // Remainder for Jan

		// Check Feb
		// Budget: 100, Spent: 120
		// Rollover from Jan: 50 (Surplus)
		// Available (Total): 100 + 50 = 150
		// Remainder: 150 - 120 = 30
		feb := body[1]
		Expect(feb.Budgeted.Equal(decimal.NewFromInt(100))).To(BeTrue())
		Expect(feb.Spent.Equal(decimal.NewFromInt(120))).To(BeTrue())
		Expect(feb.Rollover.Equal(decimal.NewFromInt(50))).To(BeTrue()) // Surplus from Jan
		Expect(feb.Available.Equal(decimal.NewFromInt(30))).To(BeTrue())
	})

	It("handles negative rollover (deficit)", func() {
		startOfMonth := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		budgetItems := []goserver.BudgetItem{
			{Date: startOfMonth, AccountId: "acc-b", Amount: decimal.NewFromInt(100)},
			{Date: startOfMonth.AddDate(0, 1, 0), AccountId: "acc-b", Amount: decimal.NewFromInt(100)},
		}

		// Jan: Spent 150 (Deficit 50)
		transactions := []goserver.Transaction{
			{
				Date:      startOfMonth.Add(time.Hour),
				Movements: []goserver.Movement{{AccountId: "acc-b", Amount: decimal.NewFromInt(150)}},
			},
		}

		mockStorage.EXPECT().GetBudgetItems("user1").Return(budgetItems, nil)
		mockStorage.EXPECT().GetAccounts("user1").Return([]goserver.Account{
			{Id: "acc-b", HideFromReports: false},
		}, nil)
		mockStorage.EXPECT().GetCurrencies("user1").Return([]goserver.Currency{}, nil)
		mockStorage.EXPECT().GetTransactions("user1", gomock.Any(), gomock.Any(), false).Return(transactions, nil)

		resp, err := sut.GetBudgetStatus(ctx, startOfMonth, startOfMonth.AddDate(0, 2, 0), "", false)
		Expect(err).ToNot(HaveOccurred())
		body := resp.Body.([]goserver.BudgetStatus)

		// Jan: Remainder -50
		jan := body[0]
		Expect(jan.Available.Equal(decimal.NewFromInt(-50))).To(BeTrue())
		// ...
		feb := body[1]
		Expect(feb.Rollover.Equal(decimal.NewFromInt(-50))).To(BeTrue())
		Expect(feb.Available.Equal(decimal.NewFromInt(50))).To(BeTrue())
	})
})

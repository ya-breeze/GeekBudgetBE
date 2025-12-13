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
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("BudgetItems API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), "userID", "user1")

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
				Amount:    100.0,
			},
			{
				Date:      startOfMonth.AddDate(0, 1, 0), // Feb
				AccountId: "account-a",
				Amount:    100.0,
			},
		}

		// Setup Transactions (Expenses)
		// User spent 50 for Account A in Jan
		// User spent 120 for Account A in Feb
		transactions := []goserver.Transaction{
			{
				Date: startOfMonth.Add(time.Hour),
				Movements: []goserver.Movement{
					{AccountId: "account-a", Amount: 50.0},
				},
			},
			{
				Date: startOfMonth.AddDate(0, 1, 0).Add(time.Hour),
				Movements: []goserver.Movement{
					{AccountId: "account-a", Amount: 120.0},
				},
			},
		}

		// Mock Calls
		mockStorage.EXPECT().GetBudgetItems("user1").Return(budgetItems, nil)
		// It will fetch transactions from MinDate (Jan 1) to requested To date.
		mockStorage.EXPECT().GetTransactions("user1", gomock.Any(), gomock.Any()).Return(transactions, nil)

		// Call SUT for Jan and Feb status
		resp, err := sut.GetBudgetStatus(ctx, startOfMonth, startOfMonth.AddDate(0, 2, 0))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Code).To(Equal(http.StatusOK))

		body := resp.Body.([]goserver.BudgetStatus)
		Expect(body).To(HaveLen(2))

		// Check Jan
		// Budget: 100, Spent: 50, Rollover: 0, Available: 50
		jan := body[0]
		Expect(jan.Budgeted).To(Equal(100.0))
		Expect(jan.Spent).To(Equal(50.0))
		Expect(jan.Rollover).To(Equal(0.0))
		Expect(jan.Available).To(Equal(50.0)) // Remainder for Jan

		// Check Feb
		// Budget: 100, Spent: 120
		// Rollover from Jan: 50 (Surplus)
		// Available (Total): 100 + 50 = 150
		// Remainder: 150 - 120 = 30
		feb := body[1]
		Expect(feb.Budgeted).To(Equal(100.0))
		Expect(feb.Spent).To(Equal(120.0))
		Expect(feb.Rollover).To(Equal(50.0)) // Surplus from Jan
		Expect(feb.Available).To(Equal(30.0))
	})

	It("handles negative rollover (deficit)", func() {
		startOfMonth := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		budgetItems := []goserver.BudgetItem{
			{Date: startOfMonth, AccountId: "acc-b", Amount: 100.0},
			{Date: startOfMonth.AddDate(0, 1, 0), AccountId: "acc-b", Amount: 100.0},
		}

		// Jan: Spent 150 (Deficit 50)
		transactions := []goserver.Transaction{
			{
				Date:      startOfMonth.Add(time.Hour),
				Movements: []goserver.Movement{{AccountId: "acc-b", Amount: 150.0}},
			},
		}

		mockStorage.EXPECT().GetBudgetItems("user1").Return(budgetItems, nil)
		mockStorage.EXPECT().GetTransactions("user1", gomock.Any(), gomock.Any()).Return(transactions, nil)

		resp, err := sut.GetBudgetStatus(ctx, startOfMonth, startOfMonth.AddDate(0, 2, 0))
		Expect(err).ToNot(HaveOccurred())
		body := resp.Body.([]goserver.BudgetStatus)

		// Jan: Remainder -50
		jan := body[0]
		Expect(jan.Available).To(Equal(-50.0))

		// Feb:
		// Budget 100
		// Rollover -50
		// Available: 100 - 50 = 50
		// Spent: 0
		// Remainder: 50
		feb := body[1]
		Expect(feb.Rollover).To(Equal(-50.0))
		Expect(feb.Available).To(Equal(50.0))
	})
})

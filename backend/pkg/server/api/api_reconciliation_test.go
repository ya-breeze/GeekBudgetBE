package api_test

import (
	"context"
	"net/http"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Reconciliation API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), constants.UserIDKey, "user1")

	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         *api.ReconciliationAPIServiceImpl
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewReconciliationAPIServiceImpl(log, mockStorage)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("GetReconciliationStatus", func() {
		It("includes accounts marked with ShowInReconciliation even without importer", func() {
			// Setup mock data
			accounts := []goserver.Account{
				{
					Id:                   "acc1",
					Name:                 "Manual Account",
					Type:                 "asset",
					ShowInReconciliation: true,
					BankInfo: goserver.BankAccountInfo{
						Balances: []goserver.BankAccountInfoBalancesInner{
							{CurrencyId: "curr1", ClosingBalance: decimal.NewFromInt(100)},
						},
					},
				},
				{
					Id:                   "acc2",
					Name:                 "Hidden Account",
					Type:                 "asset",
					ShowInReconciliation: false,
					BankInfo: goserver.BankAccountInfo{
						Balances: []goserver.BankAccountInfoBalancesInner{
							{CurrencyId: "curr1", ClosingBalance: decimal.NewFromInt(100)},
						},
					},
				},
			}

			mockStorage.EXPECT().GetAccounts("user1").Return(accounts, nil)
			mockStorage.EXPECT().GetBankImporters("user1").Return([]goserver.BankImporter{}, nil)
			mockStorage.EXPECT().GetCurrencies("user1").Return([]goserver.Currency{{Id: "curr1", Name: "USD"}}, nil)

			bulkData := &database.BulkReconciliationData{
				Balances: map[string]map[string]decimal.Decimal{
					"acc1": {"curr1": decimal.NewFromInt(100)},
					"acc2": {"curr1": decimal.NewFromInt(100)},
				},
				LatestReconciliations: make(map[string]map[string]*goserver.Reconciliation),
				UnprocessedCounts:     make(map[string]int),
				MaxTransactionDates:   make(map[string]map[string]time.Time),
			}
			mockStorage.EXPECT().GetBulkReconciliationData("user1").Return(bulkData, nil)

			// acc2 should be skipped because it has no importer, no manual reconciliation, and ShowInReconciliation is false

			resp, err := sut.GetReconciliationStatus(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))

			statuses := resp.Body.([]goserver.ReconciliationStatus)
			Expect(statuses).To(HaveLen(1))
			Expect(statuses[0].AccountId).To(Equal("acc1"))
		})

		It("returns complex status with importers, manual reconciliations and unprocessed transactions", func() {
			t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
			tAfter := t1.Add(time.Hour)

			accounts := []goserver.Account{
				{
					Id:   "acc_importer",
					Name: "Importer Account",
					Type: "asset",
					BankInfo: goserver.BankAccountInfo{
						Balances: []goserver.BankAccountInfoBalancesInner{
							{CurrencyId: "USD", ClosingBalance: decimal.NewFromInt(1000), LastUpdatedAt: &t1},
						},
					},
				},
				{
					Id:   "acc_manual",
					Name: "Manual Account",
					Type: "asset",
					BankInfo: goserver.BankAccountInfo{
						Balances: []goserver.BankAccountInfoBalancesInner{
							{CurrencyId: "USD", ClosingBalance: decimal.NewFromInt(500)},
						},
					},
				},
			}

			importers := []goserver.BankImporter{
				{AccountId: "acc_importer"},
			}

			currencies := []goserver.Currency{
				{Id: "USD", Name: "$"},
			}

			lastRecManual := &goserver.Reconciliation{
				AccountId:         "acc_manual",
				CurrencyId:        "USD",
				ReconciledBalance: decimal.NewFromInt(450),
				ReconciledAt:      t1,
				IsManual:          true,
			}

			mockStorage.EXPECT().GetAccounts("user1").Return(accounts, nil)
			mockStorage.EXPECT().GetBankImporters("user1").Return(importers, nil)
			mockStorage.EXPECT().GetCurrencies("user1").Return(currencies, nil)

			bulkData := &database.BulkReconciliationData{
				Balances: map[string]map[string]decimal.Decimal{
					"acc_importer": {"USD": decimal.NewFromInt(1050)},
					"acc_manual":   {"USD": decimal.NewFromInt(450)},
				},
				LatestReconciliations: map[string]map[string]*goserver.Reconciliation{
					"acc_manual": {"USD": lastRecManual},
				},
				UnprocessedCounts: map[string]int{
					"acc_importer": 5,
					"acc_manual":   0,
				},
				MaxTransactionDates: map[string]map[string]time.Time{
					"acc_importer": {"USD": tAfter},
					"acc_manual":   {"USD": t1},
				},
			}
			mockStorage.EXPECT().GetBulkReconciliationData("user1").Return(bulkData, nil)

			resp, err := sut.GetReconciliationStatus(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))

			statuses := resp.Body.([]goserver.ReconciliationStatus)
			Expect(statuses).To(HaveLen(2))

			// Check importer status
			var sImp goserver.ReconciliationStatus
			for _, s := range statuses {
				if s.AccountId == "acc_importer" {
					sImp = s
				}
			}
			Expect(sImp.AccountName).To(Equal("Importer Account"))
			Expect(sImp.BankBalance).To(Equal(decimal.NewFromInt(1000)))
			Expect(sImp.AppBalance).To(Equal(decimal.NewFromInt(1050)))
			Expect(sImp.Delta).To(Equal(decimal.NewFromInt(50)))
			Expect(sImp.HasUnprocessedTransactions).To(BeTrue())
			Expect(sImp.HasBankImporter).To(BeTrue())
			Expect(sImp.HasTransactionsAfterBankBalance).To(BeTrue())

			// Check manual status
			var sMan goserver.ReconciliationStatus
			for _, s := range statuses {
				if s.AccountId == "acc_manual" {
					sMan = s
				}
			}
			Expect(sMan.AccountName).To(Equal("Manual Account"))
			Expect(sMan.BankBalance).To(Equal(decimal.NewFromInt(450))) // Uses reconciled balance for manual
			Expect(sMan.AppBalance).To(Equal(decimal.NewFromInt(450)))
			Expect(sMan.Delta.IsZero()).To(BeTrue())
			Expect(sMan.HasUnprocessedTransactions).To(BeFalse())
			Expect(sMan.HasBankImporter).To(BeFalse())
			Expect(sMan.HasTransactionsAfterBankBalance).To(BeFalse())
			Expect(sMan.IsManualReconciliationEnabled).To(BeTrue())
		})
	})
})

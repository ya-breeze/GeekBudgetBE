package api_test

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	ctx := context.WithValue(context.Background(), constants.FamilyIDKey, uuid.MustParse("00000000-0000-0000-0000-000000000001"))

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

			mockStorage.EXPECT().GetAccounts(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(accounts, nil)
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{}, nil)
			mockStorage.EXPECT().GetCurrencies(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.Currency{{Id: "curr1", Name: "USD"}}, nil)

			bulkData := &database.BulkReconciliationData{
				Balances: map[string]map[string]decimal.Decimal{
					"acc1": {"curr1": decimal.NewFromInt(100)},
					"acc2": {"curr1": decimal.NewFromInt(100)},
				},
				LatestReconciliations: make(map[string]map[string]*goserver.Reconciliation),
				UnprocessedCounts:     make(map[string]int),
				MaxTransactionDates:   make(map[string]map[string]time.Time),
			}
			mockStorage.EXPECT().GetBulkReconciliationData(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(bulkData, nil)

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

			mockStorage.EXPECT().GetAccounts(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(accounts, nil)
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(importers, nil)
			mockStorage.EXPECT().GetCurrencies(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(currencies, nil)

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
			mockStorage.EXPECT().GetBulkReconciliationData(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return(bulkData, nil)

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

	Describe("ReconcileAccount", func() {
		var importerAccount goserver.Account
		var importerEntry goserver.BankImporter

		BeforeEach(func() {
			importerAccount = goserver.Account{
				Id:   "acc_importer",
				Name: "Bank",
				Type: "asset",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{CurrencyId: "USD", ClosingBalance: decimal.NewFromInt(1000)},
					},
				},
			}
			importerEntry = goserver.BankImporter{AccountId: "acc_importer"}
		})

		It("reconciles no-importer account with large delta — returns 200 with IsManual=true and delta=0 in history", func() {
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{}, nil)
			mockStorage.EXPECT().GetAccountBalance(uuid.MustParse("00000000-0000-0000-0000-000000000001"), "acc_noimporter", "USD").
				Return(decimal.NewFromInt(500), nil)

			mockStorage.EXPECT().CreateReconciliation(uuid.MustParse("00000000-0000-0000-0000-000000000001"), gomock.Any()).
				DoAndReturn(func(_ uuid.UUID, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
					Expect(rec.ReconciledBalance).To(Equal(decimal.NewFromInt(500)))
					Expect(rec.ExpectedBalance).To(Equal(decimal.NewFromInt(500))) // delta=0 in history
					Expect(rec.IsManual).To(BeTrue())
					return goserver.Reconciliation{}, nil
				})

			resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
				CurrencyId: "USD",
				Balance:    decimal.NewFromInt(0), // frontend always sends 0 — "use app balance"
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("reconciles no-importer account with balance=0 in request — IsManual=true even if fetched balance is zero", func() {
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{}, nil)
			mockStorage.EXPECT().GetAccountBalance(uuid.MustParse("00000000-0000-0000-0000-000000000001"), "acc_noimporter", "USD").
				Return(decimal.NewFromInt(0), nil) // zero balance

			mockStorage.EXPECT().CreateReconciliation(uuid.MustParse("00000000-0000-0000-0000-000000000001"), gomock.Any()).
				DoAndReturn(func(_ uuid.UUID, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
					Expect(rec.IsManual).To(BeTrue()) // must be true even when balance is zero
					Expect(rec.ReconciledBalance.IsZero()).To(BeTrue())
					Expect(rec.ExpectedBalance.IsZero()).To(BeTrue())
					return goserver.Reconciliation{}, nil
				})

			resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
				CurrencyId: "USD",
				Balance:    decimal.NewFromInt(0),
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("reconciles no-importer account within tolerance — returns 200 (happy path unchanged)", func() {
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{}, nil)

			mockStorage.EXPECT().CreateReconciliation(uuid.MustParse("00000000-0000-0000-0000-000000000001"), gomock.Any()).
				DoAndReturn(func(_ uuid.UUID, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
					Expect(rec.IsManual).To(BeTrue())
					Expect(rec.ReconciledBalance).To(Equal(decimal.NewFromFloat(500.005)))
					Expect(rec.ExpectedBalance).To(Equal(decimal.NewFromFloat(500.005))) // no-importer: expectedBalance == balance
					return goserver.Reconciliation{}, nil
				})

			resp, err := sut.ReconcileAccount(ctx, "acc_noimporter", goserver.ReconcileAccountRequest{
				CurrencyId: "USD",
				Balance:    decimal.NewFromFloat(500.005), // within tolerance (0.01)
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("blocks importer account with large delta — returns 400 (existing behavior preserved)", func() {
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{importerEntry}, nil)
			mockStorage.EXPECT().GetAccount(uuid.MustParse("00000000-0000-0000-0000-000000000001"), "acc_importer").Return(importerAccount, nil)
			// Bank balance is 1000, request balance is 500 → delta = 500 → 400
			resp, err := sut.ReconcileAccount(ctx, "acc_importer", goserver.ReconcileAccountRequest{
				CurrencyId: "USD",
				Balance:    decimal.NewFromInt(500),
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})

		It("blocks account with importer record but no BankInfo.Balances (never run) — treated as has-importer, returns 400", func() {
			accountNoBankInfo := goserver.Account{
				Id:   "acc_importer",
				Name: "Bank",
				Type: "asset",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{}, // no balance data from importer
				},
			}
			mockStorage.EXPECT().GetBankImporters(uuid.MustParse("00000000-0000-0000-0000-000000000001")).Return([]goserver.BankImporter{importerEntry}, nil)
			mockStorage.EXPECT().GetAccount(uuid.MustParse("00000000-0000-0000-0000-000000000001"), "acc_importer").Return(accountNoBankInfo, nil)
			// expectedBalance = 0 (no balance entries), balance = 500 → delta = 500 → 400
			resp, err := sut.ReconcileAccount(ctx, "acc_importer", goserver.ReconcileAccountRequest{
				CurrencyId: "USD",
				Balance:    decimal.NewFromInt(500),
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})
	})
})

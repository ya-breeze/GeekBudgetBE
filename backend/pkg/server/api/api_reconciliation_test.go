package api_test

import (
	"context"
	"net/http"

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

var _ = Describe("Reconciliation API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), common.UserIDKey, "user1")

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

			// acc1 expectations
			mockStorage.EXPECT().GetLatestReconciliation("user1", "acc1", "curr1").Return(nil, nil)
			mockStorage.EXPECT().GetAccountBalance("user1", "acc1", "curr1").Return(decimal.NewFromFloat(100.0), nil)
			mockStorage.EXPECT().CountUnprocessedTransactionsForAccount("user1", "acc1", gomock.Any()).Return(0, nil)

			// acc2 expectations
			mockStorage.EXPECT().GetLatestReconciliation("user1", "acc2", "curr1").Return(nil, nil)

			// acc2 should be skipped because it has no importer, no manual reconciliation, and ShowInReconciliation is false

			resp, err := sut.GetReconciliationStatus(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))

			statuses := resp.Body.([]goserver.ReconciliationStatus)
			Expect(statuses).To(HaveLen(1))
			Expect(statuses[0].AccountId).To(Equal("acc1"))
		})
	})
})

package common

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/test"
)

func TestBalanceChecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BalanceChecker Suite")
}

var _ = Describe("BalanceChecker", func() {
	var (
		mockCtrl *gomock.Controller
		mockDB   *mocks.MockStorage
		logger   = test.CreateTestLogger()
		userID   = "test-user"
		accID    = "test-acc"
		ctx      = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("CheckBalanceForAccount", func() {
		It("should return early if there are unprocessed transactions", func() {
			acc := goserver.Account{
				Id:   accID,
				Name: "Test Account",
			}
			mockDB.EXPECT().GetAccount(userID, accID).Return(acc, nil)
			mockDB.EXPECT().CountUnprocessedTransactionsForAccount(userID, accID, gomock.Any()).Return(5, nil)

			err := CheckBalanceForAccount(ctx, logger, mockDB, userID, accID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should generate a notification on balance mismatch", func() {
			acc := goserver.Account{
				Id:   accID,
				Name: "Test Account",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{
							CurrencyId:     "CZK",
							ClosingBalance: 1500.0,
						},
					},
				},
			}
			mockDB.EXPECT().GetAccount(userID, accID).Return(acc, nil)
			mockDB.EXPECT().CountUnprocessedTransactionsForAccount(userID, accID, gomock.Any()).Return(0, nil)

			// App balance is 1400, Bank says 1500 -> Mismatch
			mockDB.EXPECT().GetAccountBalance(userID, accID, "CZK").Return(1400.0, nil)

			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).DoAndReturn(func(uid string, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Type).To(Equal(string(models.NotificationTypeBalanceDoesntMatch)))
				Expect(n.Title).To(Equal("Balance Mismatch Detected"))
				Expect(n.Description).To(ContainSubstring("App balance: 1400.00"))
				Expect(n.Description).To(ContainSubstring("Bank balance: 1500.00"))
				return goserver.Notification{}, nil
			})

			err := CheckBalanceForAccount(ctx, logger, mockDB, userID, accID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should not generate a notification if balances match", func() {
			acc := goserver.Account{
				Id:   accID,
				Name: "Test Account",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{
							CurrencyId:     "CZK",
							ClosingBalance: 1000.0,
						},
					},
				},
			}
			mockDB.EXPECT().GetAccount(userID, accID).Return(acc, nil)
			mockDB.EXPECT().CountUnprocessedTransactionsForAccount(userID, accID, gomock.Any()).Return(0, nil)

			mockDB.EXPECT().GetAccountBalance(userID, accID, "CZK").Return(1000.0, nil)

			mockDB.EXPECT().CreateReconciliation(userID, gomock.Any()).Return(goserver.Reconciliation{}, nil)

			// No CreateNotification expected

			err := CheckBalanceForAccount(ctx, logger, mockDB, userID, accID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should check multiple currencies", func() {
			acc := goserver.Account{
				Id:   accID,
				Name: "Multi Currency",
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{CurrencyId: "CZK", ClosingBalance: 100.0},
						{CurrencyId: "USD", ClosingBalance: 200.0},
					},
				},
			}
			mockDB.EXPECT().GetAccount(userID, accID).Return(acc, nil)
			mockDB.EXPECT().CountUnprocessedTransactionsForAccount(userID, accID, gomock.Any()).Return(0, nil)

			mockDB.EXPECT().GetAccountBalance(userID, accID, "CZK").Return(100.0, nil)
			mockDB.EXPECT().CreateReconciliation(userID, gomock.Any()).Return(goserver.Reconciliation{}, nil)
			mockDB.EXPECT().GetAccountBalance(userID, accID, "USD").Return(250.0, nil) // USD mismatch!

			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).DoAndReturn(func(uid string, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Description).To(ContainSubstring("Currency: USD"))
				return goserver.Notification{}, nil
			})

			err := CheckBalanceForAccount(ctx, logger, mockDB, userID, accID)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle database errors gracefully", func() {
			mockDB.EXPECT().GetAccount(userID, accID).Return(goserver.Account{}, fmt.Errorf("db error"))

			err := CheckBalanceForAccount(ctx, logger, mockDB, userID, accID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to get account"))
		})
	})
})

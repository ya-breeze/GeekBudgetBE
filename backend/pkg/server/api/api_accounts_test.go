package api_test

import (
	"context"
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Accounts API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), common.UserIDKey, "user1")
	cfg := &config.Config{}

	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         goserver.AccountsAPIServicer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewAccountsAPIService(log, mockStorage, cfg)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("DeleteAccount", func() {
		It("returns 400 when deleting an account in use without replacement", func() {
			accountID := "acc-1-in-use"

			// Mock GetAccount (image check)
			mockStorage.EXPECT().GetAccount("user1", accountID).Return(goserver.Account{Id: accountID}, nil)

			// Mock DeleteAccount returning ErrAccountInUse
			// Note: empty string for replaceWithAccountId means no replacement
			emptyReplace := ""
			mockStorage.EXPECT().DeleteAccount("user1", accountID, &emptyReplace).Return(database.ErrAccountInUse)

			resp, err := sut.DeleteAccount(ctx, accountID, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})

		It("returns 400 when replacing with self", func() {
			accountID := "acc-1"

			// Mock GetAccount (image check)
			mockStorage.EXPECT().GetAccount("user1", accountID).Return(goserver.Account{Id: accountID}, nil)

			// We expect NO calls to GetAccount for validation of replacement because the check happens before
			// We expect NO calls to DeleteAccount

			resp, err := sut.DeleteAccount(ctx, accountID, accountID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})

		It("returns 200 when deleting an account with valid replacement", func() {
			accountID := "acc-1"
			replaceID := "acc-2"

			// Mock GetAccount (image check)
			mockStorage.EXPECT().GetAccount("user1", accountID).Return(goserver.Account{Id: accountID}, nil)

			// Mock validation of replacement account
			mockStorage.EXPECT().GetAccount("user1", replaceID).Return(goserver.Account{Id: replaceID}, nil)

			// Mock DeleteAccount success
			mockStorage.EXPECT().DeleteAccount("user1", accountID, &replaceID).Return(nil)

			resp, err := sut.DeleteAccount(ctx, accountID, replaceID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})
	})
})

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

var _ = Describe("Transactions API", func() {
	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         goserver.TransactionsAPIServicer
		ctx         context.Context
		log         = test.CreateTestLogger()
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		sut = api.NewTransactionsAPIService(log, mockStorage)
		ctx = context.WithValue(context.Background(), common.UserIDKey, "user-1")
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("handles onlySuspicious parameter", func() {
		userID := "user-1"
		dateFrom := time.Time{}
		dateTo := time.Time{}

		mockStorage.EXPECT().
			GetTransactions(userID, dateFrom, dateTo, true).
			Return([]goserver.Transaction{}, nil)

		resp, err := sut.GetTransactions(ctx, "", 0, 0, dateFrom, dateTo, true)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Code).To(Equal(http.StatusOK))
	})
})

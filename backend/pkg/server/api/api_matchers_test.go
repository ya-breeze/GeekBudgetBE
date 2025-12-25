package api_test

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Matchers API", func() {
	log := test.CreateTestLogger()
	ctx := context.WithValue(context.Background(), common.UserIDKey, "user1")
	cfg := &config.Config{}

	var (
		ctrl        *gomock.Controller
		mockStorage *mocks.MockStorage
		sut         goserver.MatchersAPIServicer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockStorage = mocks.NewMockStorage(ctrl)
		// We pass nil for unprocessedService as it is not used in DeleteMatcher
		sut = api.NewMatchersAPIServiceImpl(log, mockStorage, cfg, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("DeleteMatcher", func() {
		It("returns 204 when deleting a matcher successfully", func() {
			matcherID := "matcher-1"

			mockStorage.EXPECT().DeleteMatcher("user1", matcherID).Return(nil)

			resp, err := sut.DeleteMatcher(ctx, matcherID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(204))
		})

		It("returns 500 when storage fails", func() {
			matcherID := "matcher-failed"

			mockStorage.EXPECT().DeleteMatcher("user1", matcherID).Return(errors.New("db error"))

			resp, err := sut.DeleteMatcher(ctx, matcherID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})

		It("returns 500 when user ID is missing from context", func() {
			resp, err := sut.DeleteMatcher(context.Background(), "some-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})
})

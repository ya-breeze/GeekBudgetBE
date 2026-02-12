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

	Describe("UpdateMatcher", func() {
		It("updates a matcher successfully", func() {
			matcherID := "matcher-update"
			userID := "user1"
			input := goserver.MatcherNoId{
				DescriptionRegExp: "^Test$",
				OutputDescription: "Matched",
				OutputAccountId:   "acc1",
			}
			updatedMatcher := goserver.Matcher{
				Id:                matcherID,
				DescriptionRegExp: input.DescriptionRegExp,
				OutputDescription: input.OutputDescription,
				OutputAccountId:   input.OutputAccountId,
			}

			// Instantiate real UnprocessedTransactionsAPIServiceImpl with mockStorage
			unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(log, mockStorage)
			// Re-instantiate SUT with the service
			sut = api.NewMatchersAPIServiceImpl(log, mockStorage, cfg, unprocessedService)

			mockStorage.EXPECT().UpdateMatcher(userID, matcherID, &input).Return(updatedMatcher, nil)

			// ProcessUnprocessedTransactionsAgainstMatcher calls:
			mockStorage.EXPECT().GetMatcher(userID, matcherID).Return(updatedMatcher, nil)
			// Assuming no confirmation history, it returns early.
			// Let's check logic: if len(matcher.ConfirmationHistory) < 10 { return nil, nil }
			// So we MUST return a matcher with enough history OR expect it to return early.
			// Default updatedMatcher has no history, so it should return early.
			// BUT GetMatcher is called first.

			resp, err := sut.UpdateMatcher(ctx, matcherID, input)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
			body := resp.Body.(goserver.UpdateMatcher200Response)
			Expect(body.Matcher.Id).To(Equal(updatedMatcher.Id))
			Expect(body.AutoProcessedIds).To(BeEmpty())
		})
	})
})

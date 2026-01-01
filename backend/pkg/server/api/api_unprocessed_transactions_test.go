package api

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("UnprocessedTransactions API", func() {
	var (
		mockCtrl *gomock.Controller
		mockDB   *mocks.MockStorage
		sut      *UnprocessedTransactionsAPIServiceImpl
		logger   = test.CreateTestLogger()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
		sut = NewUnprocessedTransactionsAPIServiceImpl(logger, mockDB)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("getDuplicateTransactions", func() {
		It("should not match transactions with different currencies even if amounts are similar", func() {
			t1Date := time.Date(2025, 2, 28, 1, 0, 0, 0, time.UTC)
			t1 := goserver.Transaction{
				Id:   "tx1",
				Date: t1Date,
				Movements: []goserver.Movement{
					{Amount: 7.00, CurrencyId: "CZK", AccountId: ""}, // Unprocessed
				},
			}

			t2Date := time.Date(2025, 2, 28, 2, 14, 0, 0, time.UTC)
			t2 := goserver.Transaction{
				Id:   "tx2",
				Date: t2Date,
				Movements: []goserver.Movement{
					{Amount: 7.50, CurrencyId: "EUR", AccountId: "acc1"}, // Processed candidate
				},
			}

			candidates := []goserver.Transaction{t2}
			duplicates := sut.getDuplicateTransactions(candidates, t1)

			Expect(duplicates).To(BeEmpty(), "7.00 CZK should not match 7.50 EUR")
		})

		It("should match transactions with same currency and amounts within 1.0 threshold", func() {
			t1Date := time.Date(2025, 2, 28, 1, 0, 0, 0, time.UTC)
			t1 := goserver.Transaction{
				Id:   "tx1",
				Date: t1Date,
				Movements: []goserver.Movement{
					{Amount: 7.00, CurrencyId: "CZK", AccountId: ""},
				},
			}

			t2Date := time.Date(2025, 2, 28, 2, 14, 0, 0, time.UTC)
			t2 := goserver.Transaction{
				Id:   "tx2",
				Date: t2Date,
				Movements: []goserver.Movement{
					{Amount: 7.50, CurrencyId: "CZK", AccountId: "acc1"},
				},
			}

			candidates := []goserver.Transaction{t2}
			duplicates := sut.getDuplicateTransactions(candidates, t1)

			Expect(duplicates).To(HaveLen(1))
			Expect(duplicates[0].Id).To(Equal("tx2"))
		})

		It("should not match if currencies count differs", func() {
			t1 := goserver.Transaction{
				Id:   "tx1",
				Date: time.Now(),
				Movements: []goserver.Movement{
					{Amount: 100, CurrencyId: "CZK"},
					{Amount: 50, CurrencyId: "EUR"},
				},
			}

			t2 := goserver.Transaction{
				Id:   "tx2",
				Date: time.Now(),
				Movements: []goserver.Movement{
					{Amount: 100, CurrencyId: "CZK", AccountId: "acc1"},
				},
			}

			candidates := []goserver.Transaction{t2}
			duplicates := sut.getDuplicateTransactions(candidates, t1)

			Expect(duplicates).To(BeEmpty())
		})
	})
})

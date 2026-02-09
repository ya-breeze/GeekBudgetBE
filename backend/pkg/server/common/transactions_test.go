package common_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

var _ = Describe("Transaction Utilities", func() {
	Describe("GetIncreases", func() {
		It("should correctly calculate increases for single currency", func() {
			movements := []goserver.Movement{
				{Amount: 100, CurrencyId: "USD"},
				{Amount: -50, CurrencyId: "USD"},
				{Amount: 200, CurrencyId: "USD"},
			}
			increases := common.GetIncreases(movements)
			Expect(increases).To(HaveLen(1))
			Expect(increases["USD"]).To(Equal(float64(300)))
		})

		It("should correctly calculate increases for multiple currencies", func() {
			movements := []goserver.Movement{
				{Amount: 100, CurrencyId: "USD"},
				{Amount: 50, CurrencyId: "EUR"},
				{Amount: -25, CurrencyId: "USD"},
				{Amount: 150, CurrencyId: "EUR"},
			}
			increases := common.GetIncreases(movements)
			Expect(increases).To(HaveLen(2))
			Expect(increases["USD"]).To(Equal(float64(100)))
			Expect(increases["EUR"]).To(Equal(float64(200)))
		})

		It("should handle negative net increase by taking absolute of negative sum", func() {
			// Logic: posSum vs negSum(absolute). Takes the larger one.
			movements := []goserver.Movement{
				{Amount: 10, CurrencyId: "USD"},
				{Amount: -100, CurrencyId: "USD"},
			}
			// pos=10, neg=100. max(10, 100) = 100
			increases := common.GetIncreases(movements)
			Expect(increases["USD"]).To(Equal(float64(100)))
		})
	})

	Describe("IsDuplicate", func() {
		var (
			now = time.Now()
			m1  = []goserver.Movement{{Amount: 100, CurrencyId: "USD"}}
			m2  = []goserver.Movement{{Amount: 100, CurrencyId: "USD"}}
		)

		It("should match transactions with same date and movements", func() {
			Expect(common.IsDuplicate(now, m1, now, m2)).To(BeTrue())
		})

		It("should match transactions within 2 days", func() {
			twoDaysLater := now.Add(47 * time.Hour)
			Expect(common.IsDuplicate(now, m1, twoDaysLater, m2)).To(BeTrue())
		})

		It("should not match transactions more than 2 days apart", func() {
			threeDaysLater := now.Add(73 * time.Hour)
			Expect(common.IsDuplicate(now, m1, threeDaysLater, m2)).To(BeFalse())
		})

		It("should not match transactions with different amounts", func() {
			mDiff := []goserver.Movement{{Amount: 101, CurrencyId: "USD"}}
			Expect(common.IsDuplicate(now, m1, now, mDiff)).To(BeFalse())
		})

		It("should not match transactions with different currencies", func() {
			mDiff := []goserver.Movement{{Amount: 100, CurrencyId: "EUR"}}
			Expect(common.IsDuplicate(now, m1, now, mDiff)).To(BeFalse())
		})

		It("should match multi-movement transactions with same total increases", func() {
			ma := []goserver.Movement{
				{Amount: 100, CurrencyId: "USD"},
				{Amount: -100, CurrencyId: "USD"}, // should be ignored for increase calculation
				{Amount: 50, CurrencyId: "EUR"},
			}
			mb := []goserver.Movement{
				{Amount: 100, CurrencyId: "USD"},
				{Amount: 50, CurrencyId: "EUR"},
			}
			Expect(common.IsDuplicate(now, ma, now, mb)).To(BeTrue())
		})
	})
})

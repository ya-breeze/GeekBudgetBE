package utils

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("Transactions Utils", func() {
	Describe("GetIncreases", func() {
		It("should calculate net increases per currency", func() {
			movements := []goserver.Movement{
				{Amount: decimal.NewFromFloat(100.50), CurrencyId: "CZK", AccountId: "acc1"},
				{Amount: decimal.NewFromFloat(-50.25), CurrencyId: "CZK", AccountId: "acc2"},
				{Amount: decimal.NewFromFloat(200.00), CurrencyId: "USD", AccountId: "acc1"},
			}
			increases := GetIncreases(movements)
			Expect(increases).To(HaveLen(2))
			Expect(increases["CZK"].Equal(decimal.NewFromFloat(100.50))).To(BeTrue())
			Expect(increases["USD"].Equal(decimal.NewFromFloat(200.00))).To(BeTrue())
		})

		It("should handle multiple movements in same currency", func() {
			movements := []goserver.Movement{
				{Amount: decimal.NewFromFloat(100.00), CurrencyId: "CZK", AccountId: "acc1"},
				{Amount: decimal.NewFromFloat(50.00), CurrencyId: "CZK", AccountId: "acc3"},
				{Amount: decimal.NewFromFloat(-30.00), CurrencyId: "CZK", AccountId: "acc2"},
			}
			increases := GetIncreases(movements)
			// pos: 150, neg: 30. GetIncreases returns max(pos, neg) which is 150
			Expect(increases["CZK"].Equal(decimal.NewFromFloat(150.00))).To(BeTrue())
		})

		It("should NOT have floating point precision issues with Decimal", func() {
			// 0.1 + 0.2 is exactly 0.3 in decimal.Decimal
			movements := []goserver.Movement{
				{Amount: decimal.NewFromFloat(0.1), CurrencyId: "CZK", AccountId: "acc1"},
				{Amount: decimal.NewFromFloat(0.2), CurrencyId: "CZK", AccountId: "acc1"},
			}
			increases := GetIncreases(movements)
			Expect(increases["CZK"].Equal(decimal.NewFromFloat(0.3))).To(BeTrue())
		})
	})

	Describe("IsDuplicate", func() {
		t1Date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		t1Movements := []goserver.Movement{
			{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "acc1"},
			{Amount: decimal.NewFromInt(-100), CurrencyId: "CZK", AccountId: "acc2"},
		}

		It("should identify duplicates within 2 days and same amounts", func() {
			t2Date := t1Date.Add(24 * time.Hour)
			t2Movements := []goserver.Movement{
				{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "acc3"},
				{Amount: decimal.NewFromInt(-100), CurrencyId: "CZK", AccountId: "acc4"},
			}
			Expect(IsDuplicate(t1Date, t1Movements, t2Date, t2Movements)).To(BeTrue())
		})

		It("should not identify as duplicate if outside 2 days", func() {
			t2Date := t1Date.Add(3 * 24 * time.Hour)
			Expect(IsDuplicate(t1Date, t1Movements, t2Date, t1Movements)).To(BeFalse())
		})

		It("should not identify as duplicate if amounts differ significantly", func() {
			t2Movements := []goserver.Movement{
				{Amount: decimal.NewFromFloat(100.02), CurrencyId: "CZK", AccountId: "acc1"},
				{Amount: decimal.NewFromFloat(-100.02), CurrencyId: "CZK", AccountId: "acc2"},
			}
			Expect(IsDuplicate(t1Date, t1Movements, t1Date, t2Movements)).To(BeFalse())
		})

		It("should identify as duplicate if amounts differ by less than 0.01", func() {
			t2Movements := []goserver.Movement{
				{Amount: decimal.NewFromFloat(100.005), CurrencyId: "CZK", AccountId: "acc1"},
				{Amount: decimal.NewFromFloat(-100.005), CurrencyId: "CZK", AccountId: "acc2"},
			}
			Expect(IsDuplicate(t1Date, t1Movements, t1Date, t2Movements)).To(BeTrue())
		})
	})
})

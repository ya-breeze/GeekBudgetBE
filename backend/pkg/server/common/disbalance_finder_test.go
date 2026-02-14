package common_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

var _ = Describe("DisbalanceFinder", func() {
	var (
		accID   = "acc-1"
		currID  = "eur"
		now     = time.Now()
		mockTxs []goserver.Transaction
	)

	BeforeEach(func() {
		mockTxs = []goserver.Transaction{
			{
				Id:          "tx-1",
				Date:        now,
				Description: "Lunch",
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromFloat(-15.50)},
				},
			},
			{
				Id:          "tx-2",
				Date:        now,
				Description: "Salary",
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromFloat(1000.00)},
				},
			},
			{
				Id:          "tx-3",
				Date:        now,
				Description: "Rent",
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromFloat(-500.00)},
				},
			},
			{
				Id:          "tx-4",
				Date:        now,
				Description: "Coffee",
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromFloat(-4.50)},
				},
			},
			{
				Id:          "tx-5",
				Date:        now,
				Description: "Refund",
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromFloat(20.00)},
				},
			},
		}
	})

	It("finds a single transaction match", func() {
		target := decimal.NewFromFloat(-15.50)
		result := common.AnalyzeDisbalance(target, mockTxs, accID, currID)

		Expect(result.Candidates).To(HaveLen(1))
		Expect(result.Candidates[0].Type).To(Equal("exact_single"))
		Expect(result.Candidates[0].Transactions).To(HaveLen(1))
		Expect(result.Candidates[0].Transactions[0].Id).To(Equal("tx-1"))
	})

	It("finds a pair match", func() {
		target := decimal.NewFromFloat(-20.00) // Lunch (-15.50) + Coffee (-4.50)
		result := common.AnalyzeDisbalance(target, mockTxs, accID, currID)

		Expect(result.Candidates).To(ContainElement(HaveField("Type", "exact_pair")))
		for _, c := range result.Candidates {
			if c.Type == "exact_pair" {
				Expect(c.Transactions).To(HaveLen(2))
				ids := []string{c.Transactions[0].Id, c.Transactions[1].Id}
				Expect(ids).To(ContainElements("tx-1", "tx-4"))
			}
		}
	})

	It("finds a subset match (3+ items)", func() {
		// Lunch (-15.50) + Coffee (-4.50) + Refund (20.00) = 0
		target := decimal.NewFromFloat(0.00)
		result := common.AnalyzeDisbalance(target, mockTxs, accID, currID)

		found := false
		for _, c := range result.Candidates {
			if c.Type == "exact_subset" {
				Expect(c.Transactions).To(HaveLen(3))
				found = true
			}
		}
		Expect(found).To(BeTrue())
	})

	It("returns empty candidates when no match is found", func() {
		target := decimal.NewFromFloat(1234.56)
		result := common.AnalyzeDisbalance(target, mockTxs, accID, currID)
		Expect(result.Candidates).To(BeEmpty())
	})

	It("handles empty transaction list", func() {
		result := common.AnalyzeDisbalance(decimal.NewFromInt(100), []goserver.Transaction{}, accID, currID)
		Expect(result.TransactionCount).To(Equal(int32(0)))
		Expect(result.Candidates).To(BeEmpty())
	})

	It("handles performance with 50 transactions", func() {
		largeTxs := make([]goserver.Transaction, 50)
		for i := 0; i < 50; i++ {
			largeTxs[i] = goserver.Transaction{
				Id:   time.Now().String(),
				Date: now,
				Movements: []goserver.Movement{
					{AccountId: accID, CurrencyId: currID, Amount: decimal.NewFromInt(int64(i + 1))},
				},
			}
		}
		target := decimal.NewFromInt(100)

		start := time.Now()
		result := common.AnalyzeDisbalance(target, largeTxs, accID, currID)
		elapsed := time.Since(start)

		Expect(elapsed.Seconds()).To(BeNumerically("<", 1.0))
		Expect(result.Candidates).NotTo(BeEmpty())
	})
})

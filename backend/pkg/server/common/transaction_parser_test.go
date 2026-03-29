package common_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

var _ = Describe("ParseTransactionText", func() {
	var accounts []goserver.Account
	var currencies []goserver.Currency
	today := time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC)

	BeforeEach(func() {
		accounts = []goserver.Account{
			{Id: "acc-fio", Name: "FIO Savings"},
			{Id: "acc-others", Name: "Others"},
			{Id: "acc-kb", Name: "KB Current"},
		}
		currencies = []goserver.Currency{
			{Id: "cur-czk", Name: "CZK"},
			{Id: "cur-eur", Name: "EUR"},
		}
	})

	It("parses full canonical form with slash date", func() {
		result, warnings := common.ParseTransactionText(
			"2026/03/22 100 CZK from FIO Savings to Others",
			accounts, currencies, today,
		)
		Expect(warnings).To(BeEmpty())
		Expect(result.Date).To(Equal(today))
		Expect(result.Movements).To(HaveLen(2))
		Expect(result.Movements[0].AccountId).To(Equal("acc-fio"))
		Expect(result.Movements[0].Amount).To(Equal(decimal.NewFromInt(-100)))
		Expect(result.Movements[0].CurrencyId).To(Equal("cur-czk"))
		Expect(result.Movements[1].AccountId).To(Equal("acc-others"))
		Expect(result.Movements[1].Amount).To(Equal(decimal.NewFromInt(100)))
		Expect(result.Movements[1].CurrencyId).To(Equal("cur-czk"))
	})

	It("uses dash-separated date", func() {
		result, _ := common.ParseTransactionText(
			"2026-01-15 200 CZK from FIO Savings to Others",
			accounts, currencies, today,
		)
		Expect(result.Date).To(Equal(time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)))
	})

	It("defaults date to today when omitted", func() {
		result, _ := common.ParseTransactionText(
			"50 EUR from KB Current to Others",
			accounts, currencies, today,
		)
		Expect(result.Date).To(Equal(today))
		Expect(result.Movements[0].Amount.Equal(decimal.NewFromFloat(-50))).To(BeTrue())
		Expect(result.Movements[0].CurrencyId).To(Equal("cur-eur"))
	})

	It("handles decimal amounts", func() {
		result, _ := common.ParseTransactionText(
			"50.5 EUR from KB Current to Others",
			accounts, currencies, today,
		)
		Expect(result.Movements[0].Amount.Equal(decimal.NewFromFloat(-50.5))).To(BeTrue())
	})

	It("handles missing 'from' clause — single to-movement", func() {
		result, _ := common.ParseTransactionText(
			"100 CZK to Others",
			accounts, currencies, today,
		)
		Expect(result.Movements).To(HaveLen(1))
		Expect(result.Movements[0].AccountId).To(Equal("acc-others"))
		Expect(result.Movements[0].Amount).To(Equal(decimal.NewFromInt(100)))
	})

	It("adds warning when account not found", func() {
		_, warnings := common.ParseTransactionText(
			"100 CZK from Unknown Bank to Others",
			accounts, currencies, today,
		)
		Expect(warnings).To(ContainElement(ContainSubstring("Unknown Bank")))
	})

	It("adds warning when currency not found", func() {
		_, warnings := common.ParseTransactionText(
			"100 USD from FIO Savings to Others",
			accounts, currencies, today,
		)
		Expect(warnings).To(ContainElement(ContainSubstring("USD")))
	})

	It("adds warning on partial account match and still resolves", func() {
		result, warnings := common.ParseTransactionText(
			"100 CZK from fio to Others",
			accounts, currencies, today,
		)
		Expect(warnings).To(ContainElement(ContainSubstring("partial")))
		Expect(result.Movements[0].AccountId).To(Equal("acc-fio"))
	})

	It("captures leftover text as description", func() {
		result, _ := common.ParseTransactionText(
			"100 CZK from FIO Savings to Others groceries",
			accounts, currencies, today,
		)
		Expect(result.Description).To(Equal("groceries"))
	})

	It("is case-insensitive for account and currency matching", func() {
		result, warnings := common.ParseTransactionText(
			"100 czk from fio savings to others",
			accounts, currencies, today,
		)
		Expect(warnings).To(BeEmpty())
		Expect(result.Movements[0].CurrencyId).To(Equal("cur-czk"))
		Expect(result.Movements[0].AccountId).To(Equal("acc-fio"))
		Expect(result.Movements[1].AccountId).To(Equal("acc-others"))
	})
})

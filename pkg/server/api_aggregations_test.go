package server_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Aggregation API", func() {
	log := test.CreateTestLogger()
	accounts := test.PrepareAccounts()
	currencies := test.PrepareCurrencies()
	transactions := test.PrepareTransactions(accounts, currencies)
	dateFrom := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)

	It("aggregate expenses", func() {
		sut := server.Aggregate(accounts, transactions, dateFrom, dateTo, utils.GranularityMonth, log)
		Expect(sut.From.UnixMilli()).To(Equal(time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC).UnixMilli()))
		Expect(sut.To.UnixMilli()).To(Equal(time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC).UnixMilli()))

		Expect(sut.Intervals).To(HaveLen(2))

		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id))
		Expect(sut.Currencies[0].Accounts).To(HaveLen(2))
		Expect(sut.Currencies).To(HaveLen(1))
		Expect(sut.Currencies[0].CurrencyId).To(Equal(currencies[0].Id))

		Expect(sut.Currencies[0].Accounts[0].AccountId).To(Equal(accounts[2].Id))
		Expect(sut.Currencies[0].Accounts[0].Amounts).To(HaveLen(2))
		Expect(sut.Currencies[0].Accounts[0].Amounts[0]).To(Equal(450.0))
		Expect(sut.Currencies[0].Accounts[0].Amounts[1]).To(Equal(10.0))

		Expect(sut.Currencies[0].Accounts[1].AccountId).To(Equal(accounts[4].Id))
		Expect(sut.Currencies[0].Accounts[1].Amounts).To(HaveLen(2))
		Expect(sut.Currencies[0].Accounts[1].Amounts[0]).To(Equal(300.0))
		Expect(sut.Currencies[0].Accounts[1].Amounts[1]).To(Equal(250.0))
	})
})

package bankimporters_test

import (
	"context"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("RevolutConverter Russian CSV", func() {
	var (
		converter *bankimporters.RevolutConverter
		logger    *slog.Logger
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		cp := bankimporters.NewSimpleCurrencyProvider([]goserver.Currency{
			{Id: "CZK-ID", Name: "CZK"},
			{Id: "EUR-ID", Name: "EUR"},
		})
		var err error
		converter, err = bankimporters.NewRevolutConverter(
			logger,
			goserver.BankImporter{
				AccountId:    "test-account-id",
				FeeAccountId: "fee-account-id",
			},
			cp,
		)
		Expect(err).ToNot(HaveOccurred())
	})

	It("correctly parses Russian CSV from Revolut", func() {
		csvData := `Тип,Продукт,Дата начала,Дата выполнения,Описание,Сумма,Комиссия,Валюта,State,Остаток средств
Пополнение,Текущий,2025-12-25 14:08:28,2025-12-25 14:08:39,Top-up from *1234,1000.00,0.00,CZK,ВЫПОЛНЕНО,1032.31
Платеж по карте,Текущий,2025-12-25 18:07:35,2025-12-26 02:38:23,Store A,-100.00,0.00,CZK,ВЫПОЛНЕНО,932.31
Обмен валюты,Текущий,2025-12-26 13:22:17,2025-12-26 13:22:17,Обменено на EUR,-500.00,0.00,CZK,ВЫПОЛНЕНО,432.31
Обмен валюты,Текущий,2025-12-26 13:22:17,2025-12-26 13:22:17,Обменено на EUR,20.00,0.00,EUR,ВЫПОЛНЕНО,20.00
Списать,Текущий,2025-12-26 13:20:42,2025-12-26 13:20:42,Fee description,0.00,10.00,CZK,ВЫПОЛНЕНО,422.31`

		info, transactions, err := converter.ParseTransactions(context.Background(), "csv", csvData)
		Expect(err).ToNot(HaveOccurred())
		Expect(info).ToNot(BeNil())

		// 5 records in CSV, but 2 are joined (exchange), so 4 transactions total
		Expect(transactions).To(HaveLen(4))

		// Check Top-up
		Expect(transactions[0].Description).To(Equal("Пополнение: Top-up from *1234"))
		Expect(transactions[0].Movements).To(HaveLen(2))
		Expect(transactions[0].Movements[1].Amount.Equal(decimal.NewFromInt(1000))).To(BeTrue())
		Expect(transactions[0].Movements[1].CurrencyId).To(Equal("CZK-ID"))

		// Check Card Payment
		Expect(transactions[1].Description).To(Equal("Платеж по карте: Store A"))
		Expect(transactions[1].Movements).To(HaveLen(2))
		Expect(transactions[1].Movements[1].Amount.Equal(decimal.NewFromInt(-100))).To(BeTrue())

		// Check Joined Exchange
		Expect(transactions[2].Description).To(Equal("Обмен валюты: Обменено на EUR"))
		Expect(transactions[2].Movements).To(HaveLen(2))
		// Source movement (from second record usually in joinExchanges logic)
		// Actually joinExchanges:
		// transactions[i].Movements[0] = transactions[j].Movements[1]
		// In our case: i is -500 CZK, j is 20 EUR
		// transactions[i].Movements[0] was -(-20) = 20 EUR.
		// Wait, let's re-read convertToTransaction:
		// if !amount.IsZero() {
		// 	res.Movements = append(res.Movements, goserver.Movement{
		// 		Amount:     amount.Neg(), // negation of amount
		// 		CurrencyId: strCurrencyID,
		// 	})
		// }
		// if !remainingAmount.IsZero() {
		// 	res.Movements = append(res.Movements, goserver.Movement{
		// 		AccountId:  fc.bankImporter.AccountId,
		// 		Amount:     remainingAmount, // original amount
		// 		CurrencyId: strCurrencyID,
		// 	})
		// }
		// So i had [-(-500)=500, -500], j had [-20, 20].
		// After join: i.Movements[0] = j.Movements[1] = 20.
		Expect(transactions[2].Movements[0].Amount.Equal(decimal.NewFromInt(20))).To(BeTrue())
		Expect(transactions[2].Movements[0].CurrencyId).To(Equal("EUR-ID"))
		Expect(transactions[2].Movements[1].Amount.Equal(decimal.NewFromInt(-500))).To(BeTrue())
		Expect(transactions[2].Movements[1].CurrencyId).To(Equal("CZK-ID"))

		// Check Fee record ("Списать")
		Expect(transactions[3].Description).To(Equal("Списать: Fee description"))
		Expect(transactions[3].Movements).To(HaveLen(2))
		Expect(transactions[3].Movements[0].Amount.Equal(decimal.NewFromInt(-10))).To(BeTrue())
		Expect(transactions[3].Movements[0].AccountId).To(Equal("test-account-id"))
		Expect(transactions[3].Movements[1].Amount.Equal(decimal.NewFromInt(10))).To(BeTrue())
		Expect(transactions[3].Movements[1].AccountId).To(Equal("fee-account-id"))

		// Check Balances in info
		// CZK:
		// First balance: 1032.31 (from first record)
		// Last balance: 422.31 (from last CZK record)
		// Opening balance = 1032.31 - 1000 = 32.31
		// Closing balance = 1032.31 (newest first) or 422.31 (oldest first)?
		// The dates are in ascending order in my sample.
		// firstDate: 2025-12-25 14:08:28
		// lastDate: 2025-12-26 13:20:42
		// So oldest first. Closing = lastBalance = 422.31.
		// Opening = firstBalance - firstAmount = 1032.31 - 1000 = 32.31
		foundCZK := false
		for _, b := range info.Balances {
			if b.CurrencyId == "CZK-ID" {
				foundCZK = true
				Expect(b.ClosingBalance.String()).To(Equal("422.31"))
				Expect(b.OpeningBalance.String()).To(Equal("32.31"))
			}
		}
		Expect(foundCZK).To(BeTrue())
	})

	It("skips non-completed transactions", func() {
		csvData := `Тип,Продукт,Дата начала,Дата выполнения,Описание,Сумма,Комиссия,Валюта,State,Остаток средств
Платеж по карте,Текущий,2025-12-25 18:07:35,,Store B,-100.00,0.00,CZK,REVERTED,932.31`

		_, transactions, err := converter.ParseTransactions(context.Background(), "csv", csvData)
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).To(BeEmpty())
	})
})

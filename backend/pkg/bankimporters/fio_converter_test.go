package bankimporters_test

import (
	"context"
	"log/slog"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("FIO converter", func() {
	var (
		err error
		fc  *bankimporters.FioConverter
	)
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	BeforeEach(func() {
		cp := bankimporters.NewSimpleCurrencyProvider([]goserver.Currency{
			{Id: "__CZK_ID__", Name: "CZK"},
			{Id: "__EUR_ID__", Name: "EUR"},
		})
		fc, err = bankimporters.NewFioConverter(
			log,
			goserver.BankImporter{
				AccountId: "__accountID__",
			}, cp)
		Expect(err).ToNot(HaveOccurred())
	})

	It("parses FIO response with plain string currency in info", func() {
		// This JSON has currency as a plain string "CZK" in the info section,
		// which matches the new API format that caused the original bug
		data := []byte(`{
			"accountStatement": {
				"info": {
					"accountId": "123456789",
					"bankId": "2010",
					"currency": {"value": "CZK", "name": "currency", "id": 14},
					"iban": "CZ1234567890123456789012",
					"bic": "FIOBCZPPXXX",
					"openingBalance": {"value": 1000.50, "name": "openingBalance", "id": 1},
					"closingBalance": {"value": 1500.75, "name": "closingBalance", "id": 2},
					"dateStart": "2024-01-01+0100",
					"dateEnd": "2024-01-31+0100",
					"yearId": 2024,
					"idFrom": 1000000001,
					"idTo": 1000000002
				},
				"transactionList": {
					"transaction": [
						{
							"column0": {"value": "2024-01-15+0100", "name": "Datum", "id": 0},
							"column22": {"value": 1000000001, "name": "ID pohybu", "id": 22},
							"column1": {"value": 100.00, "name": "Objem", "id": 1},
							"column14": {"value": "CZK", "name": "Měna", "id": 14},
							"column2": {"value": "987654321", "name": "Protiúčet", "id": 2},
							"column3": {"value": "0100", "name": "Kód banky", "id": 3},
							"column5": {"value": "12345", "name": "VS", "id": 5},
							"column10": {"value": "Test Partner", "name": "Název protiúčtu", "id": 10},
							"column12": {"value": "Test Bank", "name": "Název banky", "id": 12},
							"column7": {"value": "", "name": "Uživatelská identifikace", "id": 7},
							"column8": {"value": "Příjem", "name": "Typ", "id": 8},
							"column9": {"value": "User", "name": "Provedl", "id": 9},
							"column25": {"value": "Test payment", "name": "Komentář", "id": 25},
							"column16": {"value": "", "name": "Zpráva pro příjemce", "id": 16}
						}
					]
				}
			}
		}`)

		info, transactions, err := fc.ParseTransactions(context.Background(), data)
		Expect(err).ToNot(HaveOccurred())

		// Verify account info - especially that currency is correctly parsed
		Expect(info.AccountId).To(Equal("123456789"))
		Expect(info.BankId).To(Equal("2010"))
		Expect(info.Balances).To(HaveLen(1))
		Expect(info.Balances[0].OpeningBalance.Equal(decimal.NewFromFloat(1000.50))).To(BeTrue())
		Expect(info.Balances[0].ClosingBalance.Equal(decimal.NewFromFloat(1500.75))).To(BeTrue())
		Expect(info.Balances[0].CurrencyId).To(Equal("__CZK_ID__"))

		// Verify transactions
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].PartnerName).To(Equal("Test Partner"))
		Expect(transactions[0].PartnerAccount).To(Equal("987654321/0100 vs=12345"))
		Expect(transactions[0].Description).To(ContainSubstring("Test payment"))
		Expect(transactions[0].Movements).To(HaveLen(2))
		Expect(transactions[0].Movements[0].Amount.Equal(decimal.NewFromFloat(-100.00))).To(BeTrue())
		Expect(transactions[0].Movements[0].CurrencyId).To(Equal("__CZK_ID__"))
		Expect(transactions[0].Movements[1].Amount.Equal(decimal.NewFromFloat(100.00))).To(BeTrue())
		Expect(transactions[0].Movements[1].AccountId).To(Equal("__accountID__"))
	})

	It("handles empty transaction list", func() {
		data := []byte(`{
			"accountStatement": {
				"info": {
					"accountId": "111222333",
					"bankId": "2010",
					"currency": "CZK",
					"iban": "CZ1112223330001112223330",
					"bic": "FIOBCZPPXXX",
					"openingBalance": {"value": 100.00, "name": "openingBalance", "id": 1},
					"closingBalance": {"value": 100.00, "name": "closingBalance", "id": 2},
					"dateStart": "2024-01-01+0100",
					"dateEnd": "2024-01-01+0100",
					"yearId": 2024,
					"idFrom": 0,
					"idTo": 0
				},
				"transactionList": {
					"transaction": []
				}
			}
		}`)

		info, transactions, err := fc.ParseTransactions(context.Background(), data)
		Expect(err).ToNot(HaveOccurred())
		Expect(info.AccountId).To(Equal("111222333"))
		Expect(transactions).To(BeEmpty())
	})

	It("parses FIO response with plain values (new API format)", func() {
		// This JSON uses plain values instead of objects - the format that caused the original bug
		// currency is a plain string, openingBalance and closingBalance are plain numbers
		data := []byte(`{
			"accountStatement": {
				"info": {
					"accountId": "444555666",
					"bankId": "2010",
					"currency": "CZK",
					"iban": "CZ4445556660004445556660",
					"bic": "FIOBCZPPXXX",
					"openingBalance": 2500.00,
					"closingBalance": 3000.50,
					"dateStart": "2024-02-01+0100",
					"dateEnd": "2024-02-29+0100",
					"yearId": 2024,
					"idFrom": 3000000001,
					"idTo": 3000000001
				},
				"transactionList": {
					"transaction": [
						{
							"column0": {"value": "2024-02-15+0100", "name": "Datum", "id": 0},
							"column22": {"value": 3000000001, "name": "ID pohybu", "id": 22},
							"column1": {"value": 500.50, "name": "Objem", "id": 1},
							"column14": {"value": "CZK", "name": "Měna", "id": 14},
							"column2": {"value": "", "name": "Protiúčet", "id": 2},
							"column3": {"value": "", "name": "Kód banky", "id": 3},
							"column5": {"value": "", "name": "VS", "id": 5},
							"column10": {"value": "Plain Value Test", "name": "Název protiúčtu", "id": 10},
							"column12": {"value": "", "name": "Název banky", "id": 12},
							"column7": {"value": "", "name": "Uživatelská identifikace", "id": 7},
							"column8": {"value": "Příjem", "name": "Typ", "id": 8},
							"column9": {"value": "", "name": "Provedl", "id": 9},
							"column25": {"value": "Test with plain values", "name": "Komentář", "id": 25},
							"column16": {"value": "", "name": "Zpráva pro příjemce", "id": 16}
						}
					]
				}
			}
		}`)

		info, transactions, err := fc.ParseTransactions(context.Background(), data)
		Expect(err).ToNot(HaveOccurred())

		// Verify plain value formats are correctly parsed
		Expect(info.AccountId).To(Equal("444555666"))
		Expect(info.Balances[0].OpeningBalance.Equal(decimal.NewFromFloat(2500.00))).To(BeTrue())
		Expect(info.Balances[0].ClosingBalance.Equal(decimal.NewFromFloat(3000.50))).To(BeTrue())
		Expect(info.Balances[0].CurrencyId).To(Equal("__CZK_ID__"))

		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].PartnerName).To(Equal("Plain Value Test"))
	})

	It("parses FIO response from fio_example.json format", func() {
		// This JSON matches the format in fio_example.json (numeric timestamps, nulls)
		data := []byte(`{
    "accountStatement": {
        "info": {
            "accountId": "2400222222",
            "bankId": "2010",
            "currency": "CZK",
            "iban": "CZ7920100000002400222222",
            "bic": "FIOBCZPPXXX",
            "openingBalance": 195.00,
            "closingBalance": 195.01,
            "dateStart": 1340661600000,
            "dateEnd": 1341007200000,
            "yearList": null,
            "idList": null,
            "idFrom": 1148734530,
            "idTo": 1149190193,
            "idLastDownload": 1149190192
        },
        "transactionList": {
            "transaction": [
                {
                    "column22": {
                        "value": 1148734530,
                        "name": "IDpohybu",
                        "id": 22
                    },
                    "column0": {
                        "value": 1340661600000,
                        "name": "Datum",
                        "id": 0
                    },
                    "column1": {
                        "value": 1.00,
                        "name": "Objem",
                        "id": 1
                    },
                    "column14": {
                        "value": "CZK",
                        "name": "Měna",
                        "id": 14
                    },
                    "column2": {
                        "value": "2900233333",
                        "name": "Protiúčet",
                        "id": 2
                    },
                    "column10": {
                        "value": "Pavel, Novák",
                        "name": "Názevprotiúčtu",
                        "id": 10
                    },
                    "column3": {
                        "value": "2010",
                        "name": "Kódbanky",
                        "id": 3
                    },
                    "column12": {
                        "value": "Fio banka, a.s.",
                        "name": "Názevbanky",
                        "id": 12
                    },
                    "column4": {
                        "value": "0558",
                        "name": "KS",
                        "id": 4
                    },
                    "column5": null,
                    "column6": null,
                    "column7": null,
                    "column16": null,
                    "column8": {
                        "value": "Příjem převodem uvnitřbanky",
                        "name": "Typ",
                        "id": 8
                    },
                    "column9": null,
                    "column18": null,
                    "column25": null,
                    "column26": null,
                    "column17": {
                        "value": 2105685816,
                        "name": "IDpokynu",
                        "id": 17
                    }
                }
            ]
        }
    }
}`)

		info, transactions, err := fc.ParseTransactions(context.Background(), data)
		Expect(err).ToNot(HaveOccurred())

		// Verify account info
		Expect(info.AccountId).To(Equal("2400222222"))
		Expect(info.Balances[0].OpeningBalance.Equal(decimal.NewFromFloat(195.00))).To(BeTrue())
		Expect(info.Balances[0].ClosingBalance.Equal(decimal.NewFromFloat(195.01))).To(BeTrue())
		Expect(info.Balances[0].CurrencyId).To(Equal("__CZK_ID__"))

		// Verify transaction date was correctly parsed from timestamp
		Expect(transactions).To(HaveLen(1))
		// 1340661600000 is 2012-06-25 22:00:00 UTC, which is 2012-06-26 00:00:00 CEST
		Expect(transactions[0].Date.Year()).To(Equal(2012))
		Expect(transactions[0].Date.Month()).To(Equal(time.June))
		Expect(transactions[0].Date.Day()).To(Equal(26))

		Expect(transactions[0].Movements[1].Amount.Equal(decimal.NewFromFloat(1.00))).To(BeTrue())
		Expect(transactions[0].PartnerName).To(Equal("Pavel, Novák"))
	})

	It("filters out zero-amount movements", func() {
		data := []byte(`{
			"accountStatement": {
				"info": {
					"accountId": "123",
					"bankId": "2010",
					"currency": "CZK",
					"openingBalance": 100.00,
					"closingBalance": 100.00
				},
				"transactionList": {
					"transaction": [
						{
							"column0": {"value": "2024-01-15+0100", "name": "Datum", "id": 0},
							"column22": {"value": 1, "name": "ID", "id": 22},
							"column1": {"value": 0.00, "name": "Objem", "id": 1},
							"column14": {"value": "CZK", "name": "Měna", "id": 14},
							"column8": {"value": "Zero Amount", "name": "Typ", "id": 8}
						}
					]
				}
			}
		}`)

		_, transactions, err := fc.ParseTransactions(context.Background(), data)
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].Movements).To(BeEmpty())
	})
})

package bankimporters_test

import (
	"context"
	"encoding/csv"
	"log/slog"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("KB converter", func() {
	var (
		err error
		rc  *bankimporters.KBConverter
	)
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	loc, _ := time.LoadLocation("Europe/Prague")

	BeforeEach(func() {
		cp := bankimporters.NewSimpleCurrencyProvider([]goserver.Currency{
			{Id: "__CZK_ID__", Name: "CZK"},
			{Id: "__EUR_ID__", Name: "EUR"},
			{Id: "__USD_ID__", Name: "USD"},
		})
		rc, err = bankimporters.NewKBConverter(
			log,
			goserver.BankImporter{
				AccountId: "__accountID__",
			}, cp)
		Expect(err).ToNot(HaveOccurred())
	})

	DescribeTable("parses KB transactions",
		func(data string, expectedTransaction goserver.TransactionNoId) {
			r := csv.NewReader(strings.NewReader(data))
			r.Comma = ';'
			record, err := r.Read()
			Expect(err).ToNot(HaveOccurred())
			record = rc.PrepareRow(record)
			transaction, err := rc.ConvertToTransaction(context.Background(), record)
			Expect(err).ToNot(HaveOccurred())
			Expect(transaction.Date).To(Equal(expectedTransaction.Date))
			Expect(transaction.Description).To(Equal(expectedTransaction.Description))
			Expect(transaction.Place).To(Equal(expectedTransaction.Place))
			Expect(transaction.Tags).To(Equal(expectedTransaction.Tags))
			Expect(transaction.PartnerName).To(Equal(expectedTransaction.PartnerName))
			Expect(transaction.PartnerAccount).To(Equal(expectedTransaction.PartnerAccount))
			Expect(transaction.PartnerInternalId).To(Equal(expectedTransaction.PartnerInternalId))
			Expect(transaction.Extra).To(Equal(expectedTransaction.Extra))
			Expect(transaction.ExternalIds).To(Equal(expectedTransaction.ExternalIds))
			Expect(transaction.Movements).To(HaveLen(len(expectedTransaction.Movements)))
			Expect(transaction.Movements[0].Amount.Equal(expectedTransaction.Movements[0].Amount)).To(BeTrue())
			Expect(transaction.Movements[0].CurrencyId).To(Equal(expectedTransaction.Movements[0].CurrencyId))
			Expect(transaction.Movements[0].AccountId).To(Equal(expectedTransaction.Movements[0].AccountId))
			Expect(transaction.Movements[1].Amount.Equal(expectedTransaction.Movements[1].Amount)).To(BeTrue())
			Expect(transaction.Movements[1].CurrencyId).To(Equal(expectedTransaction.Movements[1].CurrencyId))
			Expect(transaction.Movements[1].AccountId).To(Equal(expectedTransaction.Movements[1].AccountId))
		},
		Entry("transaction N1",
			`"26.09.2024";"26.09.2024";"123/45";"companyname";"-12345,00";"CZK";"";"";"";"9";"138";"0";`+
				`"externalid";"incoming";"description user";"message";"reference";"BIC";"fee"`,
			goserver.TransactionNoId{
				Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, loc),
				Description:    "incoming: description user; message; reference",
				PartnerAccount: "123/45; VS:9; KS:138",
				PartnerName:    "companyname",
				Tags:           []string{"kb"},
				ExternalIds:    []string{"externalid"},
				Movements: []goserver.Movement{
					{
						Amount:     decimal.NewFromInt(12345),
						CurrencyId: "__CZK_ID__",
					},
					{
						AccountId:  "__accountID__",
						Amount:     decimal.NewFromInt(-12345),
						CurrencyId: "__CZK_ID__",
					},
				},
			}),
	)

	It("parses KB file with header", func() {
		data := `KB+, vypis v csv. formatu;;;;;;;;;;;;;;;;;;
Datum vytvoreni souboru;21.12.2025;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;
Cislo uctu;123-123123;;;;;;;;;;;;;;;;;
Mena uctu / Hlavni mena uctu;CZK;;;;;;;;;;;;;;;;;
Mena vypisu;;;;;;;;;;;;;;;;;;
IBAN;CZ123123123;;;;;;;;;;;;;;;;;
Nazev uctu;A B;;;;;;;;;;;;;;;;;
Vypis od;01.01.2025;;;;;;;;;;;;;;;;;
Vypis do;20.12.2025;;;;;;;;;;;;;;;;;
Cislo vypisu;;;;;;;;;;;;;;;;;;
Pocet polozek;55;;;;;;;;;;;;;;;;;
Mena;CZK;EUR;USD;GBP;AUD;BGN;CAD;CHF;DKK;HUF;JPY;NOK;PLN;RON;SEK;;;
Pocatecni zustatek;321321,21;;;;;;;;;;;;;;;;;
Konecny zustatek;123123,12;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;
Datum zauctovani;Datum provedeni;Protistrana;Nazev protiuctu;Castka;Mena;Originalni castka;Originalni mena;Smenny kurz;VS;KS;SS;Identifikace transakce;Typ transakce;Popis pro me;Zprava pro prijemce;Reference platby;BIC / SWIFT;Poplatek
26.09.2024;26.09.2024;123/45;Employer Corp;12345,00;CZK;;;;9;138;0;tx1;Incoming payment;Salary;September;;BIC1;
26.09.2024;26.09.2024;987/65;Landlord;-10000,00;CZK;;;;0;0;0;tx2;Outgoing payment;Rent;October;;BIC2;
30.09.2024;30.09.2024;;Bank;-15,00;CZK;;;;0;0;0;tx3;Fee;Account Fee;;;;`

		info, transactions, err := rc.ParseAndImport("csv", data)
		Expect(err).ToNot(HaveOccurred())
		Expect(info.AccountId).To(Equal("123-123123"))
		Expect(info.Balances).To(HaveLen(1))
		Expect(info.Balances[0].ClosingBalance.Equal(decimal.NewFromFloat(123123.12))).To(BeTrue())
		Expect(info.Balances[0].CurrencyId).To(Equal("__CZK_ID__"))

		Expect(transactions).To(HaveLen(3))

		// Check Incoming Payment
		Expect(transactions[0].ExternalIds).To(ContainElement("tx1"))
		Expect(transactions[0].PartnerName).To(Equal("Employer Corp"))
		Expect(transactions[0].Movements[0].Amount.Equal(decimal.NewFromFloat(-12345.00))).To(BeTrue())
		Expect(transactions[0].Description).To(ContainSubstring("Incoming payment"))
		Expect(transactions[0].Description).To(ContainSubstring("Salary"))
		Expect(transactions[0].Description).To(ContainSubstring("September"))

		// Check Outgoing Payment
		Expect(transactions[1].ExternalIds).To(ContainElement("tx2"))
		Expect(transactions[1].PartnerName).To(Equal("Landlord"))
		Expect(transactions[1].Movements[0].Amount.Equal(decimal.NewFromFloat(10000.00))).To(BeTrue())
		Expect(transactions[1].Movements[1].Amount.Equal(decimal.NewFromFloat(-10000.00))).To(BeTrue())

		// Check Fee
		Expect(transactions[2].ExternalIds).To(ContainElement("tx3"))
		Expect(transactions[2].PartnerName).To(Equal("Bank"))
		Expect(transactions[2].Description).To(ContainSubstring("Fee"))
		Expect(transactions[2].Movements[1].Amount.Equal(decimal.NewFromFloat(-15.00))).To(BeTrue())
	})

	It("filters out zero-amount movements", func() {
		data := `KB+, vypis v csv. formatu;;;;;;;;;;;;;;;;;;
Datum vytvoreni souboru;21.12.2025;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;
Cislo uctu;123-123123;;;;;;;;;;;;;;;;;
Mena uctu / Hlavni mena uctu;CZK;;;;;;;;;;;;;;;;;
Mena vypisu;;;;;;;;;;;;;;;;;;
IBAN;CZ123123123;;;;;;;;;;;;;;;;;
Nazev uctu;A B;;;;;;;;;;;;;;;;;
Vypis od;01.01.2025;;;;;;;;;;;;;;;;;
Vypis do;20.12.2025;;;;;;;;;;;;;;;;;
Cislo vypisu;;;;;;;;;;;;;;;;;;
Pocet polozek;1;;;;;;;;;;;;;;;;;
Mena;CZK;EUR;USD;GBP;AUD;BGN;CAD;CHF;DKK;HUF;JPY;NOK;PLN;RON;SEK;;;
Pocatecni zustatek;100,00;;;;;;;;;;;;;;;;;
Konecny zustatek;100,00;;;;;;;;;;;;;;;;;
;;;;;;;;;;;;;;;;;;
Datum zauctovani;Datum provedeni;Protistrana;Nazev protiuctu;Castka;Mena;Originalni castka;Originalni mena;Smenny kurz;VS;KS;SS;Identifikace transakce;Typ transakce;Popis pro me;Zprava pro prijemce;Reference platby;BIC / SWIFT;Poplatek
26.09.2024;26.09.2024;123/45;Zero Corp;0,00;CZK;;;;0;0;0;tx0;Zero payment;Zero;Zero;;;`

		_, transactions, err := rc.ParseAndImport("csv", data)
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].Movements).To(BeEmpty())
	})
})

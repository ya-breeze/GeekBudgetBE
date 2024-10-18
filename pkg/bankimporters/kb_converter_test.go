package bankimporters_test

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		rc, err = bankimporters.NewKBConverter(
			log,
			goserver.BankImporter{
				AccountId: "__accountID__",
			}, []goserver.Currency{
				{Id: "__CZK_ID__", Name: "CZK"},
				{Id: "__EUR_ID__", Name: "EUR"},
				{Id: "__USD_ID__", Name: "USD"},
			})
		Expect(err).ToNot(HaveOccurred())
	})

	DescribeTable("parses KB transactions",
		func(data string, expectedTransaction goserver.TransactionNoId) {
			r := csv.NewReader(strings.NewReader(data))
			r.Comma = ';'
			record, err := r.Read()
			Expect(err).ToNot(HaveOccurred())
			record = rc.PrepareRow(record)
			transaction, err := rc.ConvertToTransaction(record)
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
			Expect(transaction.Movements[0].Amount).To(Equal(expectedTransaction.Movements[0].Amount))
			Expect(transaction.Movements[0].CurrencyId).To(Equal(expectedTransaction.Movements[0].CurrencyId))
			Expect(transaction.Movements[0].AccountId).To(Equal(expectedTransaction.Movements[0].AccountId))
			Expect(transaction.Movements[1].Amount).To(Equal(expectedTransaction.Movements[1].Amount))
			Expect(transaction.Movements[1].CurrencyId).To(Equal(expectedTransaction.Movements[1].CurrencyId))
			Expect(transaction.Movements[1].AccountId).To(Equal(expectedTransaction.Movements[1].AccountId))
		},
		Entry("transaction N1",
			`"26.09.2024";"26.09.2024";"123/45";"companyname";"+12345,00";"";"";"";"9";"138";"0";`+
				`"externalid";"incoming";"companyname";"personname";"9/0138/0/11/                       "`+
				`;"/VS9/SS/KS0138                     ";"B/O company ";"                                   ";`,
			goserver.TransactionNoId{
				Date:           time.Date(2024, 9, 26, 0, 0, 0, 0, loc),
				Description:    "incoming: companyname; personname; 9/0138/0/11//VS9/SS/KS0138B/O company",
				PartnerAccount: "123/45; VS:9; KS:138",
				PartnerName:    "companyname",
				Tags:           []string{"kb"},
				ExternalIds:    []string{"externalid"},
				Movements: []goserver.Movement{
					{
						Amount:     12345,
						CurrencyId: "__CZK_ID__",
					},
					{
						AccountId:  "__accountID__",
						Amount:     -12345,
						CurrencyId: "__CZK_ID__",
					},
				},
			}),
	)
})

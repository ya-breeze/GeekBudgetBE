package bankimporters_test

import (
	"context"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("RevolutConverter Repro", func() {
	var (
		converter *bankimporters.RevolutConverter
		logger    *slog.Logger
	)

	BeforeEach(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		var err error
		cp := bankimporters.NewSimpleCurrencyProvider([]goserver.Currency{
			{Id: "CZK-ID", Name: "CZK"},
			{Id: "USD-ID", Name: "USD"},
		})
		converter, err = bankimporters.NewRevolutConverter(
			logger,
			goserver.BankImporter{
				AccountId: "test-account-id",
			},
			cp,
		)
		Expect(err).ToNot(HaveOccurred())
	})

	It("generates different hashes for semantically identical transactions with different formatting", func() {
		// Two CSV records representing the same transaction, but one has quotes and the other doesn't
		// Revolut often changes export format slightly.

		// Original format
		csv1 := `Type,Product,Started Date,Completed Date,Description,Amount,Fee,Currency,State,Balance
CARD_PAYMENT,Current,2023-01-01 10:00:00,2023-01-02 10:00:00,Coffee Shop,-50.00,0.00,CZK,COMPLETED,1000.00`

		// Slightly different format (e.g. extra space in description)
		csv2 := `Type,Product,Started Date,Completed Date,Description,Amount,Fee,Currency,State,Balance
CARD_PAYMENT,Current,2023-01-01 10:00:00,2023-01-02 10:00:00,"Coffee Shop ",-50.00,0.00,CZK,COMPLETED,1000.00`

		info1, trans1, err := converter.ParseTransactions(context.Background(), "csv", csv1)
		Expect(err).ToNot(HaveOccurred())
		Expect(trans1).To(HaveLen(1))
		Expect(info1).ToNot(BeNil())

		info2, trans2, err := converter.ParseTransactions(context.Background(), "csv", csv2)
		Expect(err).ToNot(HaveOccurred())
		Expect(trans2).To(HaveLen(1))
		Expect(info2).ToNot(BeNil())

		// Verify they are semantically identical
		Expect(trans1[0].Date).To(Equal(trans2[0].Date))
		Expect(trans1[0].Movements[0].Amount).To(Equal(trans2[0].Movements[0].Amount))
		Expect(trans1[0].Description).To(ContainSubstring("Coffee Shop"))

		// Transactions now have only 1 external ID (raw record hash)
		Expect(trans1[0].ExternalIds).To(HaveLen(1))
		Expect(trans2[0].ExternalIds).To(HaveLen(1))

		// Verify hashes are DIFFERENT because it hashes the source row
		Expect(trans1[0].ExternalIds[0]).ToNot(Equal(trans2[0].ExternalIds[0]))
	})

	It("filters out zero-amount movements", func() {
		// CSV with zero amount and zero fee
		csv := `Type,Product,Started Date,Completed Date,Description,Amount,Fee,Currency,State,Balance
CARD_PAYMENT,Current,2023-01-01 10:00:00,2023-01-02 10:00:00,Zero Item,0.00,0.00,CZK,COMPLETED,1000.00`

		_, trans, err := converter.ParseTransactions(context.Background(), "csv", csv)
		Expect(err).ToNot(HaveOccurred())
		Expect(trans).To(HaveLen(1))
		// Should have 0 movements because both amount and fee are zero
		Expect(trans[0].Movements).To(BeEmpty())

		// CSV with non-zero amount but zero fee
		csv2 := `Type,Product,Started Date,Completed Date,Description,Amount,Fee,Currency,State,Balance
CARD_PAYMENT,Current,2023-01-01 10:00:00,2023-01-02 10:00:00,Normal Item,100.00,0.00,CZK,COMPLETED,1100.00`
		_, trans2, err := converter.ParseTransactions(context.Background(), "csv", csv2)
		Expect(err).ToNot(HaveOccurred())
		Expect(trans2).To(HaveLen(1))
		// Should have 2 movements (Source and Account), but NOT Fee
		Expect(trans2[0].Movements).To(HaveLen(2))
	})
})

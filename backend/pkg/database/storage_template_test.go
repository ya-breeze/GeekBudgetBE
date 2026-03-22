package database_test

import (
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestTemplateStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TemplateStorage Suite")
}

var _ = Describe("TemplateStorage", func() {
	var db database.Storage

	BeforeEach(func() {
		logger := slog.Default()
		cfg := &config.Config{DBPath: ":memory:", Verbose: false}
		db = database.NewStorage(logger, cfg)
		err := db.Open()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(db.Close)
	})

	Describe("CreateTemplate", func() {
		It("creates a template and returns it with an ID", func() {
			tpl, err := db.CreateTemplate("user1", &goserver.TransactionTemplateNoId{
				Name: "Rent",
				Movements: []goserver.Movement{
					{
						Amount:     decimal.NewFromInt(1000),
						CurrencyId: "some-currency-id",
						AccountId:  "some-account-id",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(tpl.Id).NotTo(BeEmpty())
			Expect(tpl.Name).To(Equal("Rent"))
		})
	})
})

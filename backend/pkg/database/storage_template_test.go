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
	const userID = "user1"
	const otherUserID = "user2"

	movement := goserver.Movement{
		Amount:     decimal.NewFromInt(100),
		CurrencyId: "currency-1",
		AccountId:  "account-1",
	}

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
			tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(tpl.Id).NotTo(BeEmpty())
			Expect(tpl.Name).To(Equal("Rent"))
		})

		It("stores description, place, tags, partnerName, extra", func() {
			tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:        "Groceries",
				Description: "Weekly shopping",
				Place:       "Tesco",
				Tags:        []string{"food", "weekly"},
				PartnerName: "Tesco PLC",
				Extra:       "ref:123",
				Movements:   []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(tpl.Description).To(Equal("Weekly shopping"))
			Expect(tpl.Place).To(Equal("Tesco"))
			Expect(tpl.Tags).To(ConsistOf("food", "weekly"))
			Expect(tpl.PartnerName).To(Equal("Tesco PLC"))
			Expect(tpl.Extra).To(Equal("ref:123"))
		})
	})

	Describe("GetTemplates", func() {
		BeforeEach(func() {
			_, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:      "Salary",
				Movements: []goserver.Movement{{Amount: decimal.NewFromInt(2000), CurrencyId: "currency-1", AccountId: "account-2"}},
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.CreateTemplate(otherUserID, &goserver.TransactionTemplateNoId{
				Name:      "OtherUserTemplate",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns only templates for the requesting user", func() {
			templates, err := db.GetTemplates(userID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(2))
			names := []string{templates[0].Name, templates[1].Name}
			Expect(names).To(ConsistOf("Rent", "Salary"))
		})

		It("filters by accountId when provided", func() {
			accountID := "account-1"
			templates, err := db.GetTemplates(userID, &accountID)
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(HaveLen(1))
			Expect(templates[0].Name).To(Equal("Rent"))
		})

		It("returns empty slice when no templates match the accountId filter", func() {
			accountID := "account-999"
			templates, err := db.GetTemplates(userID, &accountID)
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(BeEmpty())
		})
	})

	Describe("UpdateTemplate", func() {
		var templateID string

		BeforeEach(func() {
			tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
			templateID = tpl.Id
		})

		It("updates the template and returns the updated version", func() {
			updated, err := db.UpdateTemplate(userID, templateID, &goserver.TransactionTemplateNoId{
				Name:      "Rent Updated",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("Rent Updated"))
		})

		It("returns ErrNotFound when template does not exist", func() {
			_, err := db.UpdateTemplate(userID, "00000000-0000-0000-0000-000000000000", &goserver.TransactionTemplateNoId{
				Name:      "X",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).To(MatchError(database.ErrNotFound))
		})

		It("returns ErrNotFound when template belongs to another user", func() {
			_, err := db.UpdateTemplate(otherUserID, templateID, &goserver.TransactionTemplateNoId{
				Name:      "Hack",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).To(MatchError(database.ErrNotFound))
		})
	})

	Describe("DeleteTemplate", func() {
		var templateID string

		BeforeEach(func() {
			tpl, err := db.CreateTemplate(userID, &goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{movement},
			})
			Expect(err).NotTo(HaveOccurred())
			templateID = tpl.Id
		})

		It("deletes the template", func() {
			err := db.DeleteTemplate(userID, templateID)
			Expect(err).NotTo(HaveOccurred())

			templates, err := db.GetTemplates(userID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(templates).To(BeEmpty())
		})

		It("returns ErrNotFound for non-existent template", func() {
			err := db.DeleteTemplate(userID, "00000000-0000-0000-0000-000000000000")
			Expect(err).To(MatchError(database.ErrNotFound))
		})

		It("returns ErrNotFound when template belongs to another user", func() {
			err := db.DeleteTemplate(otherUserID, templateID)
			Expect(err).To(MatchError(database.ErrNotFound))
		})
	})
})

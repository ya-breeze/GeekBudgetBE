//nolint:fatcontext
package test_test

import (
	"context"
	"encoding/base64"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Currencies API", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
	var accessToken string
	var storage database.Storage
	logger := test.CreateTestLogger()

	BeforeEach(func() {
		forcedImportChan := make(chan common.ForcedImport)

		ctx, cancel = context.WithCancel(context.Background())
		hashed, err := auth.HashPassword([]byte(Pass1))
		if err != nil {
			panic("Error hashing password")
		}

		cfg = &config.Config{
			Port:  0,
			Users: User1 + ":" + base64.StdEncoding.EncodeToString(hashed),
		}

		storage = database.NewStorage(logger, cfg)
		if err = storage.Open(); err != nil {
			panic(err)
		}
		addr, finishCham, err = server.Serve(ctx, logger, storage, cfg, forcedImportChan)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		accessToken = getAccessToken(client, ctx)
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	It("gets empty list of existing currencies", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		currencies, _, err := client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(BeEmpty())
	})

	It("performs CRUD for currencies", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		currency := goclient.CurrencyNoID{
			Name:        "USD",
			Description: utils.StrToRef("Czech koruna"),
		}
		createdCurrency, _, err := client.CurrenciesAPI.
			CreateCurrency(ctx).
			CurrencyNoID(currency).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(createdCurrency).ToNot(BeNil())
		Expect(createdCurrency.Name).To(Equal(currency.Name))

		currencies, _, err := client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(HaveLen(1))

		currency.Name = "CZK"
		updatedCurrency, _, err := client.CurrenciesAPI.
			UpdateCurrency(ctx, createdCurrency.Id).
			CurrencyNoID(currency).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(updatedCurrency).ToNot(BeNil())
		Expect(updatedCurrency.Name).To(Equal("CZK"))

		_, err = client.CurrenciesAPI.DeleteCurrency(ctx, updatedCurrency.Id).Execute()
		Expect(err).ToNot(HaveOccurred())

		currencies, _, err = client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(BeEmpty())
	})

	It("fails to delete currency if it's in use and no replacement is provided", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		// 1. Create currency
		cur, _, err := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(goclient.CurrencyNoID{Name: "USD"}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 2. Create internal accounts for transaction
		acc1, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(goclient.AccountNoID{Name: "Acc1", Type: "asset"}).Execute()
		Expect(err).ToNot(HaveOccurred())
		acc2, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(goclient.AccountNoID{Name: "Acc2", Type: "asset"}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 3. Create transaction using this currency
		_, _, err = client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(goclient.TransactionNoID{
			Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Movements: []goclient.Movement{
				{Amount: decimal.NewFromInt(100), CurrencyId: cur.Id, AccountId: utils.StrToRef(acc1.Id)},
				{Amount: decimal.NewFromInt(-100), CurrencyId: cur.Id, AccountId: utils.StrToRef(acc2.Id)},
			},
		}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 4. Try to delete currency - should fail with 400
		_, err = client.CurrenciesAPI.DeleteCurrency(ctx, cur.Id).Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("400"))
	})

	It("successfully deletes currency when replacement is provided", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		// 1. Create two currencies
		cur1, _, err := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(goclient.CurrencyNoID{Name: "USD"}).Execute()
		Expect(err).ToNot(HaveOccurred())
		cur2, _, err := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(goclient.CurrencyNoID{Name: "EUR"}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 2. Create accounts
		acc1, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(goclient.AccountNoID{Name: "Acc1", Type: "asset"}).Execute()
		Expect(err).ToNot(HaveOccurred())
		acc2, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(goclient.AccountNoID{Name: "Acc2", Type: "asset"}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 3. Create transaction using cur1
		tr, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(goclient.TransactionNoID{
			Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Movements: []goclient.Movement{
				{Amount: decimal.NewFromInt(100), CurrencyId: cur1.Id, AccountId: utils.StrToRef(acc1.Id)},
				{Amount: decimal.NewFromInt(-100), CurrencyId: cur1.Id, AccountId: utils.StrToRef(acc2.Id)},
			},
		}).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 4. Delete cur1 with cur2 as replacement
		_, err = client.CurrenciesAPI.DeleteCurrency(ctx, cur1.Id).ReplaceWithCurrencyId(cur2.Id).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 5. Verify transaction now uses cur2
		updatedTr, _, err := client.TransactionsAPI.GetTransaction(ctx, tr.Id).Execute()
		Expect(err).ToNot(HaveOccurred())
		for _, m := range updatedTr.Movements {
			Expect(m.CurrencyId).To(Equal(cur2.Id))
		}

		// 6. Verify cur1 is gone
		currencies, _, err := client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).To(HaveLen(1))
		Expect(currencies[0].Id).To(Equal(cur2.Id))
	})
})

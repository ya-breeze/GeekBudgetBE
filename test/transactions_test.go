package test_test

import (
	"context"
	"encoding/base64"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Transactions API", func() {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		cfg         *config.Config
		addr        net.Addr
		finishCham  chan int
		client      *goclient.APIClient
		accessToken string
		accounts    []goserver.Account
		currencies  []goserver.Currency
		storage     database.Storage
	)
	logger := test.CreateTestLogger()
	now := time.Now()

	BeforeEach(func() {
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
		if err := storage.Open(); err != nil {
			panic(err)
		}
		addr, finishCham, err = server.Serve(ctx, logger, storage, cfg)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		accessToken = getAccessToken(client, ctx)

		authCtx := context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		accounts = test.PrepareAccounts()
		currencies = test.PrepareCurrencies()
		for i, account := range accounts {
			a := goclient.AccountNoID{
				Name: account.Name,
				Type: account.Type,
			}
			acc, _, err := client.AccountsAPI.CreateAccount(authCtx).AccountNoID(a).Execute()
			Expect(err).ToNot(HaveOccurred())
			accounts[i].Id = acc.Id
		}
		for i, currency := range currencies {
			c := goclient.CurrencyNoID{
				Name: currency.Name,
			}
			cur, _, err := client.CurrenciesAPI.CreateCurrency(authCtx).CurrencyNoID(c).Execute()
			Expect(err).ToNot(HaveOccurred())
			currencies[i].Id = cur.Id
		}
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	It("gets empty list of existing transactions", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		transactions, _, err := client.TransactionsAPI.GetTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(BeEmpty())
	})

	It("performs CRUD for transaction", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		t := goclient.TransactionNoID{
			Date:        now,
			Tags:        []string{"tag1", "tag2"},
			ExternalIds: []string{"ext1", "ext2"},
			Movements: []goclient.Movement{
				{
					AccountId:  &accounts[2].Id,
					CurrencyId: currencies[2].Id,
					Amount:     100,
				},
				{
					AccountId:  &accounts[0].Id,
					CurrencyId: currencies[2].Id,
					Amount:     -100,
				},
			},
		}

		// Create transaction
		created, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(created).ToNot(BeNil())
		Expect(created.Id).ToNot(BeEmpty())
		// Expect(created.Name).To(Equal(t.Name))
		// Expect(created.Type).To(Equal(t.Type))

		// Get transactions
		transactions, _, err := client.TransactionsAPI.GetTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].Id).To(Equal(created.Id))

		// Update transaction
		t.Description = utils.StrToRef("New description")
		updated, _, err := client.TransactionsAPI.UpdateTransaction(ctx, created.Id).TransactionNoID(t).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).ToNot(BeNil())
		Expect(updated.Id).To(Equal(created.Id))
		Expect(*updated.Description).To(Equal(*t.Description))

		// Get transaction by ID
		transaction, _, err := client.TransactionsAPI.GetTransaction(ctx, created.Id).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transaction).ToNot(BeNil())
		Expect(transaction.Id).To(Equal(created.Id))

		// Get aggregated expenses
		expenses, _, err := client.AggregationsAPI.GetExpenses(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(expenses).ToNot(BeNil())
		Expect(expenses.From.UnixMilli()).To(Equal(
			utils.RoundToGranularity(now, utils.GranularityMonth, false).UnixMilli()))
		Expect(expenses.To.UnixMilli()).To(Equal(
			utils.RoundToGranularity(now, utils.GranularityMonth, true).UnixMilli()))
		Expect(expenses.Intervals).To(HaveLen(1))
		Expect(expenses.Currencies).To(HaveLen(1))
		Expect(expenses.Currencies[0].CurrencyId).To(Equal(currencies[2].Id))
		Expect(expenses.Currencies[0].Accounts).To(HaveLen(1))

		Expect(expenses.Currencies[0].Accounts[0].AccountId).To(Equal(accounts[2].Id))
		Expect(expenses.Currencies[0].Accounts[0].Amounts).To(HaveLen(1))
		Expect(expenses.Currencies[0].Accounts[0].Amounts[0]).To(Equal(100.0))

		// Delete transaction
		_, err = client.TransactionsAPI.DeleteTransaction(ctx, created.Id).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Get transactions
		transactions, _, err = client.TransactionsAPI.GetTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(BeEmpty())
	})
})

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
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Unprocessed Transactions API", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
	var accessToken string
	logger := test.CreateTestLogger()

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

		addr, finishCham, err = server.Serve(ctx, logger, cfg)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		accessToken = getAccessToken(client, ctx)
	})

	AfterEach(func() {
		cancel()
		<-finishCham
	})

	It("gets empty list of existing unprocessed transactions", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		transactions, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(BeEmpty())
	})

	It("converts transaction with empty account to unprocessed", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		t := goclient.TransactionNoID{
			Date:        time.Now(),
			Description: utils.StrToRef("Purchase in BILLA"),
			Tags:        []string{"tag1", "tag2"},
			ExternalIds: []string{"ext1", "ext2"},
			Movements: []goclient.Movement{
				{
					AccountId:  nil,
					CurrencyId: "currencyID",
					Amount:     100,
				},
				{
					AccountId:  utils.StrToRef("accountID"),
					CurrencyId: "currencyID",
					Amount:     -100,
				},
			},
		}

		// Create transaction
		created, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(created).ToNot(BeNil())
		Expect(created.Id).ToNot(BeEmpty())

		// Create matcher
		m := goclient.MatcherNoID{
			Name:              "matcher1",
			OutputDescription: "Billa",
			OutputAccountId:   "accountID",
			DescriptionRegExp: utils.StrToRef(`(?i)\bBilla\b`),
		}
		_, _, err = client.MatchersAPI.CreateMatcher(ctx).MatcherNoID(m).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Get unprocessed transactions
		transactions, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].Transaction.Id).To(Equal(created.Id))
		Expect(transactions[0].Matched).ToNot(BeEmpty())
		Expect(*transactions[0].Matched[0].Transaction.Description).To(Equal(m.OutputDescription))
		Expect(*transactions[0].Matched[0].Transaction.Movements[0].AccountId).To(Equal(m.OutputAccountId))
		Expect(*transactions[0].Matched[0].Transaction.Movements[1].AccountId).To(Equal(*t.Movements[1].AccountId))

		// // Update transaction
		// t.Description = utils.StrToRef("New description")
		// updated, _, err := client.TransactionsAPI.UpdateTransaction(ctx, created.Id).TransactionNoID(t).Execute()
		// Expect(err).ToNot(HaveOccurred())
		// Expect(updated).ToNot(BeNil())
		// Expect(updated.Id).To(Equal(created.Id))
		// Expect(*updated.Description).To(Equal(*t.Description))

		// // Get transaction by ID
		// transaction, _, err := client.TransactionsAPI.GetTransaction(ctx, created.Id).Execute()
		// Expect(err).ToNot(HaveOccurred())
		// Expect(transaction).ToNot(BeNil())
		// Expect(transaction.Id).To(Equal(created.Id))

		// // Delete transaction
		// _, err = client.TransactionsAPI.DeleteTransaction(ctx, created.Id).Execute()
		// Expect(err).ToNot(HaveOccurred())

		// // Get transactions
		// transactions, _, err = client.TransactionsAPI.GetTransactions(ctx).Execute()
		// Expect(err).ToNot(HaveOccurred())
		// Expect(transactions).ToNot(BeNil())
		// Expect(transactions).To(BeEmpty())
	})
})

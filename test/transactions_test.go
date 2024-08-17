package test_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

var _ = Describe("Transactions API", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
	var accessToken string
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
			Date:        time.Now(),
			Tags:        []string{"tag1", "tag2"},
			ExternalIds: []string{"ext1", "ext2"},
			Movements: []goclient.Movement{
				{
					AccountId:  "account1",
					CurrencyId: "currency1",
					Amount:     100,
				},
				{
					AccountId:  "account2",
					CurrencyId: "currency1",
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

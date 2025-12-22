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
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Duplicate Unprocessed Transactions REPRO", func() {
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

	It("does NOT show unprocessed transactions as duplicates of each other", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		now := time.Now().Truncate(time.Second)

		t1 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Identical Unprocessed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // undefined
					CurrencyId: "CZK",
					Amount:     100,
				},
				{
					AccountId:  utils.StrToRef("valid-account"),
					CurrencyId: "CZK",
					Amount:     -100,
				},
			},
		}

		t2 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Identical Unprocessed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // undefined
					CurrencyId: "CZK",
					Amount:     100,
				},
				{
					AccountId:  utils.StrToRef("valid-account"),
					CurrencyId: "CZK",
					Amount:     -100,
				},
			},
		}

		// Create two identical unprocessed transactions
		_, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t1).Execute()
		Expect(err).ToNot(HaveOccurred())
		_, _, err = client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t2).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Get unprocessed transactions
		transactions, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).To(HaveLen(2))

		// Check duplicates - FAIL EXPECTED HERE before fix
		for _, ut := range transactions {
			Expect(ut.Duplicates).To(BeEmpty(), "Unprocessed transaction should not have another unprocessed transaction as a duplicate")
		}
	})

	It("shows processed transactions as duplicates", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		now := time.Now().Truncate(time.Second)

		// Create a PROCESSED transaction (all accounts defined)
		t1 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Processed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  utils.StrToRef("account-A"),
					CurrencyId: "CZK",
					Amount:     100,
				},
				{
					AccountId:  utils.StrToRef("account-B"),
					CurrencyId: "CZK",
					Amount:     -100,
				},
			},
		}
		processed, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t1).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Create an UNPROCESSED transaction that looks like a duplicate of the processed one
		t2 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Unprocessed Transaction Candidate"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // undefined
					CurrencyId: "CZK",
					Amount:     100,
				},
				{
					AccountId:  utils.StrToRef("account-B"),
					CurrencyId: "CZK",
					Amount:     -100,
				},
			},
		}
		_, _, err = client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(t2).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Get unprocessed transactions
		transactions, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).To(HaveLen(1))

		// Check duplicates - Should find the processed transaction
		Expect(transactions[0].Duplicates).To(HaveLen(1))
		Expect(transactions[0].Duplicates[0].Id).To(Equal(processed.Id))
	})
})

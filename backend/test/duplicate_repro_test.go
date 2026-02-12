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

		// Create CZK currency
		curReq := goclient.CurrencyNoID{Name: "CZK"}
		cur, _, _ := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(curReq).Execute()

		// Create valid-account
		accReq := goclient.AccountNoID{Name: "Valid Account", Type: "asset"}
		acc, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(accReq).Execute()

		now := time.Now().Truncate(time.Second)

		t1 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Identical Unprocessed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // undefined
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(100),
				},
				{
					AccountId:  &acc.Id,
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(-100),
				},
			},
		}

		t2 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Identical Unprocessed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // undefined
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(100),
				},
				{
					AccountId:  &acc.Id,
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(-100),
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

		// Create CZK currency
		curReq := goclient.CurrencyNoID{Name: "CZK"}
		cur, _, _ := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(curReq).Execute()

		// Create accounts
		accAReq := goclient.AccountNoID{Name: "Account A", Type: "asset"}
		accA, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(accAReq).Execute()
		accBReq := goclient.AccountNoID{Name: "Account B", Type: "asset"}
		accB, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(accBReq).Execute()

		now := time.Now().Truncate(time.Second)

		// Create a PROCESSED transaction (all accounts defined)
		t1 := goclient.TransactionNoID{
			Date:        now,
			Description: utils.StrToRef("Processed Transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  &accA.Id,
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(100),
				},
				{
					AccountId:  &accB.Id,
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(-100),
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
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(100),
				},
				{
					AccountId:  &accB.Id,
					CurrencyId: cur.Id,
					Amount:     decimal.NewFromInt(-100),
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

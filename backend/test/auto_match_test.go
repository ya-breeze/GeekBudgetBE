//nolint:intrange,dupl
package test_test

import (
	"context"
	"encoding/base64"
	"fmt"
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
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Auto Match API", func() {
	var (
		ctx         context.Context
		cancel      context.CancelFunc
		cfg         *config.Config
		addr        net.Addr
		finishChan  chan int
		client      *goclient.APIClient
		accessToken string
		storage     database.Storage
		account     *goclient.Account
		currency    *goclient.Currency
	)
	logger := test.CreateTestLogger()
	now := time.Now()

	BeforeEach(func() {
		forcedImportChan := make(chan common.ForcedImport)

		ctx, cancel = context.WithCancel(context.Background())
		hashed, err := auth.HashPassword([]byte(Pass1))
		if err != nil {
			panic("Error hashing password")
		}

		cfg = &config.Config{
			Port:                          0,
			Users:                         User1 + ":" + base64.StdEncoding.EncodeToString(hashed),
			MatcherConfirmationHistoryMax: 20,
		}

		storage = database.NewStorage(logger, cfg)
		if err = storage.Open(); err != nil {
			panic(err)
		}
		addr, finishChan, err = server.Serve(ctx, logger, storage, cfg, forcedImportChan)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		accessToken = getAccessToken(client, ctx)
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		// Create test account
		acc := goclient.AccountNoID{
			Name: "Test Account",
			Type: "asset",
		}
		account, _, err = client.AccountsAPI.CreateAccount(ctx).AccountNoID(acc).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Create test currency
		curr := goclient.CurrencyNoID{
			Name: "US Dollar",
		}
		currency, _, err = client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(curr).Execute()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
		<-finishChan
		storage.Close()
	})

	It("auto-matches transactions when matcher becomes perfect", func() {
		// 1. Create 12 identical unprocessed transactions
		transactions := make([]*goclient.Transaction, 12)
		for i := 0; i < 12; i++ {
			txn := goclient.TransactionNoID{
				Date:        now.Add(time.Duration(i) * time.Minute),
				Description: ptrString(fmt.Sprintf("AutoMatch Test %d", i)),
				PartnerName: ptrString("Netflix"), // Common pattern
				Movements: []goclient.Movement{
					{
						Amount:     decimal.NewFromInt(int64(100 + i)),
						CurrencyId: currency.Id,
						// AccountId is intentionally missing to make it unprocessed
					},
				},
			}
			// Note: CreateTransaction with missing AccountId creates an unprocessed transaction
			created, _, err := client.TransactionsAPI.CreateTransaction(ctx).
				TransactionNoID(txn).Execute()
			Expect(err).ToNot(HaveOccurred())
			transactions[i] = created
		}

		// 2. Create a matcher that targets these transactions
		matcherNoID := goclient.MatcherNoID{
			OutputDescription: ptrString("Netflix Subscription"),
			PartnerNameRegExp: ptrString("Netflix"),
			OutputAccountId:   account.Id,
			OutputTags:        []string{"subscription"},
		}
		matcher, _, err := client.MatchersAPI.CreateMatcher(ctx).MatcherNoID(matcherNoID).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 3. Confirm 10 transactions sequentially
		for i := 0; i < 10; i++ {
			// Get the transaction to be converted
			// ConvertUnprocessedTransaction requires TransactionNoId body, which is basically the new state
			// We populate it with valid data including the matcherID to signal confirmation

			// Assuming ConvertUnprocessedTransaction API takes the new state of the transaction
			// We need to fetch the unprocessed transaction first to get its current state?
			// Or just construct what we want. The API takes TransactionNoId.

			// We need to provide AccountId to make it processed.

			txnConvert := goclient.TransactionNoID{
				Date:        transactions[i].Date,
				Description: matcher.OutputDescription, // Apply matcher output
				Movements: []goclient.Movement{
					{
						Amount:     transactions[i].Movements[0].Amount,
						CurrencyId: transactions[i].Movements[0].CurrencyId,
						AccountId:  &account.Id, // Apply matcher output
					},
				},
			}

			// Call convert API
			resp, httpResp, err := client.UnprocessedTransactionsAPI.ConvertUnprocessedTransaction(ctx, transactions[i].Id).
				TransactionNoID(txnConvert).
				MatcherId(matcher.Id).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(200))

			if i < 9 {
				// Matched < 10 times, so no auto-processing yet
				Expect(resp.GetAutoProcessedIds()).To(BeEmpty())
			} else {
				// 10th confirmation: should trigger auto-processing for the remaining 2
				Expect(resp.GetAutoProcessedIds()).To(HaveLen(2))
				Expect(resp.GetAutoProcessedIds()).To(ContainElements(transactions[10].Id, transactions[11].Id))
			}
		}

		// 4. Verify the last 2 transactions are now processed
		// We can try to fetch them via GetUnprocessedTransactions -> should not be there
		// Or GetTransaction -> check fields

		for i := 10; i < 12; i++ {
			t, _, err := client.TransactionsAPI.GetTransaction(ctx, transactions[i].Id).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(t.Description).To(Equal(ptrString("Netflix Subscription")))
			Expect(t.GetMovements()[0].AccountId).To(Equal(ptrString(account.Id)))
			Expect(t.Tags).To(ContainElement("subscription"))
			Expect(t.MatcherId).To(Equal(ptrString(matcher.Id)))
			Expect(*t.IsAuto).To(BeTrue())
		}
	})
})

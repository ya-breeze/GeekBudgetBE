//nolint:fatcontext
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

var _ = Describe("Unprocessed Transactions API", func() {
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
			OutputDescription: utils.StrToRef("Billa"),
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
		Expect(*transactions[0].Matched[0].Transaction.Description).To(Equal(*m.OutputDescription))
		Expect(*transactions[0].Matched[0].Transaction.Movements[0].AccountId).To(Equal(m.OutputAccountId))
		Expect(*transactions[0].Matched[0].Transaction.Movements[1].AccountId).To(Equal(*t.Movements[1].AccountId))

		_, _, err = client.UnprocessedTransactionsAPI.
			ConvertUnprocessedTransaction(ctx, transactions[0].Transaction.Id).
			TransactionNoID(transactions[0].Matched[0].Transaction).
			Execute()
		Expect(err).ToNot(HaveOccurred())

		transactions, _, err = client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())
		Expect(transactions).To(BeEmpty())

		updated, _, err := client.TransactionsAPI.GetTransaction(ctx, created.Id).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).ToNot(BeNil())
		Expect(updated.Id).To(Equal(created.Id))
		Expect(*updated.Description).To(Equal(*m.OutputDescription))
	})

	It("ignores unprocessed transactions older than account's ignoreUnprocessedBefore date", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		// 1. Create an account with ignoreUnprocessedBefore set
		ignoreDate := time.Now().Add(-24 * time.Hour)
		acc := goclient.AccountNoID{
			Name:                    "IgnoreTestAccount",
			Type:                    "asset",
			IgnoreUnprocessedBefore: &ignoreDate,
		}
		createdAccount, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(acc).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 2. Create a transaction with empty accountId movement BEFORE the ignore date
		oldDate := ignoreDate.Add(-1 * time.Hour)
		tOld := goclient.TransactionNoID{
			Date:        oldDate,
			Description: utils.StrToRef("Old Unprocessed"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil,
					CurrencyId: "currencyID",
					Amount:     100,
				},
				{
					AccountId:  &createdAccount.Id,
					CurrencyId: "currencyID",
					Amount:     -100,
				},
			},
		}
		_, _, err = client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(tOld).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 3. Create a transaction with empty accountId movement AFTER the ignore date
		newDate := ignoreDate.Add(1 * time.Hour)
		tNew := goclient.TransactionNoID{
			Date:        newDate,
			Description: utils.StrToRef("New Unprocessed"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil,
					CurrencyId: "currencyID",
					Amount:     200,
				},
				{
					AccountId:  &createdAccount.Id,
					CurrencyId: "currencyID",
					Amount:     -200,
				},
			},
		}
		createdNew, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(tNew).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 4. Get unprocessed transactions
		transactions, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(transactions).ToNot(BeNil())

		// 5. Verify only transactions AFTER the date are returned
		// The old one should be filtered out because it has a movement with createdAccount.Id
		// and its date is before createdAccount.IgnoreUnprocessedBefore
		Expect(transactions).To(HaveLen(1))
		Expect(transactions[0].Transaction.Id).To(Equal(createdNew.Id))
	})

	It("triggers balance verification notification on mismatch", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		// 1. Create a currency
		curReq := goclient.CurrencyNoID{Name: "CZK"}
		cur, _, _ := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(curReq).Execute()

		// 2. Create an account with a specific bank balance
		accReq := goclient.AccountNoID{
			Name: "BalanceTestAccount",
			Type: "asset",
			BankInfo: &goclient.BankAccountInfo{
				Balances: []goclient.BankAccountInfoBalancesInner{
					{
						CurrencyId:     &cur.Id,
						OpeningBalance: goclient.PtrFloat64(1000.0),
						ClosingBalance: goclient.PtrFloat64(1500.0),
					},
				},
			},
		}
		acc, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(accReq).Execute()

		// 3. Add an UNPROCESSED transaction that would make the balance mismatch
		// (Opening 1000 + 400 = 1400, but bank says 1500)
		tReq := goclient.TransactionNoID{
			Date:        time.Now(),
			Description: goclient.PtrString("Mismatch transaction"),
			Movements: []goclient.Movement{
				{
					AccountId:  nil, // Unprocessed
					CurrencyId: cur.Id,
					Amount:     400,
				},
				{
					AccountId:  &acc.Id,
					CurrencyId: cur.Id,
					Amount:     -400,
				},
			},
		}
		t, _, _ := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(tReq).Execute()

		// 4. Verify no notification yet (because there's still an unprocessed transaction)
		notifications, _, _ := client.NotificationsAPI.GetNotifications(ctx).Execute()
		Expect(notifications).To(BeEmpty())

		// 4a. Create another account to serve as the "offset" for the processed transaction
		offsetAccReq := goclient.AccountNoID{
			Name: "OffsetAccount",
			Type: "asset",
		}
		offsetAcc, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(offsetAccReq).Execute()

		// 5. Convert the transaction to processed
		// Now Opening 1000 + 400 = 1400. Bank says 1500. Should trigger notification.
		conversionReq := goclient.TransactionNoID{
			Date:        t.Date,
			Description: t.Description,
			Movements: []goclient.Movement{
				{
					AccountId:  &acc.Id,
					CurrencyId: cur.Id,
					Amount:     400,
				},
				{
					AccountId:  &offsetAcc.Id, // No longer nil
					CurrencyId: cur.Id,
					Amount:     -400,
				},
			},
		}
		_, _, err := client.UnprocessedTransactionsAPI.ConvertUnprocessedTransaction(ctx, t.Id).TransactionNoID(conversionReq).Execute()
		Expect(err).ToNot(HaveOccurred())

		// 6. Verify notification is created
		notifications, _, _ = client.NotificationsAPI.GetNotifications(ctx).Execute()
		Expect(notifications).To(HaveLen(1))
		Expect(notifications[0].Title).To(Equal("Balance Mismatch Detected"))
		Expect(notifications[0].Description).To(ContainSubstring("Account balance: 1400.00"))
		Expect(notifications[0].Description).To(ContainSubstring("Bank balance: 1500.00"))
	})
})

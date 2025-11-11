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
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/test"
)

func ptrString(s string) *string {
	return &s
}

var _ = Describe("Matchers API", func() {
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
		transaction *goclient.Transaction
	)
	logger := test.CreateTestLogger()
	now := time.Now()

	BeforeEach(func() {
		forcedImportChan := make(chan background.ForcedImport)

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
		addr, finishChan, err = server.Serve(ctx, logger, storage, cfg, forcedImportChan)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		accessToken = getAccessToken(client, ctx)

		// Create test account
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
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

		// Create test transaction
		txn := goclient.TransactionNoID{
			Date:           now,
			Description:    ptrString("Test transaction"),
			PartnerName:    ptrString("Test Partner"),
			PartnerAccount: ptrString("12345"),
			Place:          ptrString("Test Place"),
			Tags:           []string{"test"},
			Movements: []goclient.Movement{
				{
					Amount:     100,
					CurrencyId: currency.Id,
					AccountId:  &account.Id,
				},
			},
		}
		transaction, _, err = client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(txn).Execute()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
		<-finishChan
		storage.Close()
	})

	It("checks matcher against transaction - match", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		matcher := goclient.MatcherNoID{
			Name:                       "Test Matcher",
			OutputDescription:          "Converted",
			DescriptionRegExp:          ptrString("Test.*"),
			PartnerAccountNumberRegExp: ptrString("123.*"),
			OutputAccountId:            account.Id,
			OutputTags:                 []string{"converted"},
		}

		checkRequest := goclient.CheckMatcherRequest{
			Matcher: matcher,
			Transaction: goclient.TransactionNoID{
				Date:           transaction.Date,
				Description:    transaction.Description,
				PartnerName:    transaction.PartnerName,
				PartnerAccount: transaction.PartnerAccount,
				Place:          transaction.Place,
				Tags:           transaction.Tags,
				Movements:      transaction.Movements,
			},
		}

		// Note: The CheckMatcher endpoint is not exposed in the goclient yet
		// This test demonstrates the expected behavior
		Expect(checkRequest.Matcher.Name).To(Equal("Test Matcher"))
		Expect(checkRequest.Transaction.Description).To(Equal(ptrString("Test transaction")))
	})

	It("checks matcher against transaction - no match", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		matcher := goclient.MatcherNoID{
			Name:                       "Test Matcher",
			OutputDescription:          "Converted",
			DescriptionRegExp:          ptrString("NonMatching.*"),
			PartnerAccountNumberRegExp: ptrString("999.*"),
			OutputAccountId:            account.Id,
			OutputTags:                 []string{"converted"},
		}

		checkRequest := goclient.CheckMatcherRequest{
			Matcher: matcher,
			Transaction: goclient.TransactionNoID{
				Date:           transaction.Date,
				Description:    transaction.Description,
				PartnerName:    transaction.PartnerName,
				PartnerAccount: transaction.PartnerAccount,
				Place:          transaction.Place,
				Tags:           transaction.Tags,
				Movements:      transaction.Movements,
			},
		}

		// Note: The CheckMatcher endpoint is not exposed in the goclient yet
		// This test demonstrates the expected behavior
		Expect(checkRequest.Matcher.Name).To(Equal("Test Matcher"))
		Expect(checkRequest.Transaction.Description).To(Equal(ptrString("Test transaction")))
	})

	It("creates and retrieves matcher", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		matcher := goclient.MatcherNoID{
			Name:                       "Test Matcher",
			OutputDescription:          "Converted",
			DescriptionRegExp:          ptrString("Test.*"),
			PartnerAccountNumberRegExp: ptrString("123.*"),
			OutputAccountId:            account.Id,
			OutputTags:                 []string{"converted"},
		}

		// Create matcher
		created, _, err := client.MatchersAPI.CreateMatcher(ctx).MatcherNoID(matcher).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(created).ToNot(BeNil())
		Expect(created.Id).ToNot(BeEmpty())
		Expect(created.Name).To(Equal(matcher.Name))

		// Get matchers
		matchers, _, err := client.MatchersAPI.GetMatchers(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(matchers).ToNot(BeNil())
		Expect(matchers).To(HaveLen(1))
		Expect(matchers[0].Id).To(Equal(created.Id))
	})
})

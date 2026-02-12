//nolint:fatcontext
package test_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

var _ = Describe("Web unprocessed convert integration", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
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
			Port:                          0,
			Users:                         User1 + ":" + base64.StdEncoding.EncodeToString(hashed),
			CookieName:                    "geekbudgetcookie",
			MatcherConfirmationHistoryMax: 10,
		}

		storage = database.NewStorage(logger, cfg)
		Expect(storage.Open()).To(Succeed())

		addr, finishCham, err = server.Serve(ctx, logger, storage, cfg, forcedImportChan)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)

		// authenticate API client for subsequent API calls
		accessToken := getAccessToken(client, ctx)
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	It("converting matched suggestion from web updates matcher confirmation history", func() {
		// Create currency
		curReq := goclient.CurrencyNoID{Name: "CZK"}
		cur, _, _ := client.CurrenciesAPI.CreateCurrency(ctx).CurrencyNoID(curReq).Execute()

		// Create account
		accReq := goclient.AccountNoID{Name: "TestAccount", Type: "asset"}
		acc, _, _ := client.AccountsAPI.CreateAccount(ctx).AccountNoID(accReq).Execute()

		// create transaction that will be matched
		txn := goclient.TransactionNoID{
			Date:        time.Now(),
			Description: utils.StrToRef("Purchase at WEBSTORE"),
			Tags:        []string{"tag1"},
			Movements: []goclient.Movement{
				{AccountId: nil, CurrencyId: cur.Id, Amount: decimal.NewFromInt(100)},
				{AccountId: &acc.Id, CurrencyId: cur.Id, Amount: decimal.NewFromInt(-100)},
			},
		}

		_, _, err := client.TransactionsAPI.CreateTransaction(ctx).TransactionNoID(txn).Execute()
		Expect(err).ToNot(HaveOccurred())

		// create matcher that matches description
		m := goclient.MatcherNoID{
			OutputDescription: utils.StrToRef("WEBSTORE"),
			OutputAccountId:   acc.Id,
			DescriptionRegExp: utils.StrToRef(`(?i)webstore`),
		}
		createdMatcher, _, err := client.MatchersAPI.CreateMatcher(ctx).MatcherNoID(m).Execute()
		Expect(err).ToNot(HaveOccurred())

		// get unprocessed and ensure matched
		unproc, _, err := client.UnprocessedTransactionsAPI.GetUnprocessedTransactions(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(unproc).ToNot(BeEmpty())
		u := unproc[0]
		Expect(u.Matched).ToNot(BeEmpty())

		// login via web to obtain session cookie
		jar, _ := cookiejar.New(nil)
		httpClient := &http.Client{Jar: jar}
		loginValues := url.Values{}
		loginValues.Set("username", User1)
		loginValues.Set("password", Pass1)
		loginResp, err := httpClient.PostForm("http://"+addr.String()+"/web/login", loginValues)
		Expect(err).ToNot(HaveOccurred())
		if _, readErr := io.ReadAll(loginResp.Body); readErr != nil {
			loginResp.Body.Close()
			Fail("failed reading login response body: " + readErr.Error())
		}
		loginResp.Body.Close()

		// prepare form for convert: use matched suggestion (first)
		matched := u.Matched[0]
		form := url.Values{}
		form.Set("transaction_id", u.Transaction.Id)
		form.Set("matcher_id", matched.MatcherId)
		// other_matchers empty
		// set movement accounts from matched transaction
		for i, mv := range matched.Transaction.Movements {
			if mv.AccountId != nil {
				form.Set(fmt.Sprintf("account_%d", i), *mv.AccountId)
			} else {
				form.Set(fmt.Sprintf("account_%d", i), "")
			}
		}

		// submit convert via web
		convertResp, err := httpClient.PostForm("http://"+addr.String()+"/web/unprocessed/convert", form)
		Expect(err).ToNot(HaveOccurred())
		if _, readErr := io.ReadAll(convertResp.Body); readErr != nil {
			convertResp.Body.Close()
			Fail("failed reading convert response body: " + readErr.Error())
		}
		convertResp.Body.Close()

		// load matcher and verify confirmation history updated (last entry true)
		userInfo, _, err := client.UserAPI.GetUser(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		userID := userInfo.Id

		// Sometimes the API does not return user id in the payload for tests; fall back to storage lookup
		if userID == "" {
			userID, err = storage.GetUserID(User1)
			Expect(err).ToNot(HaveOccurred())
		}

		loaded, err := storage.GetMatcher(userID, createdMatcher.Id)
		Expect(err).ToNot(HaveOccurred())
		Expect(loaded.ConfirmationHistory).ToNot(BeEmpty())
		Expect(loaded.ConfirmationHistory[len(loaded.ConfirmationHistory)-1]).To(BeTrue())
	})
})

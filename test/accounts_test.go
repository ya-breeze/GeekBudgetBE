package test_test

import (
	"context"
	"encoding/base64"
	"net"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Accounts API", func() {
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
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	It("gets empty list of existing accounts", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		accounts, _, err := client.AccountsAPI.GetAccounts(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).ToNot(BeNil())
		Expect(accounts).To(BeEmpty())
	})

	It("performs CRUD for account", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		acc := goclient.AccountNoID{
			Name: "Cash",
			Type: "asset",
		}

		// Create account
		created, _, err := client.AccountsAPI.CreateAccount(ctx).AccountNoID(acc).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(created).ToNot(BeNil())
		Expect(created.Id).ToNot(BeEmpty())
		Expect(created.Name).To(Equal(acc.Name))
		Expect(created.Type).To(Equal(acc.Type))

		// Get accounts
		accounts, _, err := client.AccountsAPI.GetAccounts(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).ToNot(BeNil())
		Expect(accounts).To(HaveLen(1))
		Expect(accounts[0].Id).To(Equal(created.Id))

		// Update account
		acc.Name = "Bank"
		updated, _, err := client.AccountsAPI.UpdateAccount(ctx, created.Id).AccountNoID(acc).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(updated).ToNot(BeNil())
		Expect(updated.Id).To(Equal(created.Id))
		Expect(updated.Name).To(Equal(acc.Name))

		// Get account by ID
		account, _, err := client.AccountsAPI.GetAccount(ctx, created.Id).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(account).ToNot(BeNil())
		Expect(account.Id).To(Equal(created.Id))

		// Delete account
		_, err = client.AccountsAPI.DeleteAccount(ctx, created.Id).Execute()
		Expect(err).ToNot(HaveOccurred())

		// Get accounts
		accounts, _, err = client.AccountsAPI.GetAccounts(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).ToNot(BeNil())
		Expect(accounts).To(BeEmpty())
	})
})

package test_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
)

var _ = Describe("Accounts API", func() {
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

	It("get list of existing accounts", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		accounts, _, err := client.AccountsAPI.GetAccounts(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(accounts).ToNot(BeNil())
		Expect(accounts).To(BeEmpty())
	})
})

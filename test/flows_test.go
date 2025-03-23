//nolint:fatcontext
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
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goclient"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
)

var _ = Describe("Flows", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var storage database.Storage
	// var client *goclient.APIClient
	// var accessToken string
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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
		addr, finishCham, err = server.Serve(ctx, logger, storage, cfg, forcedImportChan)
		Expect(err).ToNot(HaveOccurred())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		// client = goclient.NewAPIClient(clientCfg)

		// accessToken = getAccessToken(client, ctx)
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	// It("gets empty list of existing accounts", func() {
	// 	ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
	// 	accounts, _, err := client.AccountsAPI.GetAccounts(ctx).Execute()
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(accounts).ToNot(BeNil())
	// 	Expect(accounts).To(BeEmpty())
	// })
})

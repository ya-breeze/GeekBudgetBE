//nolint:fatcontext
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
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/test"
)

const (
	User1 = "user1"
	Pass1 = "password1"
)

var _ = Describe("User API", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
	var storage database.Storage
	logger := test.CreateTestLogger()

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
		client = goclient.NewAPIClient(clientCfg)
	})

	AfterEach(func() {
		cancel()
		<-finishCham
		storage.Close()
	})

	It("authenticates client with valid credentials", func() {
		getAccessToken(client, ctx)
	})

	It("does not authenticate client with invalid credentials", func() {
		req := client.AuthAPI.Authorize(ctx).AuthData(goclient.AuthData{
			Email:    User1,
			Password: "wrongpassword",
		})
		resp, httpResp, err := req.Execute()
		Expect(err).To(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(401))
		Expect(resp).To(BeNil())
	})

	It("returns known user object", func() {
		accessToken := getAccessToken(client, ctx)
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		user, httpResp, err := client.UserAPI.GetUser(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(user).ToNot(BeNil())
		Expect(user.Email).To(Equal(User1))
	})
})

func getAccessToken(client *goclient.APIClient, ctx context.Context) string {
	req := client.AuthAPI.Authorize(ctx).AuthData(goclient.AuthData{
		Email:    User1,
		Password: Pass1,
	})
	resp, httpResp, err := req.Execute()
	Expect(err).ToNot(HaveOccurred())
	Expect(httpResp).ToNot(BeNil())
	Expect(httpResp.StatusCode).To(Equal(200))
	Expect(resp).ToNot(BeNil())
	Expect(resp.Token).ToNot(BeEmpty())

	return resp.Token
}

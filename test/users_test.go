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

const (
	User1 = "user1"
	Pass1 = "password1"
)

var _ = Describe("GB", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient
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
	})

	AfterEach(func() {
		cancel()
		<-finishCham
	})

	It("authenticates client with valid credentials", func() {
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
})

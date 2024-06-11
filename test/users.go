package test

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
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

const (
	User1 = "user1"
	Pass1 = "password1"
)

var _ = Describe("GB", func() {
	ctx, cancel := context.WithCancel(context.Background())
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var cfg *config.Config
	var addr net.Addr
	var finishCham chan int
	var client *goclient.APIClient

	BeforeEach(func() {
		hashed, err := auth.HashPassword([]byte(Pass1))
		if err != nil {
			panic("Error hashing password")
		}

		cfg = &config.Config{
			Port:  0,
			Users: User1 + ":" + base64.StdEncoding.EncodeToString(hashed),
		}

		addr, finishCham, err = goserver.Serve(ctx, logger, cfg)
		Expect(err).To(BeNil())

		clientCfg := goclient.NewConfiguration()
		clientCfg.Servers[0].URL = "http://" + addr.String()
		client = goclient.NewAPIClient(clientCfg)
	})

	It("client can authenticate", func() {
		req := client.AuthAPI.Authorize(ctx).AuthData(goclient.AuthData{
			Email:    User1,
			Password: Pass1,
		})
		resp, httpResp, err := req.Execute()
		Expect(err).To(BeNil())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(resp).ToNot(BeNil())
		Expect(resp.Token).ToNot(BeEmpty())
	})

	AfterEach(func() {
		cancel()
		<-finishCham
	})
})

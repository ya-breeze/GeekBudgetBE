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
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Currencies API", func() {
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
		if err = storage.Open(); err != nil {
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

	It("gets empty list of existing currencies", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		currencies, _, err := client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(BeEmpty())
	})

	It("performs CRUD for currencies", func() {
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)
		currency := goclient.CurrencyNoID{
			Name:        "USD",
			Description: utils.StrToRef("Czech koruna"),
		}
		createdCurrency, _, err := client.CurrenciesAPI.
			CreateCurrency(ctx).
			CurrencyNoID(currency).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(createdCurrency).ToNot(BeNil())
		Expect(createdCurrency.Name).To(Equal(currency.Name))

		currencies, _, err := client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(HaveLen(1))

		currency.Name = "CZK"
		updatedCurrency, _, err := client.CurrenciesAPI.
			UpdateCurrency(ctx, createdCurrency.Id).
			CurrencyNoID(currency).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(updatedCurrency).ToNot(BeNil())
		Expect(updatedCurrency.Name).To(Equal("CZK"))

		_, err = client.CurrenciesAPI.DeleteCurrency(ctx, updatedCurrency.Id).Execute()
		Expect(err).ToNot(HaveOccurred())

		currencies, _, err = client.CurrenciesAPI.GetCurrencies(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(currencies).ToNot(BeNil())
		Expect(currencies).To(BeEmpty())
	})
})

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
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
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
		if httpResp != nil {
			defer func() { _ = httpResp.Body.Close() }()
		}
		Expect(err).To(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(401))
		Expect(resp).To(BeNil())
	})

	It("does not authenticate client with unknown user", func() {
		req := client.AuthAPI.Authorize(ctx).AuthData(goclient.AuthData{
			Email:    "unknown@example.com",
			Password: "password",
		})
		resp, httpResp, err := req.Execute()
		if httpResp != nil {
			defer func() { _ = httpResp.Body.Close() }()
		}
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
	It("updates user's favorite currency when field is present and non-empty", func() {
		accessToken := getAccessToken(client, ctx)
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		currencyReq := goclient.CurrencyNoID{
			Name: "USD",
		}
		createdCurrency, httpResp, err := client.CurrenciesAPI.
			CreateCurrency(ctx).
			CurrencyNoID(currencyReq).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(createdCurrency).ToNot(BeNil())

		favoriteID := createdCurrency.Id

		patch := goclient.UserPatchBody{}
		patch.SetFavoriteCurrencyId(favoriteID)

		user, httpResp, err := client.UserAPI.UpdateUserFavoriteCurrency(ctx).
			UserPatchBody(patch).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(user).ToNot(BeNil())
		Expect(user.FavoriteCurrencyId).ToNot(BeNil())
		Expect(*user.FavoriteCurrencyId).To(Equal(favoriteID))

		// verify via GetUser as well
		user2, httpResp, err := client.UserAPI.GetUser(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(user2).ToNot(BeNil())
		Expect(user2.FavoriteCurrencyId).ToNot(BeNil())
		Expect(*user2.FavoriteCurrencyId).To(Equal(favoriteID))
	})

	It("resets user's favorite currency when field is empty string", func() {
		accessToken := getAccessToken(client, ctx)
		ctx = context.WithValue(ctx, goclient.ContextAccessToken, accessToken)

		currencyReq := goclient.CurrencyNoID{
			Name: "CZK",
		}
		createdCurrency, httpResp, err := client.CurrenciesAPI.
			CreateCurrency(ctx).
			CurrencyNoID(currencyReq).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(createdCurrency).ToNot(BeNil())

		favoriteID := createdCurrency.Id

		// First set the favorite currency
		patchSet := goclient.UserPatchBody{}
		patchSet.SetFavoriteCurrencyId(favoriteID)

		_, httpResp, err = client.UserAPI.UpdateUserFavoriteCurrency(ctx).
			UserPatchBody(patchSet).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))

		// Now reset by sending empty string
		patchReset := goclient.UserPatchBody{}
		patchReset.SetFavoriteCurrencyId("")

		user, httpResp, err := client.UserAPI.UpdateUserFavoriteCurrency(ctx).
			UserPatchBody(patchReset).
			Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(user).ToNot(BeNil())
		Expect(user.FavoriteCurrencyId).To(BeNil())

		user2, httpResp, err := client.UserAPI.GetUser(ctx).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp).ToNot(BeNil())
		Expect(httpResp.StatusCode).To(Equal(200))
		Expect(user2).ToNot(BeNil())
		Expect(user2.FavoriteCurrencyId).To(BeNil())
	})
})

func getAccessToken(client *goclient.APIClient, ctx context.Context) string {
	req := client.AuthAPI.Authorize(ctx).AuthData(goclient.AuthData{
		Email:    User1,
		Password: Pass1,
	})
	resp, httpResp, err := req.Execute()
	if httpResp != nil {
		defer func() { _ = httpResp.Body.Close() }()
	}
	Expect(err).ToNot(HaveOccurred())
	Expect(httpResp).ToNot(BeNil())
	Expect(httpResp.StatusCode).To(Equal(200))
	Expect(resp).ToNot(BeNil())
	Expect(resp.Token).ToNot(BeEmpty())

	return resp.Token
}

package test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

var _ = Describe("GB", func() {
	ctx, cancel := context.WithCancel(context.Background())
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var cfg *config.Config

	BeforeEach(func() {
		hashed, err := auth.HashPassword([]byte("password1"))
		if err != nil {
			panic("Error hashing password")
		}

		cfg = &config.Config{
			Port:  0,
			Users: "user1:" + base64.StdEncoding.EncodeToString(hashed),
		}
	})

	It("starts and stops server", func() {
		_, finishCham, err := goserver.Serve(ctx, logger, cfg)
		Expect(err).To(BeNil())
		cancel()
		<-finishCham
	})
})

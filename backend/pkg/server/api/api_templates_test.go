package api_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

var _ = Describe("TemplatesAPI", func() {
	var (
		ctrl    *gomock.Controller
		mockDB  *mocks.MockStorage
		handler *api.TemplatesAPIServiceImpl
		ctx     context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(ctrl)
		handler = api.NewTemplatesAPIServiceImpl(
			slog.New(slog.NewTextHandler(os.Stderr, nil)),
			mockDB,
		)
		ctx = context.WithValue(context.Background(), constants.FamilyIDKey, uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("GetTemplates", func() {
		It("returns 200 with list of templates", func() {
			mockDB.EXPECT().GetTemplates(uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil).Return([]goserver.TransactionTemplate{
				{Id: "tpl-1", Name: "Rent"},
			}, nil)

			resp, err := handler.GetTemplates(ctx, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("returns 500 when userID is missing from context", func() {
			resp, err := handler.GetTemplates(context.Background(), "")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("CreateTemplate", func() {
		It("returns 400 when movements is empty", func() {
			resp, err := handler.CreateTemplate(ctx, goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})

		It("returns 200 with created template", func() {
			input := goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{{CurrencyId: "c1", AccountId: "a1"}},
			}
			mockDB.EXPECT().CreateTemplate(uuid.MustParse("00000000-0000-0000-0000-000000000001"), gomock.Any()).Return(goserver.TransactionTemplate{
				Id: "new-id", Name: "Rent",
			}, nil)

			resp, err := handler.CreateTemplate(ctx, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("UpdateTemplate", func() {
		It("returns 500 when userID is missing from context", func() {
			resp, err := handler.UpdateTemplate(context.Background(), "some-id", goserver.TransactionTemplateNoId{
				Name:      "Rent",
				Movements: []goserver.Movement{{CurrencyId: "c1", AccountId: "a1"}},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Describe("DeleteTemplate", func() {
		It("returns 204 on success", func() {
			mockDB.EXPECT().DeleteTemplate(uuid.MustParse("00000000-0000-0000-0000-000000000001"), "tpl-1").Return(nil)
			resp, err := handler.DeleteTemplate(ctx, "tpl-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Code).To(Equal(http.StatusNoContent))
		})
	})
})

package background

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Background Duplicate Detection", func() {
	var (
		mockCtrl *gomock.Controller
		mockDB   *mocks.MockStorage
		logger   = test.CreateTestLogger()
		userID   = "user1"
		ctx      = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("processUserDuplicates", func() {
		It("should detect duplicates and create notifications", func() {
			t1 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source1"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
			}
			t2 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source2"}, // different source
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)

			// Expect both to be marked suspicious
			mockDB.EXPECT().UpdateTransaction(userID, t1.Id, gomock.Any()).Return(goserver.Transaction{}, nil)
			mockDB.EXPECT().UpdateTransaction(userID, t2.Id, gomock.Any()).Return(goserver.Transaction{}, nil)

			// Expect notification
			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).DoAndReturn(func(uid string, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Type).To(Equal(string(models.NotificationTypeDuplicateDetected)))
				Expect(n.Description).To(ContainSubstring("2 potential duplicate"))
				return goserver.Notification{}, nil
			})

			processUserDuplicates(ctx, logger, mockDB, userID)
		})

		It("should skip transactions with different amounts", func() {
			t1 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source1"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
			}
			t2 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source2"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(200), CurrencyId: "USD"}},
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No UpdateTransaction or CreateNotification expected

			processUserDuplicates(ctx, logger, mockDB, userID)
		})

		It("should respect DuplicateDismissed flag", func() {
			t1 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source1"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
			}
			t2 := goserver.Transaction{
				Id:                 uuid.New().String(),
				Date:               time.Now(),
				ExternalIds:        []string{"source2"},
				Movements:          []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
				DuplicateDismissed: true, // User dismissed it
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No UpdateTransaction or CreateNotification expected

			processUserDuplicates(ctx, logger, mockDB, userID)
		})

		It("should not re-mark if already has DuplicateReason", func() {
			t1 := goserver.Transaction{
				Id:                uuid.New().String(),
				Date:              time.Now(),
				ExternalIds:       []string{"source1"},
				Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
				SuspiciousReasons: []string{models.DuplicateReason}, // Already marked
			}
			t2 := goserver.Transaction{
				Id:                uuid.New().String(),
				Date:              time.Now(),
				ExternalIds:       []string{"source2"},
				Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
				SuspiciousReasons: []string{models.DuplicateReason}, // Already marked
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No UpdateTransaction or CreateNotification expected

			processUserDuplicates(ctx, logger, mockDB, userID)
		})
	})
})

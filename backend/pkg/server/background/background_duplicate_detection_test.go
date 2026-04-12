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
		familyID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
		ctx      = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("processFamilyDuplicates", func() {
		It("should detect duplicates and create notifications", func() {
			// Simulate inter-account transfer: outgoing from account A, incoming to account B
			t1 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source1"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
			}
			t2 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source2"}, // different source
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
			}

			mockDB.EXPECT().GetTransactions(familyID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			mockDB.EXPECT().AddDuplicateRelationship(familyID, t1.Id, t2.Id).Return(nil)

			// Expect both to be marked suspicious via internal method (system operation, no user notifications)
			mockDB.EXPECT().UpdateTransactionInternal(familyID, t1.Id, gomock.Any()).Return(goserver.Transaction{}, nil)
			mockDB.EXPECT().UpdateTransactionInternal(familyID, t2.Id, gomock.Any()).Return(goserver.Transaction{}, nil)

			// Expect notification
			mockDB.EXPECT().CreateNotification(familyID, gomock.Any()).DoAndReturn(func(uid uuid.UUID, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Type).To(Equal(string(models.NotificationTypeDuplicateDetected)))
				Expect(n.Description).To(ContainSubstring("2 potential duplicate"))
				return goserver.Notification{}, nil
			})

			processFamilyDuplicates(ctx, logger, mockDB, familyID)
		})

		It("should not flag two same-direction transactions as duplicates", func() {
			// Two separate purchases of the same amount on consecutive days are NOT duplicates
			t1 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source1"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
			}
			t2 := goserver.Transaction{
				Id:          uuid.New().String(),
				Date:        time.Now(),
				ExternalIds: []string{"source2"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
			}

			mockDB.EXPECT().GetTransactions(familyID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No AddDuplicateRelationship, UpdateTransactionInternal, or CreateNotification expected

			processFamilyDuplicates(ctx, logger, mockDB, familyID)
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

			mockDB.EXPECT().GetTransactions(familyID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No UpdateTransaction or CreateNotification expected

			processFamilyDuplicates(ctx, logger, mockDB, familyID)
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

			mockDB.EXPECT().GetTransactions(familyID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// No UpdateTransaction or CreateNotification expected

			processFamilyDuplicates(ctx, logger, mockDB, familyID)
		})

		It("should not re-mark if already has DuplicateReason", func() {
			// Opposite directions so the pair passes the direction check
			t1 := goserver.Transaction{
				Id:                uuid.New().String(),
				Date:              time.Now(),
				ExternalIds:       []string{"source1"},
				Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
				SuspiciousReasons: []string{models.DuplicateReason}, // Already marked
			}
			t2 := goserver.Transaction{
				Id:                uuid.New().String(),
				Date:              time.Now(),
				ExternalIds:       []string{"source2"},
				Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "USD"}},
				SuspiciousReasons: []string{models.DuplicateReason}, // Already marked
			}

			mockDB.EXPECT().GetTransactions(familyID, gomock.Any(), time.Time{}, false).Return([]goserver.Transaction{t1, t2}, nil)
			// AddDuplicateRelationship still called even when already marked (idempotent link)
			mockDB.EXPECT().AddDuplicateRelationship(familyID, t1.Id, t2.Id).Return(nil)
			// No UpdateTransactionInternal or CreateNotification expected (both already marked)

			processFamilyDuplicates(ctx, logger, mockDB, familyID)
		})
	})
})

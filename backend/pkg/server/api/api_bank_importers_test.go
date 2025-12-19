package api

import (
	"regexp"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("BankImporters API", func() {
	var (
		mockCtrl *gomock.Controller
		mockDB   *mocks.MockStorage
		sut      *BankImportersAPIServiceImpl
		logger   = test.CreateTestLogger()
		userID   = "user1"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
		sut = NewBankImportersAPIServiceImpl(logger, mockDB)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("saveImportedTransactions", func() {
		It("should auto-convert transactions when a perfect matcher is found", func() {
			// Setup perfect matcher
			matcherID := uuid.New()
			matcher := goserver.Matcher{
				Id:                  matcherID.String(),
				OutputDescription:   "Converted Desc",
				OutputAccountId:     "acc1",
				OutputTags:          []string{"tag1"},
				ConfirmationHistory: []bool{true, true, true, true, true, true, true, true, true, true}, // 10 confirmations
			}
			// Logic in common.Match uses specific fields to match. Let's set description regex.
			matcher.DescriptionRegExp = "^Test Transaction$"
			descRegex, err := regexp.Compile(matcher.DescriptionRegExp)
			Expect(err).ToNot(HaveOccurred())

			runtimeMatcher := database.MatcherRuntime{
				Matcher:           &matcher,
				DescriptionRegexp: descRegex,
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{runtimeMatcher}, nil)

			// Expect auto-confirmation
			mockDB.EXPECT().AddMatcherConfirmation(userID, matcherID.String(), true).Return(nil)

			// Expect transaction creation with auto-converted fields
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				Expect(t.Description).To(Equal("Converted Desc"))
				Expect(t.IsAuto).To(BeTrue())
				Expect(t.MatcherId).To(Equal(matcherID.String()))
				// Also check movements accountId override
				Expect(t.Movements[0].AccountId).To(Equal("acc1"))
				return goserver.Transaction{Id: uuid.New().String()}, nil
			})

			// Mock updateLastImportFields
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil)
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					Description: "Test Transaction",
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: -100, CurrencyId: "USD"}},
				},
			}

			_, err = sut.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, transactions)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should NOT auto-convert when matcher is NOT perfect", func() {
			// Setup imperfect matcher
			matcherID := uuid.New()
			matcher := goserver.Matcher{
				Id:                  matcherID.String(),
				OutputDescription:   "Converted Desc",
				ConfirmationHistory: []bool{true}, // Only 1 confirmation
			}
			matcher.DescriptionRegExp = "^Test Transaction$"
			descRegex, err := regexp.Compile(matcher.DescriptionRegExp)
			Expect(err).ToNot(HaveOccurred())

			runtimeMatcher := database.MatcherRuntime{
				Matcher:           &matcher,
				DescriptionRegexp: descRegex,
			}

			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{runtimeMatcher}, nil)

			// Expect normal transaction creation without auto-conversion
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				Expect(t.Description).To(Equal("Test Transaction")) // Unchanged
				Expect(t.IsAuto).To(BeFalse())
				Expect(t.MatcherId).To(BeEmpty())
				return goserver.Transaction{Id: uuid.New().String()}, nil
			})

			// Mock updateLastImportFields
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil)
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					Description: "Test Transaction",
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: -100, CurrencyId: "USD"}},
				},
			}

			_, err = sut.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, transactions)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

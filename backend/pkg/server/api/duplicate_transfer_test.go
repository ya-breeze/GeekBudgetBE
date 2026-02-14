package api

import (
	"context"
	"regexp"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/test"
)

var _ = Describe("Duplicate Transfer Handling", func() {
	var (
		mockCtrl *gomock.Controller
		mockDB   *mocks.MockStorage
		sutBI    *BankImportersAPIServiceImpl
		sutUT    *UnprocessedTransactionsAPIServiceImpl
		logger   = test.CreateTestLogger()
		userID   = "user1"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockDB = mocks.NewMockStorage(mockCtrl)
		cfg := &config.Config{
			BankImporterFilesPath: "storage/bank-importer-files",
		}
		sutBI = NewBankImportersAPIServiceImpl(logger, mockDB, cfg)
		sutUT = &UnprocessedTransactionsAPIServiceImpl{logger: logger, db: mockDB}

		mockDB.EXPECT().CountUnprocessedTransactionsForAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil).AnyTimes()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("saveImportedTransactions duplicate transfer detection", func() {
		It("should skip auto-match when a similar transfer already exists in DB", func() {
			// (Same as before but using sutBI)
			// 1. Setup existing transfer in DB
			existingDate := time.Now().Add(-time.Hour * 24)
			existingTx := goserver.Transaction{
				Id:   uuid.New().String(),
				Date: existingDate,
				Movements: []goserver.Movement{
					{AccountId: "accA", Amount: decimal.NewFromInt(-1000), CurrencyId: "CZK"},
					{AccountId: "accB", Amount: decimal.NewFromInt(1000), CurrencyId: "CZK"},
				},
				Description: "Existing Transfer",
				ExternalIds: []string{"ext-existing"},
			}

			// 2. Setup perfect matcher for the new import
			matcherID := uuid.New().String()
			matcher := goserver.Matcher{
				Id:                  matcherID,
				OutputDescription:   "Matched Transfer",
				OutputAccountId:     "accA", // Matcher fills the second leg
				ConfirmationHistory: []bool{true, true, true, true, true, true, true, true, true, true},
				DescriptionRegExp:   "Imported Transfer",
			}
			r, _ := regexp.Compile(matcher.DescriptionRegExp)
			runtimeMatcher := database.MatcherRuntime{Matcher: &matcher, DescriptionRegexp: r}

			// 3. New imported transaction (from Account B)
			importedTx := goserver.TransactionNoId{
				Date:        existingDate.Add(time.Hour), // Slightly different time, same day
				Description: "Imported Transfer",
				ExternalIds: []string{"ext-imported"},
				Movements: []goserver.Movement{
					{AccountId: "accB", Amount: decimal.NewFromInt(1000), CurrencyId: "CZK"},
					{AccountId: "", Amount: decimal.NewFromInt(-1000), CurrencyId: "CZK"}, // Second leg empty
				},
			}

			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{existingTx}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{runtimeMatcher}, nil)

			// EXPECT: Transaction created but NOT as auto, and with skip reason
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				Expect(t.IsAuto).To(BeFalse())
				Expect(t.AutoMatchSkipReason).To(ContainSubstring("Potential duplicate detected"))
				Expect(t.Movements[1].AccountId).To(Equal(""), "Movements should NOT be modified by matcher if skip")
				return goserver.Transaction{Id: uuid.New().String()}, nil
			})

			// Mock importer update
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{}, nil).AnyTimes()
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			_, err := sutBI.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, []goserver.TransactionNoId{importedTx}, false)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("ProcessUnprocessedTransactionsAgainstMatcher duplicate transfer detection", func() {
		It("should skip auto-match when a similar transfer already exists in DB", func() {
			ctx := context.WithValue(context.Background(), constants.UserIDKey, userID)
			existingDate := time.Now().Add(-time.Hour * 24)

			// 1. Existing transfer in DB
			existingTx := goserver.Transaction{
				Id:   uuid.New().String(),
				Date: existingDate,
				Movements: []goserver.Movement{
					{AccountId: "accA", Amount: decimal.NewFromInt(-1000), CurrencyId: "CZK"},
					{AccountId: "accB", Amount: decimal.NewFromInt(1000), CurrencyId: "CZK"},
				},
				Description: "Existing Transfer",
			}

			// 2. Unprocessed transaction to be matched
			unprocessedTx := goserver.Transaction{
				Id:   uuid.New().String(),
				Date: existingDate,
				Movements: []goserver.Movement{
					{AccountId: "accB", Amount: decimal.NewFromInt(1000), CurrencyId: "CZK"},
					{AccountId: "", Amount: decimal.NewFromInt(-1000), CurrencyId: "CZK"},
				},
				Description: "Unprocessed Transfer",
			}

			// 3. Perfect matcher
			matcherID := uuid.New().String()
			matcher := goserver.Matcher{
				Id:                  matcherID,
				OutputDescription:   "Matched Transfer",
				OutputAccountId:     "accA",
				DescriptionRegExp:   "Unprocessed Transfer",
				ConfirmationHistory: []bool{true, true, true, true, true, true, true, true, true, true},
			}
			r, _ := regexp.Compile(matcher.DescriptionRegExp)
			runtimeMatcher := database.MatcherRuntime{Matcher: &matcher, DescriptionRegexp: r}

			// Mock DB expectations
			mockDB.EXPECT().GetMatcher(userID, matcherID).Return(matcher, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{runtimeMatcher}, nil)
			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), gomock.Any(), gomock.Any()).Return([]goserver.Transaction{existingTx, unprocessedTx}, nil)
			mockDB.EXPECT().GetAccounts(userID).Return([]goserver.Account{
				{Id: "accA"},
				{Id: "accB"},
			}, nil)
			mockDB.EXPECT().GetAccount(userID, "accA").Return(goserver.Account{Id: "accA"}, nil).AnyTimes()
			mockDB.EXPECT().GetTransactions(userID, gomock.Any(), gomock.Any(), gomock.Any()).Return([]goserver.Transaction{existingTx, unprocessedTx}, nil).AnyTimes()

			// EXPECT: Transaction updated but NOT as auto, and with skip reason
			mockDB.EXPECT().UpdateTransaction(userID, unprocessedTx.Id, gomock.Any()).DoAndReturn(func(uid, id string, t goserver.TransactionNoIdInterface) (goserver.Transaction, error) {
				Expect(t.GetIsAuto()).To(BeFalse())
				Expect(t.GetAutoMatchSkipReason()).To(ContainSubstring("Potential duplicate detected"))
				Expect(t.GetMovements()[1].AccountId).To(Equal(""), "Movements should NOT be modified by matcher if skip")
				return goserver.Transaction{}, nil
			})

			processed, err := sutUT.ProcessUnprocessedTransactionsAgainstMatcher(ctx, userID, matcherID, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(processed).To(ContainElement(unprocessedTx.Id))
		})
	})
})

package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/test"
)

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

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

		// Default expectation for balance checks triggered during imports
		// CountUnprocessedTransactionsForAccount returning 1 causes early return in CheckBalanceForAccount
		// Tests that need balance check should add their own GetAccount mock first
		mockDB.EXPECT().CountUnprocessedTransactionsForAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil).AnyTimes()
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

			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
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
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil).AnyTimes()
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					Description: "Test Transaction",
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
				},
			}

			_, err = sut.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, transactions, false)
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

			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{runtimeMatcher}, nil)

			// Expect normal transaction creation without auto-conversion
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				Expect(t.Description).To(Equal("Test Transaction")) // Unchanged
				Expect(t.IsAuto).To(BeFalse())
				Expect(t.MatcherId).To(BeEmpty())
				return goserver.Transaction{Id: uuid.New().String()}, nil
			})

			// Mock updateLastImportFields
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil).AnyTimes()
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					Description: "Test Transaction",
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
				},
			}

			_, err = sut.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, transactions, false)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should NOT auto-convert when multiple matchers match", func() {
			// Setup two matchers that both match the same transaction
			matcher1ID := uuid.New().String()
			matcher1 := goserver.Matcher{
				Id:                  matcher1ID,
				OutputDescription:   "Desc 1",
				ConfirmationHistory: []bool{true, true, true, true, true, true, true, true, true, true},
				DescriptionRegExp:   "Test Transaction",
			}
			r1, _ := regexp.Compile(matcher1.DescriptionRegExp)

			matcher2ID := uuid.New().String()
			matcher2 := goserver.Matcher{
				Id:                  matcher2ID,
				OutputDescription:   "Desc 2",
				ConfirmationHistory: []bool{true, true, true, true, true, true, true, true, true, true},
				DescriptionRegExp:   "Test Transaction",
			}
			r2, _ := regexp.Compile(matcher2.DescriptionRegExp)

			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{
				{Matcher: &matcher1, DescriptionRegexp: r1},
				{Matcher: &matcher2, DescriptionRegexp: r2},
			}, nil)

			// Expect NORMAL transaction creation (not auto-converted) because of conflict
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				Expect(t.IsAuto).To(BeFalse())
				Expect(t.MatcherId).To(BeEmpty())
				return goserver.Transaction{Id: uuid.New().String()}, nil
			})

			// Mock updateLastImportFields
			mockDB.EXPECT().GetBankImporter(userID, "imp1").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil).AnyTimes()
			mockDB.EXPECT().UpdateBankImporter(userID, "imp1", gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					Description: "Test Transaction",
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
				},
			}

			_, err := sut.saveImportedTransactions(userID, "imp1", &goserver.BankAccountInfo{}, transactions, false)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Deduplication Logic", func() {
		It("should filter out duplicates from DB and within batch", func() {
			// 1. Setup existing transactions in DB
			existingTx := goserver.Transaction{
				Id:          uuid.New().String(),
				ExternalIds: []string{"ext-existing"},
				Date:        time.Now().AddDate(0, 0, -1),
			}

			// Mock DB return
			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).
				Return([]goserver.Transaction{existingTx}, nil)

			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{}, nil)

			// 2. Setup imported transactions
			// - txNew: completely new
			// - txDuplicateDB: matches existingTx by ExternalID
			// - txBatch1: first occurrence in batch
			// - txBatch2: duplicate of txBatch1 in batch (should be skipped)

			txNew := goserver.TransactionNoId{
				Date:        time.Now(),
				Description: "New Tx",
				ExternalIds: []string{"ext-new"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(10), CurrencyId: "USD"}},
			}

			txDuplicateDB := goserver.TransactionNoId{
				Date:        time.Now().AddDate(0, 0, -1), // Date doesn't strictly matter for ExternalID match but good for realism
				Description: "Duplicate DB",
				ExternalIds: []string{"ext-existing"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(20), CurrencyId: "USD"}},
			}

			txBatch1 := goserver.TransactionNoId{
				Date:        time.Now(),
				Description: "Batch Tx",
				ExternalIds: []string{"ext-batch"},
				Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(30), CurrencyId: "USD"}},
			}
			// Exact copy of txBatch1
			txBatch2 := txBatch1

			importedTransactions := []goserver.TransactionNoId{txNew, txDuplicateDB, txBatch1, txBatch2}

			// Expect CreateTransaction ONLY for txNew and txBatch1
			// DuplicateDB and Batch2 should be skipped

			savedTxs := 0
			// We can capture the arguments to verify WHICH ones are saved
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).DoAndReturn(func(uid string, t *goserver.TransactionNoId) (goserver.Transaction, error) {
				if t.Description == "Duplicate DB" {
					Fail("Should not save duplicate from DB")
				}
				// We can't easily distinguish Batch1 and Batch2 by content since they are identical,
				// but the logic guarantees only one is saved.
				savedTxs++
				return goserver.Transaction{Id: uuid.New().String()}, nil
			}).Times(2) // We expect exactly 2 calls

			// Mock updateLastImportFields
			mockDB.EXPECT().GetBankImporter(userID, "imp-dedup").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil).AnyTimes()
			mockDB.EXPECT().UpdateBankImporter(userID, "imp-dedup", gomock.Any()).DoAndReturn(func(uid, id string, bi goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
				// Verify counts in description?
				// The implementation calls updateLastImportFields with totalTransactionsCnt=4, newTransactionsCnt=2
				// We can just return success
				return goserver.BankImporter{}, nil
			})

			_, err := sut.saveImportedTransactions(userID, "imp-dedup", &goserver.BankAccountInfo{}, importedTransactions, false)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Suspicious Logic", func() {
		It("should mark missing transactions as suspicious when checkMissing is true", func() {
			importerID := "imp-suspicious"
			accountID := "acc-suspicious"

			// 1. Setup existing transactions in DB (Active) both belonging to the account
			txPresent := goserver.Transaction{
				Id:          uuid.New().String(),
				ExternalIds: []string{"ext-present"},
				Date:        time.Now(),
				Movements:   []goserver.Movement{{AccountId: accountID, Amount: decimal.NewFromInt(-100), CurrencyId: "USD"}},
			}
			txMissing := goserver.Transaction{
				Id:          uuid.New().String(),
				ExternalIds: []string{"ext-missing"},
				Date:        time.Now(),
				Movements:   []goserver.Movement{{AccountId: accountID, Amount: decimal.NewFromInt(-200), CurrencyId: "USD"}},
			}
			txOtherAccount := goserver.Transaction{
				Id:          uuid.New().String(),
				ExternalIds: []string{"ext-other"},
				Date:        time.Now(),
				Movements:   []goserver.Movement{{AccountId: "other-acc", Amount: decimal.NewFromInt(-300), CurrencyId: "USD"}},
			}

			// Mock DB return
			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).
				Return([]goserver.Transaction{txPresent, txMissing, txOtherAccount}, nil)

			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{}, nil)

			// Mock GetBankImporter to return account ID
			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(goserver.BankImporter{Id: importerID, AccountId: accountID}, nil).AnyTimes()

			// Expect GetAccount when checkMissing is true (AnyTimes for balance check too)
			mockDB.EXPECT().GetAccount(userID, accountID).Return(goserver.Account{Id: accountID, BankInfo: goserver.BankAccountInfo{}}, nil).AnyTimes()
			mockDB.EXPECT().UpdateAccount(userID, accountID, gomock.Any()).Return(goserver.Account{}, nil)

			// 2. Setup imported transactions (Only txPresent)
			importedTransactions := []goserver.TransactionNoId{
				{
					Date:        txPresent.Date,
					Description: "Present Tx",
					ExternalIds: []string{"ext-present"},
					Movements:   txPresent.Movements,
				},
			}

			// Expect UpdateTransaction for txMissing with Suspicious=true
			mockDB.EXPECT().UpdateTransaction(userID, txMissing.Id, gomock.Any()).DoAndReturn(func(uid, id string, t goserver.TransactionNoIdInterface) (goserver.Transaction, error) {
				Expect(t.GetSuspiciousReasons()).To(Equal([]string{"Not present in importer transactions"}))
				return goserver.Transaction{}, nil
			})

			// Expect NO update for txPresent (it matches) and txOtherAccount (wrong account)
			// (Mock controller ensures unexpected calls fail)

			// Expect notification for suspicious transactions
			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).DoAndReturn(func(uid string, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Title).To(Equal("Suspicious Transactions Detected"))
				Expect(n.Type).To(Equal(string(models.NotificationTypeInfo)))
				return goserver.Notification{}, nil
			})

			// Mock updateLastImportFields
			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil)

			info := &goserver.BankAccountInfo{
				Balances: []goserver.BankAccountInfoBalancesInner{
					{CurrencyId: "USD", OpeningBalance: decimal.NewFromInt(1000), ClosingBalance: decimal.NewFromInt(900)},
				},
			}
			_, err := sut.saveImportedTransactions(userID, importerID, info, importedTransactions, true)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("FetchBankImporter", func() {
		It("should reset FetchAll and create notification on failure when FetchAll is true", func() {
			ctx := context.WithValue(context.Background(), common.UserIDKey, userID)
			importerID := "imp-fail-all"
			bi := goserver.BankImporter{
				Id:       importerID,
				Name:     "Failing Importer",
				FetchAll: true,
				Extra:    "some-token",
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(bi, nil).AnyTimes()
			mockDB.EXPECT().GetCurrencies(userID).Return([]goserver.Currency{}, nil)

			// Expect FetchAll to be reset
			updateCall1 := mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).DoAndReturn(func(uid, id string, data goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
				Expect(data.GetFetchAll()).To(BeFalse())
				return goserver.BankImporter{}, nil
			})
			// Expect result logging
			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil).After(updateCall1)

			// Expect notification
			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).Return(goserver.Notification{}, nil)

			// Mock HTTP transport failure
			oldTransport := http.DefaultClient.Transport
			defer func() { http.DefaultClient.Transport = oldTransport }()
			http.DefaultClient.Transport = &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewBufferString("fail"))}, nil
				},
			}

			resp, err := sut.FetchBankImporter(ctx, importerID)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Code).To(Equal(500))
		})

		It("should set IsStopped on failure when not interactive and not FetchAll", func() {
			// This test simulates a background fetch (isInteractive=false)
			// Since we can't easily call Fetch with private isInteractive from public API, we call Fetch directly or mock internal logic?
			// But Fetch is public on struct but not interface? No, Fetch is public helper on implementation.

			ctx := context.WithValue(context.Background(), common.UserIDKey, userID)
			importerID := "imp-fail-stopped"
			bi := goserver.BankImporter{
				Id:        importerID,
				Name:      "Failing Importer Stopped",
				FetchAll:  false,
				IsStopped: false,
				Extra:     "fail-token",
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(bi, nil).AnyTimes()
			mockDB.EXPECT().GetCurrencies(userID).Return([]goserver.Currency{}, nil)

			// Expect IsStopped to be set to true
			updateCall1 := mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).DoAndReturn(func(uid, id string, data goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
				Expect(data.GetIsStopped()).To(BeTrue())
				return goserver.BankImporter{}, nil
			})
			// Expect result logging
			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil).After(updateCall1)

			// Expect notification for stopped importer
			mockDB.EXPECT().CreateNotification(userID, gomock.Any()).DoAndReturn(func(uid string, n *goserver.Notification) (goserver.Notification, error) {
				Expect(n.Title).To(Equal("Bank Import Stopped"))
				return goserver.Notification{}, nil
			})

			// Mock HTTP transport failure
			oldTransport := http.DefaultClient.Transport
			defer func() { http.DefaultClient.Transport = oldTransport }()
			http.DefaultClient.Transport = &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewBufferString("fail"))}, nil
				},
			}

			// Call Fetch directly with isInteractive=false
			_, err := sut.Fetch(ctx, userID, importerID, false)
			Expect(err).To(HaveOccurred())
		})

		It("should skip fetch if IsStopped is true and not interactive", func() {
			ctx := context.WithValue(context.Background(), common.UserIDKey, userID)
			importerID := "imp-skipped"
			bi := goserver.BankImporter{
				Id:        importerID,
				IsStopped: true,
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(bi, nil).AnyTimes()

			// Expect result logging (UpdateBankImporter called by addImportResult)
			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil).AnyTimes()

			// Call Fetch directly with isInteractive=false
			_, err := sut.Fetch(ctx, userID, importerID, false)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("bank importer is stopped"))
		})

		It("should reset IsStopped on successful fetch", func() {
			ctx := context.WithValue(context.Background(), common.UserIDKey, userID)
			importerID := "imp-reset"
			bi := goserver.BankImporter{
				Id:        importerID,
				IsStopped: true,
				Type:      "fio", // Need valid type for mock logic
				Extra:     "good-token",
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(bi, nil).AnyTimes()
			mockDB.EXPECT().GetCurrencies(userID).Return([]goserver.Currency{}, nil).AnyTimes()

			// Mock successful import sequence...
			// This requires mocking FIO converter or ensuring FioConverter works with mock DB/Transport.
			// FioConverter uses http.DefaultClient.
			oldTransport := http.DefaultClient.Transport
			defer func() { http.DefaultClient.Transport = oldTransport }()
			http.DefaultClient.Transport = &mockTransport{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					// Return empty JSON list for transactions
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewBufferString(`{"accountStatement": {"transactionList": {"transaction": []}}}`)),
					}, nil
				},
			}

			// Expect IsStopped to be reset to false
			updateCall1 := mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).DoAndReturn(func(uid, id string, data goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
				Expect(data.GetIsStopped()).To(BeFalse())
				return goserver.BankImporter{}, nil
			})

			// addImportResult will also call UpdateBankImporter...
			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil).After(updateCall1)
			// mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil) // Not called for empty transactions
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{}, nil).AnyTimes() // Might not be called either, but safe to allow any times or just remove if strict

			_, err := sut.Fetch(ctx, userID, importerID, true)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("isDuplicate", func() {
		It("should not match transactions with same amount but different currencies", func() {
			t1 := &goserver.TransactionNoId{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(123), CurrencyId: "CZK"},
				},
			}
			t2 := &goserver.Transaction{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(123), CurrencyId: "EUR"},
				},
			}

			Expect(common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements)).To(BeFalse(), "123 CZK should not match 123 EUR")
		})

		It("should match transactions with same amount and same currency", func() {
			t1 := &goserver.TransactionNoId{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(123), CurrencyId: "CZK"},
				},
			}
			t2 := &goserver.Transaction{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(123), CurrencyId: "CZK"},
				},
			}

			Expect(common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements)).To(BeTrue())
		})

		It("should not match transactions with different amounts in same currency", func() {
			t1 := &goserver.TransactionNoId{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(123), CurrencyId: "CZK"},
				},
			}
			t2 := &goserver.Transaction{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(124), CurrencyId: "CZK"},
				},
			}

			Expect(common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements)).To(BeFalse())
		})

		It("should handle complex multi-movement matches", func() {
			t1 := &goserver.TransactionNoId{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(100), CurrencyId: "CZK"},
					{Amount: decimal.NewFromInt(-100), CurrencyId: "CZK"},
					{Amount: decimal.NewFromInt(50), CurrencyId: "EUR"},
				},
			}
			t2 := &goserver.Transaction{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(100), CurrencyId: "CZK"},
					{Amount: decimal.NewFromInt(50), CurrencyId: "EUR"},
				},
			}
			// Both have 100 CZK increase and 50 EUR increase
			Expect(common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements)).To(BeTrue())
		})

		It("should not match if currencies set differs", func() {
			t1 := &goserver.TransactionNoId{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(100), CurrencyId: "CZK"},
					{Amount: decimal.NewFromInt(50), CurrencyId: "EUR"},
				},
			}
			t2 := &goserver.Transaction{
				Date: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				Movements: []goserver.Movement{
					{Amount: decimal.NewFromInt(100), CurrencyId: "CZK"},
				},
			}
			Expect(common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements)).To(BeFalse())
		})
	})

	Describe("Opening Balance Logic", func() {
		It("should update OpeningBalance when checkMissing is true", func() {
			importerID := "imp-full"
			accountID := "acc-full"
			currencyID := "USD"

			// Existing account with an opening balance
			existingAccount := goserver.Account{
				Id: accountID,
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{CurrencyId: currencyID, OpeningBalance: decimal.NewFromInt(1000), ClosingBalance: decimal.NewFromInt(1500)},
					},
				},
			}

			// Imported info with NEW opening balance
			importedInfo := &goserver.BankAccountInfo{
				Balances: []goserver.BankAccountInfoBalancesInner{
					{CurrencyId: currencyID, OpeningBalance: decimal.NewFromInt(2000), ClosingBalance: decimal.NewFromInt(2500)},
				},
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(goserver.BankImporter{Id: importerID, AccountId: accountID}, nil).AnyTimes()
			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{}, nil)
			mockDB.EXPECT().GetAccount(userID, accountID).Return(existingAccount, nil).AnyTimes()

			// EXPECT: Entire balance object replaced (OpeningBalance updated to 2000)
			mockDB.EXPECT().UpdateAccount(userID, accountID, gomock.Any()).DoAndReturn(func(uid, aid string, acc *goserver.AccountNoId) (goserver.Account, error) {
				Expect(acc.BankInfo.Balances[0].OpeningBalance).To(Equal(decimal.NewFromInt(2000)))
				Expect(acc.BankInfo.Balances[0].ClosingBalance).To(Equal(decimal.NewFromInt(2500)))
				return goserver.Account{}, nil
			})

			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					ExternalIds: []string{"ext1"},
					Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: currencyID}},
				},
			}
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).Return(goserver.Transaction{}, nil)

			_, err := sut.saveImportedTransactions(userID, importerID, importedInfo, transactions, true) // checkMissing=true
			Expect(err).ToNot(HaveOccurred())
		})

		It("should NOT update OpeningBalance (only ClosingBalance) when checkMissing is false", func() {
			importerID := "imp-inc"
			accountID := "acc-inc"
			currencyID := "USD"

			// Existing account with an opening balance
			existingAccount := goserver.Account{
				Id: accountID,
				BankInfo: goserver.BankAccountInfo{
					Balances: []goserver.BankAccountInfoBalancesInner{
						{CurrencyId: currencyID, OpeningBalance: decimal.NewFromInt(1000), ClosingBalance: decimal.NewFromInt(1500)},
					},
				},
			}

			// Imported info with DIFFERENT opening balance (which should be ignored)
			importedInfo := &goserver.BankAccountInfo{
				Balances: []goserver.BankAccountInfoBalancesInner{
					{CurrencyId: currencyID, OpeningBalance: decimal.NewFromInt(2000), ClosingBalance: decimal.NewFromInt(2500)},
				},
			}

			mockDB.EXPECT().GetBankImporter(userID, importerID).Return(goserver.BankImporter{Id: importerID, AccountId: accountID}, nil).AnyTimes()
			mockDB.EXPECT().GetTransactionsIncludingDeleted(userID, gomock.Any(), gomock.Any()).Return([]goserver.Transaction{}, nil)
			mockDB.EXPECT().GetMatchersRuntime(userID).Return([]database.MatcherRuntime{}, nil)
			mockDB.EXPECT().GetAccount(userID, accountID).Return(existingAccount, nil).AnyTimes()

			// EXPECT: OpeningBalance preserved at 1000, but ClosingBalance updated to 2500
			mockDB.EXPECT().UpdateAccount(userID, accountID, gomock.Any()).DoAndReturn(func(uid, aid string, acc *goserver.AccountNoId) (goserver.Account, error) {
				Expect(acc.BankInfo.Balances[0].OpeningBalance).To(Equal(decimal.NewFromInt(1000))) // Preserved
				Expect(acc.BankInfo.Balances[0].ClosingBalance).To(Equal(decimal.NewFromInt(2500))) // Updated
				return goserver.Account{}, nil
			})

			mockDB.EXPECT().UpdateBankImporter(userID, importerID, gomock.Any()).Return(goserver.BankImporter{}, nil)

			transactions := []goserver.TransactionNoId{
				{
					Date:        time.Now(),
					ExternalIds: []string{"ext2"},
					Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(200), CurrencyId: currencyID}},
				},
			}
			mockDB.EXPECT().CreateTransaction(userID, gomock.Any()).Return(goserver.Transaction{}, nil)

			_, err := sut.saveImportedTransactions(userID, importerID, importedInfo, transactions, false) // checkMissing=false
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

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
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/mocks"
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
				Movements:   []goserver.Movement{{Amount: 10, CurrencyId: "USD"}},
			}

			txDuplicateDB := goserver.TransactionNoId{
				Date:        time.Now().AddDate(0, 0, -1), // Date doesn't strictly matter for ExternalID match but good for realism
				Description: "Duplicate DB",
				ExternalIds: []string{"ext-existing"},
				Movements:   []goserver.Movement{{Amount: 20, CurrencyId: "USD"}},
			}

			txBatch1 := goserver.TransactionNoId{
				Date:        time.Now(),
				Description: "Batch Tx",
				ExternalIds: []string{"ext-batch"},
				Movements:   []goserver.Movement{{Amount: 30, CurrencyId: "USD"}},
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
			mockDB.EXPECT().GetBankImporter(userID, "imp-dedup").Return(goserver.BankImporter{LastImports: []goserver.ImportResult{}}, nil)
			mockDB.EXPECT().UpdateBankImporter(userID, "imp-dedup", gomock.Any()).DoAndReturn(func(uid, id string, bi goserver.BankImporterNoIdInterface) (goserver.BankImporter, error) {
				// Verify counts in description?
				// The implementation calls updateLastImportFields with totalTransactionsCnt=4, newTransactionsCnt=2
				// We can just return success
				return goserver.BankImporter{}, nil
			})

			_, err := sut.saveImportedTransactions(userID, "imp-dedup", &goserver.BankAccountInfo{}, importedTransactions)
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
			mockDB.EXPECT().GetCurrencies(userID).Return([]goserver.Currency{}, nil)

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
})

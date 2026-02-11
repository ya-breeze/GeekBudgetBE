package database

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type BulkReconciliationData struct {
	Balances              map[string]map[string]decimal.Decimal          // AccountID -> CurrencyID -> Balance
	LatestReconciliations map[string]map[string]*goserver.Reconciliation // AccountID -> CurrencyID -> Reconciliation
	UnprocessedCounts     map[string]int                                 // AccountID -> Count
	MaxTransactionDates   map[string]map[string]time.Time                // AccountID -> CurrencyID -> Max Date
}

# Reconciliation Logic & Status

Guidelines and facts discovered regarding the reconciliation system.

## Status Indicators
- **Green**: Delta is within `common.ReconciliationTolerance` (0.01) AND there are no unprocessed transactions.
- **Yellow**: Delta is within tolerance, but there are **unprocessed transactions**. Manual reconciliation is disabled.
- **Red**: Delta exceeds tolerance.

## FIO Account Status
- During investigation, Fio appeared yellow with 0.00 delta.
- **Fact**: This was due to uncategorized/unprocessed transactions in the database where movements had empty account IDs (e.g., "Frutisimo Zlicin").
- **Constraint**: Manual reconciliation is intentionally blocked by unprocessed transactions to ensure the "App Balance" is final.

## Tolerance
- Defined in `backend/pkg/server/common/constants.go`.
- Currently set to **0.01** to handle floating-point precision issues.

## Frontend UI Patterns
- **Disabled Buttons**: Angular Material tooltips on disabled buttons require a wrapper element (like a `<span>`) to capture mouse events.
- **Stale Balance Warning**: A ⚠️ icon appears if `hasTransactionsAfterBankBalance` is true, indicating the bank balance metadata might be older than the current transaction set.
+
+## Performance & Batching
+- **Optimization**: `GetReconciliationStatus` uses a bulk fetching strategy via `s.db.GetBulkReconciliationData(userID)`.
+- **Logic**: It aggregates balances, latest reconciliations, and unprocessed counts in memory after a single pass over the user's transactions.
+- **Data Availability**: The `BulkReconciliationData` struct (defined in `pkg/database/bulk_types.go`) provides efficient lookups by `AccountID` and `CurrencyID`.

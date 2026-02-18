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

## Performance & Batching
- **Optimization**: `GetReconciliationStatus` uses a bulk fetching strategy via `s.db.GetBulkReconciliationData(userID)`.
- **Logic**: It aggregates balances, latest reconciliations, and unprocessed counts in memory after a single pass over the user's transactions.
- **Data Availability**: The `BulkReconciliationData` struct (defined in `pkg/database/bulk_types.go`) provides efficient lookups by `AccountID` and `CurrencyID`.

## Disbalance Analysis ("Find Disbalance Cause")
- **Feature**: Uses a subset-sum algorithm (Tiered approach: singles, pairs, then DP-based subset search) to identify transactions that explain a delta.
- **Location**: `backend/pkg/server/common/disbalance_finder.go`.
- **Constraint**: Subset algorithm handles mixed signs and is limited to 50 transactions and 100,000 DP states for performance.
- **FE Interaction**: Highlight candidate transactions in the UI when analysis is active.

## Reconciliation History
- **Access**: Navigable from the reconciliation dashboard per account+currency.
- **Pattern**: `/reconciliation/history/:accountId/:currencyId`.
- **API Design Rule**: Core API models (like `Reconciliation`) should remain "lean". They should only store IDs (e.g., `accountId`, `currencyId`). The Frontend is responsible for resolving these IDs to human-readable names using separate entity lookups (Accounts/Currencies API).

## Testing & Verification
- **Subset Match**: When testing the disbalance finder, ensure target deltas are consistent with mock transaction data.
- **Decimal Comparison**: Always use `.Equal()` (scale-insensitive) for `decimal.Decimal` comparisons in Ginkgo tests.
- **Go Literals**: Use `decimal.NewFromFloat()` for decimal constants in tests to avoid type mismatch.

## Reconciliation History Preservation

### Invalidation Logic
- **Constraint**: Changing a transaction (update/delete) or inserting a retrospective transaction must invalidate future reconciliations to maintain data integrity.
- **Rule**: `InvalidateReconciliation` (in `storage_reconciliation.go`) accepts a `fromDate`. Only reconciliations with `reconciled_at >= fromDate` are deleted.
- **Trigger**: `InvalidateReconciliation` is called by `invalidateReconciliationIfAmountsChanged`, which is triggered by:
    - `UpdateTransaction`
    - `DeleteTransaction`
    - `CreateTransaction` (if the new transaction is older than the latest reconciliation)

### History Preservation
- **Behavior**: Historical reconciliations that predate the modified transaction are **preserved**. This ensures that past verified states remain intact even if data is corrected later.

### Notification Policy
- **Policy**: "Reconciliation Invalidated" notifications are only shown for manual user actions.
- **Rules**:
    - **Show** notification: Manual `UpdateTransaction` or `DeleteTransaction` by the user.
    - **Hide** notification: Background bank imports (`CreateTransaction`/`CreateTransactionsBatch`) and internal conversion tasks (`UpdateTransactionInternal` for auto-matching or unprocessed transaction conversion).
- **Goal**: Provides situational awareness for data integrity issues while avoiding noise during routine automation.

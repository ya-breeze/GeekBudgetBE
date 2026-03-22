# Design: Reconciliation for Accounts Without Bank Importers

**Date:** 2026-03-22
**Status:** Approved

## Problem

Accounts without a bank importer can have manual reconciliation enabled (`isManualReconciliationEnabled = true`). The reconciliation UI shows a "bank balance" derived from the last manual reconciliation record and a delta vs. the current app balance. However:

1. **Frontend**: The reconcile button is disabled whenever `|delta| > 0.01`, even for no-importer accounts. There is no other mechanism to update the balance reference point, so a user with a drifted balance is permanently stuck.
2. **Backend (latent bug)**: `ReconcileAccount` derives `expectedBalance` from `BankInfo.Balances.ClosingBalance`, which is always empty for no-importer accounts (`expectedBalance = 0`). The tolerance check `|balance - 0| > 0.01` therefore fails with HTTP 400 for any non-zero balance — manual reconciliation has never actually worked for non-trivial balances.

## Scope

Changes apply **only to accounts where `hasBankImporter = false`** (detected in the backend as `len(acc.BankInfo.Balances) == 0`). Accounts with bank importers are unchanged.

## Chosen Approach: Condition-Based Skip (Approach A)

No API contract changes. Fix the backend validation logic and update the frontend gate and UX. Smallest surface area.

---

## Backend Changes

**File:** `backend/pkg/server/api/api_reconciliation.go` — `ReconcileAccount` handler

After fetching the account (line ~148), branch on whether it has a bank importer:

```
hasImporter := len(acc.BankInfo.Balances) > 0

if hasImporter {
    // Existing logic: derive expectedBalance from BankInfo, enforce tolerance check
} else {
    // No-importer path:
    // - expectedBalance = balance (the app balance being confirmed)
    // - Skip tolerance validation entirely
    // - IsManual = true always
}
```

**Why `expectedBalance = balance`**: For no-importer accounts, there is no external source of truth. Setting `expectedBalance` to the confirmed app balance means the reconciliation history record reads "at this moment I confirmed balance is X" with a stored delta of 0 — semantically correct and historically useful.

**`IsManual` flag**: Currently set as `body.Balance.IsPositive()`. For no-importer accounts, always set `IsManual = true` regardless of whether balance was explicitly provided, since these are by definition manual.

---

## Frontend Changes

**Files:**
- `frontend/src/app/features/reconciliation/reconciliation.component.ts`
- `frontend/src/app/features/reconciliation/reconciliation.component.html`

### 1. Disable Condition

Current logic (applied to all accounts):
```
disabled if |delta| > 0.01 || hasUnprocessedTransactions
```

New logic:
- **Importer accounts**: unchanged
- **No-importer accounts** (`!hasBankImporter && isManualReconciliationEnabled`): disabled only if `hasUnprocessedTransactions`

### 2. Tooltip (`getReconcileTooltip`)

Add a case for no-importer accounts with large delta:
```
"Balance differs by {delta} from last reconciliation — click to confirm"
```

Existing cases (unprocessed transactions, within tolerance) are unchanged.

### 3. Confirmation Dialog

In the `reconcile()` method, before making the API call, check:
```
if (!element.hasBankImporter && Math.abs(element.delta) > RECONCILIATION_TOLERANCE)
```

If true, open a `MatDialog` confirmation with:
- **Message**: `"The current balance differs from the last reconciled balance by {delta}. This exceeds the normal tolerance. Are you sure the current balance is correct?"`
- **Actions**: Cancel | Confirm
- Only proceed with the API call on Confirm.

Use Angular Material's `MatDialog` (already a project dependency) with an inline confirmation component or `MatDialogRef` confirm pattern consistent with existing dialogs in the codebase.

---

## Data Flow

```
User clicks Reconcile (no-importer, large delta)
  → Frontend checks: !hasBankImporter && |delta| > 0.01
  → Opens confirmation dialog
  → User confirms
  → POST /v1/accounts/{id}/reconcile { currencyId, balance: 0 }
  → Backend: fetches account, detects no BankInfo.Balances
  → Sets expectedBalance = appBalance, skips tolerance check
  → Creates reconciliation record (IsManual=true, delta=0 in history)
  → 200 OK
  → Frontend refreshes reconciliation status
```

---

## What Does Not Change

- Accounts with bank importers: no behavior change
- The `enable-reconciliation` flow: unchanged
- The `hasUnprocessedTransactions` gate: still blocks reconciliation for all accounts
- OpenAPI spec: no new endpoints or request fields
- `make generate`: not required

---

## Testing

**Backend:**
- Unit test: `ReconcileAccount` with a no-importer account and large delta → expect 200 and correct record stored
- Unit test: `ReconcileAccount` with a no-importer account, delta within tolerance → expect 200 (unchanged happy path)
- Unit test: `ReconcileAccount` with an importer account and large delta → expect 400 (existing behavior preserved)

**Frontend:**
- Component test: reconcile button enabled when `!hasBankImporter && isManualReconciliationEnabled && !hasUnprocessedTransactions && |delta| > 0.01`
- Component test: confirmation dialog appears on click when delta is large and no importer
- Component test: API not called if user cancels the dialog

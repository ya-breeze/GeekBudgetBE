# Design: Reconciliation for Accounts Without Bank Importers

**Date:** 2026-03-22
**Status:** Approved

## Problem

Accounts without a bank importer can have manual reconciliation enabled (`isManualReconciliationEnabled = true`). The reconciliation UI shows a "bank balance" derived from the last manual reconciliation record and a delta vs. the current app balance. However:

1. **Frontend**: The reconcile button is disabled whenever `|delta| > 0.01`, even for no-importer accounts. There is no other mechanism to update the balance reference point, so a user with a drifted balance is permanently stuck.
2. **Backend (latent bug)**: `ReconcileAccount` derives `expectedBalance` from `BankInfo.Balances.ClosingBalance`, which is always empty for no-importer accounts (`expectedBalance = 0`). The tolerance check `|balance - 0| > 0.01` therefore fails with HTTP 400 for any non-zero balance — manual reconciliation has never actually worked for non-trivial balances.

## Scope

Changes apply **only to accounts where the user has no bank importer configured**, detected via `GetBankImporters(userID)` lookup (see Backend Changes). Accounts with bank importers are unchanged.

The `enableManual` flow (which creates the first reconciliation record and sets `isManualReconciliationEnabled = true`) is unchanged. The large-delta path applies only after manual reconciliation has been enabled — the user must go through `enableManual` first before the reconcile button appears.

## Chosen Approach: Condition-Based Skip (Approach A)

No API contract changes. Fix the backend validation logic and update the frontend gate and UX. Smallest surface area.

---

## Backend Changes

**File:** `backend/pkg/server/api/api_reconciliation.go` — `ReconcileAccount` handler

### Importer Detection

Do NOT use `len(acc.BankInfo.Balances) == 0` to detect no-importer accounts. `BankInfo.Balances` reflects whether import data has been written, not whether an importer is configured — a configured importer that has never run would appear identical to an account without one. Use the same mechanism as `GetReconciliationStatus`: call `s.db.GetBankImporters(userID)` and check whether any importer's `AccountId` matches the account being reconciled.

```go
importers, err := s.db.GetBankImporters(userID)
// handle err
hasImporter := false
for _, imp := range importers {
    if imp.AccountId == id {
        hasImporter = true
        break
    }
}
```

### Branch Logic

The existing `ReconcileAccount` handler has this structure (simplified):

```go
balance := body.Balance
if balance.IsZero() {
    balance = GetAccountBalance(...)   // resolve actual app balance
}

acc = GetAccount(...)                  // fetch account for BankInfo

var expectedBalance decimal.Decimal
for _, b := range acc.BankInfo.Balances {  // loop only relevant for importer accounts
    if b.CurrencyId == body.CurrencyId {
        expectedBalance = b.ClosingBalance
        break
    }
}

if balance.Sub(expectedBalance).Abs().GreaterThan(tolerance) {  // tolerance check
    return 400
}

CreateReconciliation(..., expectedBalance, IsManual: body.Balance.IsPositive())
```

After adding the `GetBankImporters` check, restructure as follows. The `balance` resolution block stays unconditional (needed by both paths). The `GetAccount` call and `BankInfo.Balances` loop move **inside the `hasImporter` branch** since they are only needed there. The `else` branch sets `expectedBalance = balance` directly and skips the tolerance check:

```go
// balance resolution — unconditional
balance := body.Balance
if balance.IsZero() {
    balance = GetAccountBalance(...)
}

var expectedBalance decimal.Decimal
isManual := true

if hasImporter {
    // existing logic: GetAccount, loop BankInfo.Balances, enforce tolerance check
    acc = GetAccount(...)
    for _, b := range acc.BankInfo.Balances {
        if b.CurrencyId == body.CurrencyId {
            expectedBalance = b.ClosingBalance
            break
        }
    }
    if balance.Sub(expectedBalance).Abs().GreaterThan(tolerance) {
        return 400
    }
    isManual = body.Balance.IsPositive()
} else {
    // no-importer path: confirm current app balance, no validation
    expectedBalance = balance
    // isManual = true (already set above)
}

CreateReconciliation(..., expectedBalance, IsManual: isManual)
```

**Why `expectedBalance = balance`**: For no-importer accounts, there is no external source of truth. Setting `expectedBalance` to the confirmed app balance means the reconciliation history record reads "at this moment I confirmed balance is X" with a stored delta of 0 — semantically correct and historically useful.

**`IsManual` flag**: For no-importer accounts, set `IsManual = true` unconditionally regardless of the balance value (zero, positive, or negative). The current logic `IsManual: body.Balance.IsPositive()` would incorrectly set `IsManual = false` when the fetched app balance is zero or negative.

---

## Frontend Changes

**Files:**
- `frontend/src/app/features/reconciliation/reconciliation.component.ts`
- `frontend/src/app/features/reconciliation/reconciliation.component.html`

### 1. Disable Condition

Current `[disabled]` binding in the HTML template:
```html
[disabled]="
    (element.delta || 0) > 0.01 ||
    (element.delta || 0) < -0.01 ||
    element.hasUnprocessedTransactions
"
```

Updated binding — delta conditions guarded by `hasBankImporter`:
```html
[disabled]="
    element.hasUnprocessedTransactions ||
    (element.hasBankImporter && ((element.delta || 0) > 0.01 || (element.delta || 0) < -0.01))
"
```

Effect:
- **Importer accounts**: unchanged — delta still blocks
- **No-importer accounts**: disabled only when `hasUnprocessedTransactions`

### 2. Tooltip (`getReconcileTooltip`)

Add a case for no-importer accounts with large delta:
```
"Balance differs by {delta} from last reconciliation — click to confirm"
```

Existing cases (unprocessed transactions, within tolerance) are unchanged.

### 3. Confirmation Dialog

In the `reconcile()` method, before making the API call, check:
```typescript
if (!element.hasBankImporter && Math.abs(element.delta ?? 0) > RECONCILIATION_TOLERANCE)
```

If true, open a confirmation dialog using the existing **`ConfirmationDialogComponent`** located at:
`frontend/src/app/shared/components/confirmation-dialog/confirmation-dialog.component.ts`

This component accepts a `ConfirmationDialogData` interface with `title`, `message`, `confirmText`, and `cancelText`.

`ReconciliationComponent` currently imports `MatDialogModule` but does not inject `MatDialog`. Add `MatDialog` to the constructor injection. Then:

```typescript
const dialogRef = this.dialog.open(ConfirmationDialogComponent, {
  data: {
    title: 'Confirm Balance',
    message: `The current balance differs from the last reconciled balance by ${Math.abs(element.delta ?? 0).toFixed(2)}. This exceeds the normal tolerance. Are you sure the current balance is correct?`,
    confirmText: 'Confirm',
    cancelText: 'Cancel',
  } as ConfirmationDialogData,
});
dialogRef.afterClosed().subscribe(confirmed => {
  if (confirmed) { /* proceed with API call */ }
});
```

Do not use `window.confirm()` (already an anti-pattern in `enableManual()` — do not replicate it here).

### 4. Row Status Color (`getStatusClass`)

The existing `getStatusClass` method colors a row red when `|delta| > RECONCILIATION_TOLERANCE` for all accounts. After this fix, a no-importer account with a large delta is in a valid "confirm pending" state — a red row is misleading. Add a branch:

```typescript
if (!status.hasBankImporter) {
    return status.hasUnprocessedTransactions ? 'status-yellow' : 'status-green';
}
// existing delta-based logic for importer accounts
```

This treats no-importer accounts as always "green" (or yellow if unprocessed), since the large delta is expected and resolvable via the confirm flow.

---

## Data Flow

```
User clicks Reconcile (no-importer, |delta| > 0.01)
  → Frontend: !hasBankImporter && |delta| > 0.01 → open ConfirmationDialog
  → User confirms
  → POST /v1/accounts/{id}/reconcile { currencyId, balance: 0 }
      Note: balance: 0 signals "use current app balance"
            Backend fetches actual balance via GetAccountBalance → e.g. 1000.00
  → Backend: calls GetBankImporters(userID), no importer found for this account
  → Skips tolerance check
  → Sets expectedBalance = 1000.00 (fetched app balance), IsManual = true
  → Creates reconciliation record (delta in history = 0, confirming balance is correct)
  → 200 OK
  → Frontend refreshes reconciliation status
```

---

## What Does Not Change

- Accounts with bank importers: no behavior change
- The `enable-reconciliation` flow: unchanged (still required before reconcile button appears)
- The `hasUnprocessedTransactions` gate: still blocks reconciliation for all accounts
- OpenAPI spec: no new endpoints or request fields
- `make generate`: not required

---

## Testing

**Backend:**
- `ReconcileAccount` with no-importer account and large delta → expect 200, record stored with `IsManual=true`, `ReconciledBalance == ExpectedBalance` (delta = 0 in history)
- `ReconcileAccount` with no-importer account and `body.Balance = 0` (triggers app-balance fetch) → expect 200, `IsManual=true` even if fetched balance is zero
- `ReconcileAccount` with no-importer account and delta within tolerance → expect 200 (happy path unchanged)
- `ReconcileAccount` with importer account and large delta → expect 400 (existing behavior preserved)
- `ReconcileAccount` with an account that has a `BankImporter` record but no `BankInfo.Balances` entries (importer configured but never run) → expect 400 (treated as has-importer, existing strict behavior)

**Frontend:**
- Reconcile button enabled when `!hasBankImporter && isManualReconciliationEnabled && !hasUnprocessedTransactions && |delta| > 0.01`
- Confirmation dialog opens on click when delta is large and no importer
- API not called if user cancels the dialog
- API called if user confirms the dialog
- Reconcile button still disabled when `hasUnprocessedTransactions`, regardless of importer status

---
description: how to manage and debug the duplicate transaction detection flow
---

# Duplicate Transaction Detection Workflow

This workflow describes how to interact with the background duplicate detection system.

## 1. Detection Logic
- The background task runs every 24 hours.
- It scans transactions from the last 30 days.
- It uses `common.IsDuplicate` (which checks date ±2 days and movement amounts).
- It only flags transactions that have *different* `ExternalIDs`.

## 2. Marking Suspicious
- If a potential duplicate is found, the system adds `"Potential duplicate from different importer"` to the `SuspiciousReasons` JSON array.
- This results in the transaction being highlighted to the user in the UI.

## 3. Dismissal Persistence (Crucial)
To prevent the background task from re-flagging a transaction that a user has already reviewed:
1. The user must clear the `SuspiciousReasons` list.
2. The user must set `DuplicateDismissed = true`.

**If you are implementing a dismissal feature or debugging why a transaction was re-flagged:**
- Check if `DuplicateDismissed` is being correctly set to `true` in the database.
- The background task **skips** any transaction where `DuplicateDismissed` is true.

## 4. Duplicate Linking
- Pairwise bidirectional relationships are stored in the `TransactionDuplicate` table.
- Use `db.GetDuplicateTransactionIDs` to find all linked duplicates for a transaction.
- Background task populates these links automatically when it detects suspicious duplicates.

## 5. User Resolution Flows

### Option A: Dismiss as Non-Duplicate
- Action: Update transaction with `DuplicateDismissed = true`.
- Consequence: Background task will skip this transaction in the future.
- Backend Auto-Action: `UpdateTransaction` automatically calls `ClearDuplicateRelationships`. This removes all bidirectional links and **synchronizes suspicious reasons** (clearing `models.DuplicateReason` from linked transactions if they have no other duplicate links).

### Option B: Merge Transactions
- Action: `POST /v1/transactions/merge` with `keepId` and `mergeId`.
- Process:
  1. Kept transaction inherits `ExternalIDs` from merged transaction.
  2. Merged transaction is soft-deleted using `GORM`'s `Delete` (sets `DeletedAt`) and tagged with `MergedIntoID`, `MergedAt`.
  3. Linked transactions are updated: `models.DuplicateReason` is removed if they no longer have any active duplicate links.

## 6. Debugging
- All duplicates are logged with the prefix `DUPLICATE DETECTION:`.
- Check `backend/pkg/server/background/background_duplicate_detection.go` for implementation details.
- Use `pkg/server/common/transactions_test.go` to verify the detection logic.
- Verify relationship consistency in the `transaction_duplicates` table (see `.agent/rules/database-investigation.md` for tools).

## 7. Revalidation on Update
- When a transaction is updated (date or amounts), `UpdateTransaction` automatically calls `RevalidateDuplicateRelationships`.
- This re-checks all linked transactions using `utils.IsDuplicate`. If they are no longer duplicates, the link is removed.
- Both affected transactions have their `models.DuplicateReason` cleared if they no longer have any active duplicate links.

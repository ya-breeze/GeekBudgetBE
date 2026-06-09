# ADR-010: Archive-on-Merge for Deduplicated Transactions

## Status
Accepted

## Context and Problem Statement

When the user resolves a duplicate by merging two transactions, the "loser" must disappear from the
active ledger so balances are correct — but simply deleting it would lose the original imported
record and its provenance, which matters for auditing and for understanding past imports.

## Decision Drivers

- Merged-away transactions must not affect active balances
- The original record (and its external ids) must be preserved
- Standard transaction queries should only return active transactions

## Considered Options

- **Hard delete the merged transaction** — simple, but loses history
- **Keep it active with a flag** — risks it leaking into balances/queries
- **Soft-delete + archive snapshot** — remove from active set, preserve a copy

## Decision Outcome

Chosen: on merge, the kept transaction absorbs the merged transaction's `ExternalIDs`; the merged
transaction is **GORM soft-deleted** (excluded from active queries) **and** a full snapshot is
written to the `MergedTransaction` archive table (`KeptTransactionID`, `OriginalTransactionID`, and
all original fields). Active gets (`GET /v1/transactions/{id}`) exclude it; the archive is read via
`GET /v1/mergedTransactions/{id}`. A startup migration backfills archives for any pre-existing
merged transactions.

### Pros

- Active balances and listings are clean (soft-deleted rows excluded)
- Full original record preserved for audit and history
- External ids transferred so re-imports still dedupe against the kept transaction

### Cons

- Two representations to keep coherent (soft-deleted row + archive snapshot)
- The snapshot duplicates fields and is frozen — schema changes to `Transaction` don't propagate
- Retrieval of merged records requires a separate endpoint/table

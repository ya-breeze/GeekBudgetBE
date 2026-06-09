# ADR-007: Empty AccountId as the "Unprocessed" Contract

## Status
Accepted

## Context and Problem Statement

Bank importers know the source account of a transaction (e.g. the checking account) but not the
category it should be booked against (groceries, rent, etc.). The system needs a way to represent
a transaction that is imported but not yet fully categorized, so the user (or a matcher) can
complete it later.

## Decision Drivers

- Imported transactions must be persisted immediately, before categorization
- The "needs categorization" state must be queryable and unambiguous
- The model should not require a separate "unprocessed transaction" table

## Considered Options

- **Empty `AccountId` on a movement** — a movement with `AccountId == ""` marks the transaction
  as unprocessed
- **A boolean flag** (e.g. `IsProcessed`) on the transaction
- **A separate unprocessed-transactions table** distinct from transactions

## Decision Outcome

Chosen: a transaction is **unprocessed** when any of its movements has an empty `AccountId`.
The same `Transaction` model and table serve both states. The unprocessed-transactions view filters
for transactions with an empty-account movement; matching/auto-matching fills that movement's
account with the matcher's `OutputAccountId`. Movement validation explicitly permits empty
`AccountId` while rejecting non-existent account ids.

### Pros

- One transaction model and table for both states; no duplication
- The state is derivable from data, no extra flag to keep consistent
- Conversion is a simple in-place fill of the empty account

### Cons

- It is an implicit convention — empty string carries semantic meaning that must be documented
- Validation has to special-case empty `AccountId` against the "account must exist" rule
- Every consumer that aggregates balances must be aware that unprocessed movements have no account

# ADR-005: Double-Entry Bookkeeping via Movements

## Status
Accepted

## Context and Problem Statement

A transaction needs to represent where money came from and where it went, support transfers
between accounts, currency exchange, and fees — not just a single signed amount against one
category.

## Decision Drivers

- Model real financial events (transfers, splits, fees, FX) accurately
- Make balances derivable by summing movements per account
- Support multi-currency within a single transaction

## Considered Options

- **Single-entry** — transaction has one `amount`, one `category`, one `account`
- **Double-entry via movements** — transaction owns a list of movements, each (account, currency,
  signed amount)

## Decision Outcome

Chosen: **double-entry**. A `Transaction` owns `Movements []Movement`; each movement references one
`AccountId`, one `CurrencyId`, and a signed `Amount`. A normal transaction has 2+ movements that
balance (e.g. -100 from "Cash", +100 to "Groceries"). Account balances are computed by aggregating
movements. The MCP `instructions` block documents this model for AI assistants.

### Pros

- Transfers, splits, fees, and FX are natural (just more movements)
- Balances are derived, not stored redundantly
- Matches accounting conventions; reconciliation has a sound basis

### Cons

- More complex than single-entry for the simplest "I spent $5" case
- The UI and importers must construct balanced movement sets
- Validation must enforce that every movement references valid account/currency
  (see ADR-007 for the empty-account exception)

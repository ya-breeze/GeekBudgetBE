# ADR-004: Use shopspring/decimal for All Money

## Status
Accepted

## Context and Problem Statement

A finance app cannot tolerate floating-point rounding error. Amounts must be exact through storage,
arithmetic, JSON serialization, and balance comparison across currencies.

## Decision Drivers

- Exact decimal arithmetic for balances and reconciliation
- Consistent representation across backend storage, API, and frontend
- Reliable equality and tolerance comparisons

## Considered Options

- **`float64`** — native, fast, but lossy
- **Integer minor units (cents)** — exact, but awkward across currencies and fractional rates
- **`github.com/shopspring/decimal`** — arbitrary-precision decimal type

## Decision Outcome

Chosen: **`shopspring/decimal`** for every monetary field (movement amounts, budget amounts,
reconciliation balances, currency rates).

Critical global setting: `decimal.MarshalJSONWithoutQuotes = true` in `cmd/main.go` so amounts
serialize as JSON **numbers**, not strings — the frontend depends on this. Comparisons use
`.Equal()` (not `==`) and tolerance checks use `.Sub().Abs().GreaterThan(tolerance)` (reconciliation
tolerance is `0.01`).

### Pros

- No rounding drift in balances or reconciliation
- One representation end to end
- Tolerance-based balance checks are explicit and correct

### Cons

- The `MarshalJSONWithoutQuotes` flag is global and easy to forget in a new entrypoint — if missed,
  the frontend receives strings and arithmetic silently breaks
- Tests must use `.Equal()`; naive `Equal()` matchers can fail on scale differences (`1.0` vs `1.00`)
- Slightly more verbose than native floats

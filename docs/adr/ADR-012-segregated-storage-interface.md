# ADR-012: Segregated Storage Interface with Split Implementation

## Status
Accepted

## Context and Problem Statement

All persistence goes through a single `Storage` abstraction. As the domain grew (accounts,
currencies, transactions, matchers, importers, reconciliation, budgets, notifications, audit, …),
a single monolithic interface and file became hard to navigate, test, and mock.

## Decision Drivers

- Keep the persistence contract readable as it grows
- Allow focused mocking/testing per domain
- Avoid import cycles and one giant file

## Considered Options

- **One large `Storage` interface in one file** — simple but unwieldy
- **Per-domain sub-interfaces composed into `Storage`, split across files** — interface segregation
- **Separate repository structs per domain** — more plumbing, multiple injection points

## Decision Outcome

Chosen: a composed `Storage` interface built from per-domain sub-interfaces (`UserStorage`,
`AccountStorage`, `CurrencyStorage`, `TransactionStorage`, `BankImporterStorage`, `MatcherStorage`,
`ReconciliationStorage`, `AuditLogStorage`, `SystemStorage`, …). Implementation is one `storage`
struct whose methods live in domain-specific files (`storage_account.go`, `storage_transaction.go`,
…) with shared helpers in `storage_common.go`. `WithContext` returns a context-scoped copy (used to
attach change-source for audit). A single generated mock backs tests.

### Pros

- The interface reads as a set of cohesive domain contracts
- Implementation files map 1:1 to domains; tests run at package level
- `WithContext` cleanly threads request/change-source context into audit logging
- Generic helpers (`performUpdate`, `recordAuditLog`) cut CRUD boilerplate

### Cons

- The composed interface is large (`//nolint:interfacebloat`) even if segments are small
- One concrete struct implements everything — not truly independent per-domain stores
- Contributors must know which `storage_*.go` file a method belongs in

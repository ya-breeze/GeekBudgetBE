# ADR-006: GORM AutoMigrate Instead of Versioned Migrations

## Status
Accepted

## Context and Problem Statement

The schema evolves as features are added. The app needs schema changes applied on startup without
operational ceremony, given a single-file SQLite database and a single-host deployment.

## Decision Drivers

- Minimal operational overhead for a self-hosted single-family app
- Additive schema changes are the common case (new tables, new columns)
- SQLite's limited DDL makes complex migrations awkward anyway (see ADR-001)

## Considered Options

- **GORM `AutoMigrate`** — derive and apply additive schema changes from struct tags on startup
- **Versioned migration files** (golang-migrate, goose) — explicit, ordered, reversible scripts
- **Hand-written SQL migrations**

## Decision Outcome

Chosen: **GORM `AutoMigrate`**, called in `storage.Open()` over the full model list. Where data
transformation is needed beyond what AutoMigrate does, a targeted Go function runs alongside it
(e.g. `migrateExistingMergedTransactions` backfills the merged-transactions archive). A separate
`kin-core` migration runs first for auth tables.

### Pros

- Schema is kept in sync with the Go models automatically on startup
- No migration files to author, order, or apply manually
- Additive changes (the common case) just work

### Cons

- Destructive or renaming changes (drop/rename column) are not handled — these need manual,
  one-off Go migration code
- No down-migrations / rollback story
- Schema history is implicit (in the structs + git) rather than an explicit ordered ledger
- Care needed: AutoMigrate will not remove or alter existing columns, so cleanups are manual

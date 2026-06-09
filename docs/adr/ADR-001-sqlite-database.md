# ADR-001: Use SQLite as the Database

## Status
Accepted

## Context and Problem Statement

GeekBudget needs persistent storage for accounts, currencies, transactions/movements, matchers,
bank importers, reconciliations, budgets, and audit logs. It is a personal finance app for a
single family, self-hosted, with low concurrency.

## Decision Drivers

- Zero operational overhead — no separate database process to run or upgrade
- File-based — the whole dataset is one file, trivial to back up
- Single-host deployment — no need for network database access
- Low concurrency — one family, a handful of users

## Considered Options

- **SQLite** — embedded, file-based relational database
- **PostgreSQL** — full-featured server-based RDBMS
- **MySQL/MariaDB** — server-based RDBMS

## Decision Outcome

Chosen: **SQLite**, accessed via GORM (`gorm.io/driver/sqlite`).

There is no scenario requiring high concurrent write throughput or horizontal scaling. SQLite's
single-file format integrates directly with the backup task (see ADR-013): backups are produced
with `VACUUM INTO` to a temp file, then archived as `geekbudget-backup-YYYY-MM-DD.tar.gz` alongside
bank-importer files.

### Pros

- No separate process to start, monitor, or upgrade
- Database file is included directly in the tar.gz backup
- `VACUUM INTO` gives an atomic, consistent snapshot without stopping the app
- GORM keeps the door open to another driver later

### Cons

- Not suitable for horizontal scaling or high concurrent write throughput
- Limited DDL (e.g. partial ALTER TABLE) constrains schema evolution (see ADR-006)
- `go-sqlite3` requires CGo, complicating cross-compilation

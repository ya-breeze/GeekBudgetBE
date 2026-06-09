# Architecture Decision Records

This directory records the significant architectural decisions behind GeekBudget — the *why*
behind choices that aren't obvious from reading the code. Each ADR captures the context, the
options considered, the decision, and its trade-offs.

The current-state behavior these decisions produce is described in the OpenSpec specs under
[`openspec/specs/`](../../openspec/specs/).

| ADR | Title |
|-----|-------|
| [001](ADR-001-sqlite-database.md) | Use SQLite as the Database |
| [002](ADR-002-openapi-first-codegen.md) | OpenAPI-First Development with Code Generation |
| [003](ADR-003-family-multitenancy-kin-core.md) | Family-Based Multi-Tenancy via kin-core |
| [004](ADR-004-decimal-for-money.md) | Use shopspring/decimal for All Money |
| [005](ADR-005-double-entry-movements.md) | Double-Entry Bookkeeping via Movements |
| [006](ADR-006-gorm-automigrate.md) | GORM AutoMigrate Instead of Versioned Migrations |
| [007](ADR-007-empty-accountid-unprocessed-contract.md) | Empty AccountId as the "Unprocessed" Contract |
| [008](ADR-008-matcher-confirmation-history-autotrust.md) | Matcher Trust via Rolling Confirmation History |
| [009](ADR-009-duplicate-detection-heuristic.md) | Duplicate Detection Heuristic |
| [010](ADR-010-archive-on-merge.md) | Archive-on-Merge for Deduplicated Transactions |
| [011](ADR-011-mcp-server-interface.md) | MCP Server as a First-Class Interface |
| [012](ADR-012-segregated-storage-interface.md) | Segregated Storage Interface with Split Implementation |
| [013](ADR-013-timer-based-background-tasks.md) | Timer-Based Background Tasks |

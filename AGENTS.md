# CLAUDE.md

## Project Overview

GeekBudget is a full-stack personal finance management application with a Go backend and Angular frontend. It supports multiple users, bank transaction imports, transaction matching, and budgeting.

See `.agent/rules/principal-architect-mode.md` before starting any work.

## Repository Structure

```
GeekBudgetBE/
├── api/                    # OpenAPI 3.0 specification (openapi.yaml)
├── backend/               # Go backend (REST API + web UI)
│   ├── cmd/               # CLI entry point and commands (Cobra)
│   ├── pkg/
│   │   ├── auth/          # JWT authentication
│   │   ├── bankimporters/ # FIO, KB, Revolut converters/fetchers
│   │   ├── config/        # Viper configuration
│   │   ├── database/      # GORM models, split SQLite storage implementation, migrations, mocks
│   │   │   ├── storage.go           # Composed Storage interface & struct definition
│   │   │   ├── storage_*.go         # Domain-specific implementations (user, account, etc.)
│   │   │   ├── storage_common.go    # Shared internal helpers
│   │   │   └── bulk_types.go        # Shared data structures for bulk operations
│   │   ├── generated/     # Auto-generated code from OpenAPI (DO NOT EDIT)
│   │   ├── server/        # HTTP server, middlewares, API handlers, web handlers, background jobs
│   │   └── utils/         # Utility functions
│   ├── test/              # Integration tests
│   └── webapp/            # Static assets and HTML templates
├── app/                   # Next.js 15 SPA (In Development - Future UI)
│   ├── app/               # Next.js App Router pages
│   ├── components/        # React components (shared + shadcn/ui)
│   ├── hooks/             # Custom React hooks
│   └── lib/               # API client, auth, utilities
├── frontend/              # Angular 20 SPA (Active - Current UI)
│   └── src/app/
│       ├── core/api/      # Generated API client
│       ├── features/      # Feature modules (accounts, transactions, matchers, etc.)
│       ├── layout/        # Layout components
│       └── shared/        # Shared components, pipes, directives
├── nginx/                 # Reverse proxy configuration
├── Makefile               # Build, test, lint, docker orchestration
└── docker-compose.yml     # Multi-service deployment
```

## Tech Stack

- **Backend:** Go 1.24+, Gorilla Mux, GORM, SQLite, JWT, Cobra/Viper, slog, shopspring/decimal
- **Frontend (Active):** Angular 20, React 19, TanStack Query, shadcn/ui, Tailwind CSS 3, TypeScript
- **Frontend (In Development):** Next.js 15, Angular Material, RxJS, Chart.js, SCSS
- **Testing:** Ginkgo/Gomega (backend BDD)
- **Code generation:** OpenAPI Generator (Go client/server + Angular client)
- **Formatting:** gofumpt (Go), Prettier + ESLint (frontend)
- **Linting:** golangci-lint (Go), ESLint with Angular rules (frontend)
- **Workflows:** See `.agent/workflows/` for specialized task guides (e.g., deduplication)
- **Investigation:** See `.agent/rules/database-investigation.md` for local DB querying guidelines
- **Financial Details:** See `.agent/rules/financial-data-handling.md` for Decimal usage rules
- **Bank Importers:** See `.agent/rules/bank-importers.md` for importer-specific dates and rules

## Common Commands

```bash
# Build everything
make build

# Run all tests (backend + frontend)
make test

# Run backend tests only
cd backend && go tool github.com/onsi/ginkgo/v2/ginkgo -r

# Run frontend tests only
cd frontend && npm run test -- --watch=false --browsers=ChromeHeadless

# Watch backend tests
make watch

# Start backend dev server (port 8080)
make run-backend

# Start frontend dev server (port 4200, proxies API to 8080)
make run-frontend

# Format and lint all code
make lint

# Full pipeline: build, test, validate OpenAPI, lint
make all

# Generate code from OpenAPI spec
make generate

# Generate mocks
make generate_mocks

# MCP Server (AI Assistant Integration)
make mcp-config          # Generate .mcp.json for Claude Desktop/Code
make mcp-server          # Run MCP stdio server for test@test.com

# Docker
make docker-build
make docker-up
make docker-down

# Database Investigation
sqlite3 geekbudget.db ".tables"
sqlite3 geekbudget.db ".header on" ".mode column" "SELECT * FROM transactions LIMIT 10;"
```

## Development Workflow

1. After making changes, run `make all` to build, test, validate, and lint.
2. For Go files, always run `gofumpt -w` on changed files.
3. Never manually edit files in `backend/pkg/generated/` -- they are auto-generated.
4. When changing API definitions, update `api/openapi.yaml` first, then `make generate`.
5. All database models must include `UserID` for multi-user isolation.
6. All API endpoints require JWT auth except `/v1/authorize`.
8. **Financial Accuracy**: Always use `decimal.Decimal` (from `github.com/shopspring/decimal`) for money. In tests, use `.Equal()` instead of `==`. In the frontend, wrap amounts in `Number()` for safety.
9. **API Strictness**: When updating transactions, ensure the request body does NOT contain an `id` or other `Entity` fields. The backend decodes directly into `TransactionNoId` (or similar interface) and will fail with `json: unknown field "id"` if extra fields are present. Strip these fields in the frontend service or component before sending.
10. **For Next.js frontend**: API calls should use the axios client in `app/lib/api/client.ts` with JWT interceptor. Use TanStack Query hooks in `app/lib/api/hooks/` for data fetching.
11. **Refactoring Workflow**: Always ensure existing code is covered by tests *before* refactoring or fixing bugs. If tests are missing, add them first (e.g., `storage_images_test.go`). When refactoring `Storage`, split domain logic into separate `storage_DOMAIN.go` files and use Interface Segregation for the main `Storage` interface.

## Testing

- Backend uses Ginkgo BDD framework with Gomega matchers. Tests live alongside source files (`*_test.go`) and in `backend/test/` for integration tests.
- Frontend uses Karma/Jasmine. Tests are `*.spec.ts` files alongside components.
- Legacy Angular frontend uses Karma/Jasmine. Tests are `*.spec.ts` files alongside components.
- Use `test@test.com` / `test` credentials for browser testing.
- If `make run-backend` fails with "address already in use", the backend is already running.
- If `make run-frontend` fails with "Port 4200 is already in use", the frontend is already running. Use `make run-app` for the Next.js frontend (port 3000).

## Deduplication & Archiving Flow
 
 1. **Background Task:** `StartDuplicateDetection` runs every 24 hours (or manually via `DuplicateDetectionCommand`), scanning transactions from the last 30 days. Uses `models.DuplicateReason` constant.
 2. **Identification:** Uses `common.IsDuplicate` to find transactions with similar dates (±2 days) and amounts.
 3. **Different Sources:** Only flags transactions if they have different `ExternalIDs` (indicating different import sources).
 4. **Duplicate Linking:** Pairwise relationships are stored in the `TransactionDuplicate` junction table. Bidirectional links (T1↔T2) allow efficient retrieval.
 5. **User Resolution:**
    - **Dismissal:** Setting `DuplicateDismissed = true` clears all links and prevents re-flagging.
    - **Merging:** `POST /v1/transactions/merge` transfers external IDs to the "kept" transaction and performs a GORM soft-delete on the other.
    - **Archiving:** Merged transactions are moved to a separate archive table.
    - **Retrieval:** Use `GET /v1/mergedTransactions/{id}` to fetch these archived records. The standard `GET /v1/transactions/{id}` only returns active ones.
 6. **Synchronized Cleanup:** The storage layer (`ClearDuplicateRelationships`) automatically removes `models.DuplicateReason` from linked transactions if they have no other duplicate links remaining.

## Reconciliation Flow

1. **Status Retrieval**: `GET /v1/reconciliation/status` returns balance details, delta, and flags for all accounts.
2. **Tolerance**: Minor discrepancies up to `common.ReconciliationTolerance` (0.01) are handled as "matching".
3. **Blocking**: Manual reconciliation is **blocked** if `hasUnprocessedTransactions` is true. This ensures the App Balance is finalized before being compared to the Bank Balance.
4. **Stale Balances**: If `hasTransactionsAfterBankBalance` is true, a warning ⚠️ is shown in the UI. This happens when the system detects transactions with a date newer than the bank balance timestamp.
5. **Timestamping**: Bank balances include a `lastUpdatedAt` timestamp derived from statement metadata (e.g., Fio's `DateEnd`) or the newest transaction in an import.
6. **UI Tooltips**: Disabled reconciliation buttons have tooltips (wrapped in `<span>`) explaining exactly why reconciliation is unavailable (e.g., large delta or unprocessed transactions).
7. **Manual Reconciliation**: 
   - **Auto-Balance**: When enabling manual reconciliation, the system defaults to using the *current* `AppBalance` as the starting bank balance, avoiding manual input.
   - **Display**: For accounts without bank importers, manual reconciliation records are used to populate "Bank Balance" and "Balance Date" columns in the UI, mimicking a bank feed for consistency.
8. **Performance & Batching**: Status retrieval is optimized via `GetBulkReconciliationData`. It fetches all accounts, latest reconciliations, and transactions in a single pass to avoid N+1 queries.

## Code Patterns

> **See [backend/README.md](backend/README.md)** for detailed backend patterns and examples.

**Quick reference:**
- **Service layer:** `ServiceImpl` structs with `logger` and `db` fields, methods take `context.Context` and `userID`.
- **API handlers:** Extract `userID` from context via `ctx.Value(common.UserIDKey)`, return `goserver.ImplResponse`.
- **Web handlers:** Use `r.ValidateUserID()` for auth, `utils.CreateTemplateData()` for template data, `tmpl.ExecuteTemplate()` for rendering.
- **Storage interface:** All DB operations go through the `database.Storage` interface. The implementation is split into domain-specific files (`storage_user.go`, `storage_account.go`, etc.) and the main `Storage` interface is a composition of fine-grained interfaces (`UserStorage`, `AccountStorage`, etc.). Mock implementation is in `database/mocks/` for testing.
- **Bank importers:** Implement the `Importer` interface in `pkg/bankimporters/`.

## Key Environment Variables

| Variable | Purpose | Default |
|---|---|---|
| `GB_USERS` | User credentials (bcrypt) | - |
| `GB_DBPATH` | SQLite database path | `/data/geekbudget.db` |
| `GB_JWT_SECRET` | JWT signing secret | - |
| `GB_SESSIONSECRET` | Session secret | - |
| `GB_COOKIESECURE` | HTTPS-only cookies | `true` |
| `GB_DISABLEIMPORTERS` | Disable bank importers | `false` |
| `GB_PORT` | API server port | `8080` |
| `GB_ALLOWEDORIGINS` | CORS origins | - |

## JSON Querying Patterns

When querying records with JSON columns like `movements` (in `transactions`) or `bank_info` (in `accounts`), use SQLite's native JSON functions for precision:

-   **SQLite Example**: `SELECT * FROM transactions, json_each(transactions.movements) WHERE json_extract(json_each.value, '$.accountId') = '...';`
-   **GORM Joins**: Use `Joins("CROSS JOIN json_each(table.column)")` for array fields.
-   **GORM Grouping**: **CRITICAL**: Use `.Group("table.id")` instead of `.Distinct("id")` when joining with `json_each`. `Distinct("id")` restricts the selected columns to *only* the ID, which breaks full record retrieval/Save operations.
-   **JSON Paths**:
    *   Transactions: `$.accountId`, `$.currencyId`, `$.amount`
    *   Accounts: `$.balances[*].currencyId`, `$.balances[*].openingBalance`

## Database Development Tips

1.  **SQLite & time.Time**: When querying the latest record by date, do **NOT** use `SELECT MAX(created_at)`. SQLite returns this as a string, which GORM cannot scan into a `time.Time` destination.
    -   ❌ `db.Select("MAX(date)").Scan(&t)`
    -   ✅ `db.Order("date DESC").First(&record)`
2.  **Region Markers**: This project uses `#region` and `#endregion` comments. Ensure they are balanced and descriptive.
    -   Use `// #region Name` and `// #endregion Name`.

## Next.js Frontend Development (app/)

**Status**: Phase 1 complete. Angular remains default until all phases finished.

> **See [app/README.md](app/README.md)** for current status, commands, and tech stack.
> **See [app/PLAN.md](app/PLAN.md)** for detailed implementation plan and architecture.

**Quick Start**:
```bash
make dev-app          # Backend + Next.js (port 3000)
make run-app          # Next.js only
make build-app        # Build Next.js
make lint-app         # Lint Next.js
```

**Test Credentials**: `test@test.com` / `test`

# CLAUDE.md

## Project Overview

GeekBudget is a full-stack personal finance management application with a Go backend and Angular frontend. It supports multiple users, bank transaction imports, transaction matching, and budgeting.

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
│   │   ├── database/      # GORM models, SQLite storage interface, migrations, mocks
│   │   ├── generated/     # Auto-generated code from OpenAPI (DO NOT EDIT)
│   │   ├── server/        # HTTP server, middlewares, API handlers, web handlers, background jobs
│   │   └── utils/         # Utility functions
│   ├── test/              # Integration tests
│   └── webapp/            # Static assets and HTML templates
├── frontend/              # Angular 20 SPA
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

- **Backend:** Go 1.24+, Gorilla Mux, GORM, SQLite, JWT, Cobra/Viper, slog
- **Frontend:** Angular 20, Angular Material, RxJS, Chart.js, SCSS
- **Testing:** Ginkgo/Gomega (backend BDD), Karma/Jasmine (frontend)
- **Code generation:** OpenAPI Generator (Go client/server + Angular client)
- **Formatting:** gofumpt (Go), Prettier + ESLint (frontend)
- **Linting:** golangci-lint (Go), ESLint with Angular rules (frontend)

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

# Docker
make docker-build
make docker-up
make docker-down
```

## Development Workflow

1. After making changes, run `make all` to build, test, validate, and lint.
2. For Go files, always run `gofumpt -w` on changed files.
3. Never manually edit files in `backend/pkg/generated/` -- they are auto-generated.
4. When adding new API endpoints, update `api/openapi.yaml` first, then `make generate`.
5. All database models must include `UserID` for multi-user isolation.
6. All API endpoints require JWT auth except `/v1/authorize`.
7. **For Next.js frontend**: Never import from `new-frontend/src/lib/api/generated/` directly. Always use custom hooks from `new-frontend/src/lib/api/hooks/` which wrap the generated code and won't be overwritten.

## Testing

- Backend uses Ginkgo BDD framework with Gomega matchers. Tests live alongside source files (`*_test.go`) and in `backend/test/` for integration tests.
- Frontend uses Karma/Jasmine. Tests are `*.spec.ts` files alongside components.
- Use `test@example.com` / `test` credentials for browser testing.
- If `make run-backend` fails with "address already in use", the backend is already running.
- If `make run-frontend` fails with "Port 4200 is already in use", the frontend is already running.

## Code Patterns

- **Service layer:** `ServiceImpl` structs with `logger` and `db` fields, methods take `context.Context` and `userID`.
- **API handlers:** Extract `userID` from context via `ctx.Value(common.UserIDKey)`, return `goserver.ImplResponse`.
- **Web handlers:** Use `r.ValidateUserID()` for auth, `utils.CreateTemplateData()` for template data, `tmpl.ExecuteTemplate()` for rendering.
- **Storage interface:** All DB operations go through the `database.Storage` interface; mock implementation in `database/mocks/` for testing.
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

# GeekBudget Backend

> **For comprehensive developer documentation**, see [`CLAUDE.md`](../CLAUDE.md) in the project root.

## Overview

Go 1.24+ backend providing RESTful APIs and web interface for personal finance management. Supports multiple users, bank imports, transaction matching, and budgeting.

### Architecture

- **API Layer**: OpenAPI 3.0 generated endpoints (`api/openapi.yaml`)
- **Web Layer**: Go templates and web handlers
- **Database**: SQLite with GORM
- **Authentication**: JWT-based
- **Bank Importers**: FIO, KB, Revolut

### Directory Structure

```
cmd/                    - CLI entry points (Cobra)
pkg/
├── auth/              - JWT authentication
├── bankimporters/     - Bank-specific importers (FIO, KB, Revolut)
├── config/            - Viper configuration
├── database/          - GORM models, split SQLite storage implementation (storage.go, storage_*.go)
├── generated/         - Auto-generated from OpenAPI (DO NOT EDIT)
├── server/            - HTTP server, middleware, handlers
└── utils/             - Utilities
test/                  - Integration tests
webapp/                - Static assets and templates
```

## Development Workflow

### Code Generation

```bash
# Regenerate API code from OpenAPI spec
make generate

# NEVER manually edit pkg/generated/ - it's auto-generated
# Always update api/openapi.yaml first, then regenerate
```

### Testing

```bash
# Run backend tests with Ginkgo
cd backend && go tool github.com/onsi/ginkgo/v2/ginkgo -r

# Or use the Makefile
make test
```

### Formatting

```bash
# Format changed files
go tool mvdan.cc/gofumpt -w <file>

# Always run before committing
make all
```

## Backend-Specific Code Patterns

### Service Layer Pattern

All business logic lives in service structs with logger and database dependencies:

```go
type ServiceImpl struct {
    logger *slog.Logger
    db     database.Storage
}

func (s *ServiceImpl) Method(ctx context.Context, userID string, params) (result, error) {
    // Validate parameters
    // Perform business logic
    // Database operations
    // Return results
}
```

### API Handler Pattern

API handlers extract userID from context and return `goserver.ImplResponse`:

```go
func (s *APIServiceImpl) HandleEndpoint(ctx context.Context, request Request) (goserver.ImplResponse, error) {
    userID, ok := ctx.Value(constants.UserIDKey).(string)
    if !ok {
        return goserver.Response(500, nil), nil
    }

    // Business logic here

    return goserver.Response(200, result), nil
}
```

### Web Handler Pattern

Web handlers validate userID, create template data, and execute templates:

```go
func (r *WebAppRouter) pageHandler(w http.ResponseWriter, req *http.Request) {
    tmpl, err := r.loadTemplates()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    userID, err := r.ValidateUserID(tmpl, w, req)
    if err != nil {
        return
    }

    data := utils.CreateTemplateData(req, "page_name")
    // Add page-specific data

    if err := tmpl.ExecuteTemplate(w, "template.tpl", data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

## Bank Importers

Each bank has its own importer in `pkg/bankimporters/`:

- **Interface**: Implement `Importer` interface for new banks
- **File Formats**: Handle CSV, Excel, etc.
- **Conversion**: Transform bank-specific formats to internal transaction model
- **Error Handling**: Validate and report errors clearly

**Supported Banks**: FIO, KB (Komerční banka), Revolut

See [`.agent/rules/bank-importers.md`](../.agent/rules/bank-importers.md) for importer-specific guidelines.

## Common Tasks

### Adding New API Endpoint

1. Update `api/openapi.yaml` with endpoint definition
2. Run `make generate` to regenerate code
3. Implement handler in `pkg/server/api/`
4. Add business logic and database operations
5. Test with `make test`

### Adding New Database Model

1. Create model struct in `pkg/database/models/`
2. Add migration in `pkg/database/migration.go`
3. Update `database.Storage` interface if needed
4. Generate mocks: `make generate_mocks`
5. Implement CRUD operations

### Adding Bank Importer

1. Create new file in `pkg/bankimporters/`
2. Implement `Importer` interface
3. Handle bank-specific data format
4. Add tests for conversion logic

## Key Dependencies

- **Gorilla Mux**: HTTP routing
- **GORM**: ORM for SQLite
- **JWT-go**: Authentication tokens
- **Cobra/Viper**: CLI and configuration
- **slog**: Structured logging
- **shopspring/decimal**: Financial precision
- **Ginkgo/Gomega**: BDD testing

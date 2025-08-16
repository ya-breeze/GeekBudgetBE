# Copilot Instructions for GeekBudget Backend

## Project Overview
GeekBudget is a personal finance management application backend written in Go which supports multiple users.
It provides RESTful APIs and a web interface for managing personal financial data, including accounts,
transactions, bank imports, and budgeting features.

## Architecture & Structure

### Core Components
- **API Layer**: RESTful endpoints generated from OpenAPI 3.0 specification (`api/openapi.yaml`)
- **Web Layer**: HTML templates and web handlers for browser interface
- **Database Layer**: SQLite with GORM for data persistence
- **Authentication**: JWT-based authentication system
- **Bank Importers**: Automated transaction import from various banks (FIO, KB, Revolut)

### Directory Structure
```
cmd/                    - CLI application entry points
pkg/
├── auth/              - Authentication & JWT handling
├── bankimporters/     - Bank-specific transaction importers
├── config/            - Application configuration
├── constants/         - Application constants
├── database/          - Database models, migrations, storage interface
├── generated/         - Auto-generated API client/server code
├── server/            - HTTP server, middleware, API handlers
├── utils/             - Utility functions
└── webapp/            - Web interface handlers and templates
test/                  - Integration and unit tests
webapp/                - Static assets and HTML templates
```

## Development Guidelines

### General
- Project allows multiple users to manage their financial data securely.
- Project stores each financial event as a transaction linked to a user.
- Project supports automatic import of transactions from various sources.
- Project allows to define "matchers" for categorizing transactions based on rules.

### Code Generation
- The project uses OpenAPI Generator for API client/server code
- Run `make generate` to regenerate code from `api/openapi.yaml`
- Never manually edit files in `pkg/generated/` - they are auto-generated
- When adding new endpoints, update `api/openapi.yaml` first
- Run code generation after API changes

### API Development
- Follow RESTful principles for new endpoints
- All APIs require authentication except `/v1/authorize`
- Use proper HTTP status codes and error handling
- Validate input parameters and request bodies
- Include proper OpenAPI documentation for new endpoints

### WEB UI
- Use golang templates for rendering HTML
- Follow existing patterns in `pkg/server/webapp/`
- Use `webapp/templates/` for HTML templates
- Serve static assets from `webapp/static/`
- Implement CSRF protection for forms
- Ensure user sessions are managed securely
- User token is passed in cookies for web requests

### Database & Models
- Use GORM for database operations
- Follow the existing model patterns in `pkg/database/models/`
- All models should include user isolation (UserID field)
- Use transactions for multi-step operations
- Implement proper database migrations

### Error Handling
- Use structured logging with slog
- Return appropriate HTTP status codes
- Provide meaningful error messages
- Log errors with context for debugging

### Authentication & Security
- All API endpoints (except auth) require JWT authentication
- Use middleware for authentication checks
- Validate user permissions for data access
- Ensure users can only access their own data

### Testing
- Write tests for new functionality
- Use the existing test patterns in `test/` directory
- Mock database operations for unit tests
- Include integration tests for API endpoints

## Coding Patterns

### Service Layer Pattern
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

### Handler Pattern (API)
```go
func (s *APIServiceImpl) HandleEndpoint(ctx context.Context, request Request) (goserver.ImplResponse, error) {
    userID, ok := ctx.Value(common.UserIDKey).(string)
    if !ok {
        return goserver.Response(500, nil), nil
    }
    
    // Business logic here
    
    return goserver.Response(200, result), nil
}
```

### Handler Pattern (Web)
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
- Each bank has its own converter in `pkg/bankimporters/`
- Implement the `Importer` interface for new banks
- Handle different file formats (CSV, Excel, etc.)
- Convert bank-specific formats to internal transaction format
- Include proper error handling and validation

## Configuration
- Use Viper for configuration management
- Support environment variables and config files
- Store sensitive data (JWT secrets, DB paths) in config
- Validate configuration on startup

## Key Technologies & Dependencies
- **Go 1.24+**: Primary language
- **Gorilla Mux**: HTTP routing
- **GORM**: ORM for database operations
- **SQLite**: Database engine
- **JWT**: Authentication tokens
- **OpenAPI Generator**: API code generation
- **Cobra**: CLI framework
- **Viper**: Configuration management
- **Slog**: Structured logging

## Common Tasks

### Adding New API Endpoint
1. Update `api/openapi.yaml` with new endpoint definition
2. Regenerate code
3. Implement service method in appropriate `pkg/server/api/` file
4. Add business logic and database operations
5. Test the endpoint

### Adding New Web Page
1. Create handler in `pkg/server/webapp/`
2. Add route in webapp router
3. Create HTML template in `webapp/templates/`
4. Add any required static assets

### Adding New Database Model
1. Create model struct in `pkg/database/models/`
2. Add migration in `pkg/database/migration.go`
3. Update storage interface if needed
4. Implement CRUD operations

### Bank Importer
1. Create converter in `pkg/bankimporters/`
2. Implement required interface methods
3. Handle bank-specific data format
4. Add tests for conversion logic

## Performance Considerations
- Use database connections efficiently
- Implement proper caching where appropriate
- Handle file uploads for bank import efficiently
- Use background processing for long-running tasks
- Monitor memory usage for large datasets

## Security Best Practices
- Validate all user inputs
- Use parameterized queries to prevent SQL injection
- Implement rate limiting for sensitive endpoints
- Ensure proper session management
- Log security-relevant events
- Keep dependencies updated for security patches

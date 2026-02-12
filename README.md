# GeekBudget - Personal Finance Management

A full-stack personal finance management application with Go backend, Angular frontend (active), and Next.js frontend (in development).

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose
- Make

### Local Development

```bash
# Install dependencies
make install

# Build everything
make build

# Run backend (port 8080)
make run-backend

# Run Angular frontend (port 4200)
make run-frontend

# Or run Next.js frontend (port 3000, in development)
make dev-app
```

### Docker Deployment

```bash
# Build and start all services
make compose

# Access at http://localhost
```

## Architecture

- **Backend**: Go 1.24+ REST API with SQLite database
- **Frontend**: Angular 20 SPA (active) + Next.js 15 (in development)
- **Reverse Proxy**: Nginx
- **Testing**: Ginkgo/Gomega (backend), Karma/Jasmine (Angular)

## Key Features

- Multi-user finance tracking with JWT authentication
- Bank transaction import (FIO, KB, Revolut)
- Automated transaction matching and categorization
- Budget planning and reconciliation
- Duplicate detection and merging

## Documentation

- **[CLAUDE.md](CLAUDE.md)** - Comprehensive developer guide (tech stack, commands, patterns, workflows)
- **[backend/README.md](backend/README.md)** - Backend-specific development guide
- **[app/README.md](app/README.md)** - Next.js frontend status and roadmap

## License

See LICENSE file for details.

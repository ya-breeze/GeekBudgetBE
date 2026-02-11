# GeekBudget - Personal Finance Management

A full-stack personal finance management application with Go backend and Angular frontend.

## Project Structure

```
GeekBudgetBE/
├── api/                    # OpenAPI specifications
├── backend/               # Go backend application
│   ├── cmd/              # Application entry points
│   ├── pkg/              # Go packages
│   ├── test/             # Backend tests
│   └── webapp/           # Static web assets
├── frontend/             # Angular frontend (active)
│   ├── src/              # Angular source code
│   ├── public/           # Static assets
│   └── dist/             # Built frontend (generated)
├── app/                  # Next.js frontend (in development)
│   ├── app/              # Next.js App Router pages
│   ├── components/       # React components
│   ├── lib/              # API client, auth, utilities
│   └── hooks/            # Custom React hooks
├── nginx/                # Nginx reverse proxy configuration
├── docker-compose.yml    # Docker orchestration
└── Makefile             # Build and development commands
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose
- Make

### Development Setup

1. **Install dependencies:**
   ```bash
   make install
   ```

2. **Build the application:**
   ```bash
   make build
   ```

3. **Run with Docker Compose:**
   ```bash
   make compose
   ```

4. **Access the application:**
   - Web Interface: http://localhost
   - API Documentation: http://localhost/v1/docs (if available)

### Development Commands

- `make build` - Build both backend and frontend
- `make test` - Run all tests
- `make lint` - Lint code
- `make run-backend` - Run backend in development mode
- `make run-frontend` - Run frontend in development mode
- `make docker-up` - Start Docker containers
- `make docker-down` - Stop Docker containers
- `make docker-logs` - View container logs
- `make clean` - Clean build artifacts

## Architecture

The application follows a microservices architecture with:

- **Backend**: Go REST API server with SQLite database
- **Frontend**: Angular 20 SPA with Angular Material UI
- **Reverse Proxy**: Nginx for routing and load balancing
- **Containerization**: Docker for consistent deployment

> **Note**: A new Next.js 15 frontend is in development in the `app/` directory. Use `make dev-app` to run it.

## Features

- Personal finance tracking
- Transaction management
- Account management
- Budget planning
- Import/export functionality
- User authentication
- Responsive web interface

## Configuration

Environment variables can be set in `.env` file or passed to Docker Compose:

- `GB_USERS` - User credentials
- `GB_DBPATH` - Database file path
- `GEEKBUDGET_DATA_PATH` - Data volume mount path

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

See LICENSE file for details.

# Robotics Management Platform — API Server

Production-grade Go REST API built with Gin framework and clean architecture for enterprise robotics and management applications.

## Tech Stack

- **Go 1.23** (LTS)
- **Gin v1.10.0** — High-performance HTTP web framework
- **gin-contrib/cors v1.7.2** — CORS middleware for Gin
- **pgx v5.7.1** — PostgreSQL driver with connection pooling
- **Viper v1.19.0** — Configuration (YAML + env var overrides)
- **slog** — Structured logging (Go stdlib)
- **golang-migrate v4.18.1** — Database migrations
- **google/uuid v1.6.0** — Request ID generation
- **golang.org/x/time v0.8.0** — Token bucket rate limiting
- **Docker Compose** — Local development environment

## Architecture

```
HTTP Request
  │
  ▼
┌─────────────────────────────────────────────────────────┐
│  Gin Middleware Chain                                     │
│  Recovery → RequestID → Logging → CORS → RateLimit       │
│                              [Auth per-route]            │
└─────────────────────────────────────────────────────────┘
  │
  ▼
Handler  →  Service  →  Repository  →  PostgreSQL
 (Gin)    (Business)   (Data Access)
```

All API responses follow a standard JSON envelope format via interceptor helper functions:

```json
{
  "success": true,
  "message": "OK",
  "data": { },
  "errors": null,
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2026-02-13T10:30:00Z"
}
```

### Key Gin Patterns

- **`gin.New()`** instead of `gin.Default()` — explicit middleware control
- **`gin.SetMode()`** — auto-configured based on `APP_ENV` (debug/test/release)
- **`NoRoute` / `NoMethod`** — 404/405 responses in standard envelope
- **`c.ShouldBindJSON()`** — request body binding with validation
- **`c.AbortWithStatusJSON()`** — middleware rejection in standard envelope
- **Route groups** — `/api/v1` with nested auth-protected sub-groups

## Prerequisites

- [Go 1.23+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Make](https://www.gnu.org/software/make/) (optional, for convenience commands)

## Getting Started

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd my-application
```

### 2. Install Go dependencies

```bash
go mod tidy
```

### 3. Start PostgreSQL

```bash
docker compose -f deployments/docker/docker-compose.yml up -d postgres
```

Wait a few seconds for the database to be ready. You can check with:

```bash
docker compose -f deployments/docker/docker-compose.yml ps
```

### 4. Run database migrations

```bash
go run ./cmd/migration -direction=up
```

Or with Make:

```bash
make migrate-up
```

### 5. Start the API server

```bash
go run ./cmd/api
```

Or with Make:

```bash
make run
```

The server starts on `http://localhost:3000` by default.

### 6. Verify it works

```bash
# Liveness probe
curl http://localhost:3000/ping

# Health check (includes DB status)
curl http://localhost:3000/health

# List users (requires auth header)
curl -H "Authorization: Bearer test-token" http://localhost:3000/api/v1/users

# Create a user
curl -X POST http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer test-token" \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","email":"john@example.com","full_name":"John Doe","role":"operator"}'

# Get user by ID
curl -H "Authorization: Bearer test-token" http://localhost:3000/api/v1/users/1
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/ping` | No | Liveness probe |
| GET | `/health` | No | Service + database health check |
| GET | `/api/v1/users` | Yes | List users (paginated, filterable) |
| POST | `/api/v1/users` | Yes | Create a new user |
| GET | `/api/v1/users/:id` | Yes | Get user by ID |
| PUT | `/api/v1/users/:id` | Yes | Update user |
| DELETE | `/api/v1/users/:id` | Yes | Delete user |

### Query Parameters for `GET /api/v1/users`

| Param | Type | Description |
|-------|------|-------------|
| `role` | string | Filter by role: `admin`, `operator`, `viewer` |
| `is_active` | bool | Filter by active status: `true` or `false` |
| `limit` | int | Page size (default 20, max 100) |
| `offset` | int | Pagination offset |

## Configuration

Configuration is loaded in this order (later overrides earlier):

1. `config/config.yaml` — Base defaults
2. `config/config.{APP_ENV}.yaml` — Environment overlay (dev, staging, prod)
3. Environment variables — Prefixed with `APP_`, underscores replace dots

### Key environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment name (`dev`/`staging`/`prod`) | `dev` |
| `APP_SERVER_PORT` | Server port | `3000` |
| `APP_DATABASE_HOST` | PostgreSQL host | `localhost` |
| `APP_DATABASE_PORT` | PostgreSQL port | `5432` |
| `APP_DATABASE_USER` | Database user | `app_user` |
| `APP_DATABASE_PASSWORD` | Database password | `changeme` |
| `APP_DATABASE_DBNAME` | Database name | `robotics_mgmt` |
| `APP_LOG_LEVEL` | Log level (debug/info/warn/error) | `info` |
| `APP_LOG_FORMAT` | Log format (json/text) | `json` |

### Gin Mode

Gin mode is automatically set based on `APP_ENV`:

| APP_ENV | Gin Mode |
|---------|----------|
| `dev` / `development` | `debug` |
| `test` | `test` |
| anything else | `release` |

## Project Structure

```
my-application/
├── cmd/
│   ├── api/main.go              # API server entry point
│   └── migration/main.go        # Database migration tool
├── internal/
│   ├── api/
│   │   ├── handler/             # Gin HTTP handlers
│   │   ├── interceptor/         # Response envelope helpers (Success/Fail/Abort)
│   │   ├── middleware/           # Auth, CORS, logging, rate limit, recovery
│   │   ├── request/             # Request DTOs
│   │   ├── response/            # Response DTOs
│   │   └── router/              # Gin engine + route registration
│   ├── domain/                  # Domain entities and errors
│   ├── repository/              # Data access interfaces + PostgreSQL impl
│   └── service/                 # Business logic
├── pkg/
│   ├── database/                # PostgreSQL connection pool
│   └── logger/                  # Structured logging setup
├── config/                      # YAML configuration files
├── migrations/                  # SQL migration files
├── deployments/docker/          # Dockerfile + docker-compose
├── Makefile                     # Build automation
└── go.mod
```

## Make Commands

```bash
make build          # Build API and migration binaries to bin/
make run            # Build and run the API server
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make lint           # Run golangci-lint
make tidy           # Run go mod tidy
make migrate-up     # Run database migrations (up)
make migrate-down   # Roll back database migrations
make docker-up      # Start all services via Docker Compose
make docker-down    # Stop all Docker Compose services
make docker-build   # Build Docker images
make clean          # Remove build artifacts
```

## Running with Docker (Full Stack)

To run both PostgreSQL and the API server in Docker:

```bash
# Build and start everything
docker compose -f deployments/docker/docker-compose.yml up -d

# View logs
docker compose -f deployments/docker/docker-compose.yml logs -f api

# Stop everything
docker compose -f deployments/docker/docker-compose.yml down
```

## Database Migrations

Migrations live in the `migrations/` directory and follow the naming convention:

```
NNNNNN_description.up.sql    # Apply migration
NNNNNN_description.down.sql  # Roll back migration
```

```bash
# Apply all pending migrations
go run ./cmd/migration -direction=up

# Roll back all migrations
go run ./cmd/migration -direction=down

# Apply N steps
go run ./cmd/migration -direction=up -steps=1

# Roll back N steps
go run ./cmd/migration -direction=down -steps=1
```

## License

See [LICENSE](LICENSE) for details.

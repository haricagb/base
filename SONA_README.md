# SONA — Multi-Tenant IoT Elderly Care Platform

SONA is a multi-tenant IoT ecosystem designed for elderly care, integrating CoCo companion robots with a complex web of human relationships (caregivers, family members, tech admins). It provides real-time robot management, care incident tracking, family story sharing, and HIPAA-compliant conversation logging.

## Architecture Overview

```
                         Envoy Proxy (:80/:443)
                        /         |          \
  enterprise.*.dev     family.*.dev      /graphql  /api
       |                    |               |         \
   Svelte App         Flutter/API      Hasura v2    Go API
   (Admin UI)         (Family App)     (:8080)      (:3000)
                                           |           |
                                     PostgreSQL (:5432)
                                     (Multi-Tenant)
                                           |
                                    Firebase Auth
                                    (User Sync)
                                           |
                                    OTEL Collector
                                    (PII-Stripped Logs)
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| **Envoy** | 80, 9901 | Reverse proxy — routes by subdomain/path |
| **Hasura GraphQL** | 8080 | Auto-generated GraphQL API with RBAC |
| **Go API** | 3000 | Custom business logic (auth, robot lifecycle, HIPAA logging) |
| **PostgreSQL** | 5432 | Primary data store (multi-tenant) |
| **OTEL Collector** | 4317, 4318 | Telemetry collection with PII stripping |

## Data Model

### Tables

```
enterprises
├── users (FK: enterprise_id)
│   └── caregiver_residents (junction: caregiver ↔ resident, max 5)
├── residents (FK: enterprise_id)
│   ├── stories (FK: resident_id, author_id)
│   ├── incidents (FK: resident_id, reporter_id)
│   └── robot_sessions (FK: resident_id, robot_id)
└── robots (FK: enterprise_id, assigned_resident_id)
```

### Key Schema Decisions

- **Multi-tenancy**: All core tables have `enterprise_id` for tenant isolation
- **HIPAA compliance**: Sensitive fields use `_encrypted` suffix (medical_notes, conversation_summary, story content)
- **Max-5 rule**: Database trigger enforces max 5 residents per caregiver
- **Role-based**: Users have one of 5 roles: `mta`, `eta`, `caregiver`, `family`, `robot`

## Roles & Access Control (RBAC)

Implemented via Hasura row-level permissions + Firebase Custom Claims.

| Role | Full Name | Scope | Capabilities |
|------|-----------|-------|-------------|
| **mta** | Master Tech Admin | Global | Full CRUD on all tables. Manage enterprises, assign robot stock. |
| **eta** | Enterprise Tech Admin | Own enterprise | Manage users, residents, robots within their enterprise. Provision robots to residents. |
| **caregiver** | Caregiver | Assigned residents (max 5) | View assigned residents. Log incidents. Read stories and robot session summaries. |
| **family** | Family Member | Their relative | Write stories/incidents for their specific resident. View their own submissions. |
| **robot** | CoCo Robot (CC) | Own record | Machine-token auth. Read shared history (stories) for assigned resident. Log conversation sessions. |

### JWT Claims Mapping

The Go backend issues JWTs with these claims, mapped to Hasura session variables:

| JWT Claim | Hasura Session Variable | Description |
|-----------|------------------------|-------------|
| `user_id` | `X-Hasura-User-Id` | User's database ID |
| `role` | `X-Hasura-Default-Role` | One of: mta, eta, caregiver, family, robot |
| `enterprise_id` | `X-Hasura-Enterprise-Id` | Tenant scope for eta/caregiver/family |
| `robot_id` | `X-Hasura-Robot-Id` | Robot's database ID (robot role only) |

## Project Structure

```
my-application/
├── cmd/api/                    # Go HTTP server entry point
├── cmd/migration/              # Database migration CLI
├── config/                     # YAML configs (base, dev, prod, staging)
├── deployments/
│   ├── docker/                 # Dockerfile + docker-compose.yml
│   ├── envoy/                  # Envoy proxy configuration
│   ├── otel/                   # OpenTelemetry collector config
│   ├── helm/                   # Kubernetes Helm charts
│   └── kubernetes/             # Kustomize manifests
├── firebase/
│   ├── firebase-config.example.json
│   └── functions/sync-user/    # Cloud Function for user sync
├── hasura/
│   ├── config.yaml             # Hasura CLI config
│   └── metadata/               # Tables, permissions, actions
├── internal/                   # Go application code (clean architecture)
│   ├── api/handler/            # HTTP handlers + Hasura Actions handler
│   ├── api/middleware/         # Auth, CORS, rate limiting, logging
│   ├── api/router/             # Route definitions
│   ├── auth/                   # JWT, password hashing, auth service
│   ├── domain/                 # Domain entities and errors
│   ├── repository/             # Data access layer
│   └── service/                # Business logic
├── migrations/                 # PostgreSQL migrations (golang-migrate)
├── pkg/                        # Reusable packages (database, logger, etc.)
├── terraform/
│   ├── modules/                # networking, database, compute, storage, cache, monitoring
│   └── environments/           # dev, staging, prod
└── test/                       # e2e, integration, testdata
```

## Local Development

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Hasura CLI (`npm install -g hasura-cli` or binary)
- Firebase CLI (optional, for function deployment)

### Quick Start

```bash
# 1. Start all services (Postgres, Go API, Hasura, Envoy, OTEL)
docker compose -f deployments/docker/docker-compose.yml up -d

# 2. Run database migrations
make migrate-up

# 3. Apply Hasura metadata (tracks tables + sets permissions)
make hasura-metadata-apply

# 4. Open Hasura Console
# Visit http://localhost:8080/console
# Admin secret: hasura-dev-admin-secret

# 5. Verify Go API health
curl http://localhost:3000/health
```

### Make Targets

```bash
make build                  # Compile Go binaries
make run                    # Build + run API server
make test                   # Run all tests
make lint                   # Run golangci-lint
make migrate-up             # Run DB migrations forward
make migrate-down           # Rollback last migration
make docker-up              # Start Docker Compose stack
make docker-down            # Stop Docker Compose stack
make hasura-console         # Open Hasura Console (via CLI)
make hasura-metadata-apply  # Apply metadata from files to Hasura
make hasura-metadata-export # Export metadata from Hasura to files
```

### Ports Summary (Local Dev)

| Port | Service |
|------|---------|
| 80 | Envoy (reverse proxy) |
| 3000 | Go API |
| 5432 | PostgreSQL |
| 8080 | Hasura GraphQL + Console |
| 4317 | OTEL Collector (gRPC) |
| 4318 | OTEL Collector (HTTP) |
| 9901 | Envoy admin |

## Database Migrations

Migrations are managed by `golang-migrate` (not Hasura). Hasura tracks existing tables via metadata.

**Migration files:** `migrations/000001_*.sql` through `migrations/000010_*.sql`

| Migration | Description |
|-----------|-------------|
| 000001 | Create users table with roles, indexes, updated_at trigger |
| 000002 | Add password_hash column to users |
| 000003 | Create enterprises table (multi-tenant root) |
| 000004 | Alter users: add enterprise_id FK, firebase_uid, migrate roles to SONA |
| 000005 | Create residents table |
| 000006 | Create robots table (CoCo inventory) |
| 000007 | Create caregiver_residents junction (max 5 rule trigger) |
| 000008 | Create stories table |
| 000009 | Create incidents table |
| 000010 | Create robot_sessions table |

**Workflow for new tables:**
1. Add a new `.up.sql` + `.down.sql` migration file
2. Run `make migrate-up`
3. Add a Hasura metadata YAML for the new table in `hasura/metadata/databases/default/tables/`
4. Update `tables.yaml` to include the new file
5. Run `make hasura-metadata-apply`

## Infrastructure (Digital Ocean)

Provisioned via Terraform in `terraform/`.

### Resources

| Module | Resource | Description |
|--------|----------|-------------|
| networking | VPC + Firewall | Private network, DB port restricted to VPC |
| database | Managed PostgreSQL 16 | Primary data store |
| compute | DOKS (Managed K8s) | Application runtime |
| storage | DO Spaces | Static assets (Svelte builds, media) |
| cache | Managed Redis | Session cache, rate limiting |
| monitoring | DO Alerts | CPU, memory, DB connection alerts |

### Deploy

```bash
cd terraform/environments/dev
export TF_VAR_do_token="your-digital-ocean-api-token"
terraform init
terraform plan
terraform apply
```

## Firebase Auth Flow

1. User signs up via Firebase Authentication (client-side)
2. Firebase Cloud Function (`firebase/functions/sync-user/`) triggers on user creation
3. Cloud Function POSTs user data (uid, email, display_name) to Go API
4. Go API inserts user into PostgreSQL with `firebase_uid` field
5. Go API returns JWT with Hasura-compatible claims
6. Client uses JWT for all subsequent Hasura GraphQL requests

## Envoy Routing

| Subdomain / Path | Destination | Description |
|-------------------|-------------|-------------|
| `enterprise.machanirobotics.dev` | Svelte App | Admin/enterprise management UI |
| `family.machanirobotics.dev` | Flutter App | Family member portal |
| `*/graphql` | Hasura :8080 | GraphQL API |
| `*/api/*` | Go API :3000 | Custom REST endpoints |

## HIPAA Compliance Notes

- **Encryption at rest**: Digital Ocean managed databases use native encryption
- **PII stripping**: OTEL Collector strips email, phone, medical_notes, SSN from telemetry
- **Encrypted columns**: `medical_notes_encrypted`, `content_encrypted`, `conversation_summary_encrypted` — encrypt at application layer before storage
- **Audit logging**: All mutations are logged via Hasura webhook events
- **Access control**: Row-level security enforced by Hasura permissions per role
- **Password hashing**: bcrypt via `golang.org/x/crypto` — `password_hash` column excluded from all GraphQL responses

## API Endpoints (Go Backend)

### Public
- `GET /health` — Health check
- `GET /ping` — Liveness probe

### Auth (via Hasura Actions → Go API)
- `POST /api/v1/actions/login` — Login (email + password)
- `POST /api/v1/actions/register` — Register new user
- `POST /api/v1/actions/refresh` — Refresh JWT tokens

### Legacy REST (available but clients should use GraphQL)
- `GET/POST /api/v1/users` — List/create users
- `GET/PUT/DELETE /api/v1/users/:id` — User CRUD

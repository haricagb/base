<!-- structure.md -->

# Project Structure

```
my-application/
├── .github/                          # GitHub specific files
│   ├── workflows/                    # CI/CD workflows
│   │   ├── ci.yml
│   │   ├── cd.yml
│   │   └── security-scan.yml
│   ├── CODEOWNERS
│   └── pull_request_template.md
│
├── cmd/                              # Main applications for this project
│   ├── api/                          # API server entry point
│   │   └── main.go
│   ├── worker/                       # Background worker entry point
│   │   └── main.go
│   └── migration/                    # DB migration tool
│       └── main.go
│
├── internal/                         # Private application code (cannot be imported)
│   ├── api/                          # API layer
│   │   ├── handler/                  # HTTP handlers
│   │   │   ├── user_handler.go
│   │   │   ├── health_handler.go
│   │   │   └── handler.go
│   │   ├── middleware/               # HTTP middlewares
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   ├── cors.go
│   │   │   └── rate_limit.go
│   │   ├── request/                  # Request DTOs
│   │   │   └── user_request.go
│   │   └── response/                 # Response DTOs
│   │       └── user_response.go
│   │
│   ├── service/                      # Business logic layer
│   │   ├── user_service.go
│   │   └── interfaces.go
│   │
│   ├── repository/                   # Data access layer
│   │   ├── user_repository.go
│   │   ├── postgres/                 # DB specific implementations
│   │   │   └── user_postgres.go
│   │   └── interfaces.go
│   │
│   ├── domain/                       # Domain models/entities
│   │   ├── user.go
│   │   └── errors.go
│   │
│   ├── worker/                       # Background job workers
│   │   ├── email_worker.go
│   │   └── notification_worker.go
│   │
│   └── validator/                    # Custom validators
│       └── user_validator.go
│
├── pkg/                              # Public libraries (can be imported by external apps)
│   ├── logger/                       # Logging package
│   │   ├── logger.go
│   │   └── config.go
│   ├── database/                     # Database utilities
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── migration.go
│   ├── cache/                        # Cache utilities
│   │   └── redis_cache.go
│   ├── httputil/                     # HTTP utilities
│   │   ├── client.go
│   │   └── response.go
│   ├── jwt/                          # JWT utilities
│   │   └── jwt.go
│   └── validator/                    # Common validators
│       └── validator.go
│
├── config/                           # Configuration files
│   ├── config.go                     # Config struct and loader
│   ├── config.yaml                   # Default config
│   ├── config.dev.yaml               # Development config
│   ├── config.staging.yaml           # Staging config
│   └── config.prod.yaml             # Production config
│
├── migrations/                       # Database migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
│
├── terraform/                        # Infrastructure as Code
│   ├── environments/
│   │   ├── dev/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   ├── terraform.tfvars
│   │   │   └── backend.tf
│   │   ├── staging/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   ├── terraform.tfvars
│   │   │   └── backend.tf
│   │   └── prod/
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       ├── outputs.tf
│   │       ├── terraform.tfvars
│   │       └── backend.tf
│   │
│   └── modules/                      # Reusable Terraform modules
│       ├── compute/
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   └── outputs.tf
│       ├── networking/
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   └── outputs.tf
│       ├── database/
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   └── outputs.tf
│       ├── storage/
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   └── outputs.tf
│       ├── cache/
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   └── outputs.tf
│       └── monitoring/
│           ├── main.tf
│           ├── variables.tf
│           └── outputs.tf
│
├── deployments/                      # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile
│   │   ├── Dockerfile.dev
│   │   └── docker-compose.yml
│   ├── kubernetes/
│   │   ├── base/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── configmap.yaml
│   │   │   ├── secret.yaml
│   │   │   └── kustomization.yaml
│   │   └── overlays/
│   │       ├── dev/
│   │       │   └── kustomization.yaml
│   │       ├── staging/
│   │       │   └── kustomization.yaml
│   │       └── prod/
│   │           └── kustomization.yaml
│   └── helm/
│       └── my-app/
│           ├── Chart.yaml
│           ├── values.yaml
│           ├── values-dev.yaml
│           ├── values-staging.yaml
│           ├── values-prod.yaml
│           └── templates/
│               ├── deployment.yaml
│               ├── service.yaml
│               ├── ingress.yaml
│               └── configmap.yaml
│
├── scripts/                          # Build and utility scripts
│   ├── build.sh
│   ├── test.sh
│   ├── deploy.sh
│   ├── migration-up.sh
│   ├── migration-down.sh
│   └── seed.sh
│
├── test/                             # Additional test files
│   ├── integration/
│   │   ├── api_test.go
│   │   └── database_test.go
│   ├── e2e/
│   │   └── user_flow_test.go
│   └── testdata/
│       ├── fixtures.json
│       └── mock_data.sql
│
├── docs/                             # Documentation
│   ├── api/
│   │   ├── swagger.yaml
│   │   └── postman_collection.json
│   ├── architecture/
│   │   └── architecture.md
│   └── deployment/
│       └── deployment.md
│
├── api/                              # API specs and definitions
│   ├── openapi/
│   │   └── openapi.yaml
│   └── proto/
│       └── service.proto
│
├── tools/                            # Supporting tools
│   └── tools.go
│
├── .env.example
├── .env.dev
├── .env.staging
├── .gitignore
├── .dockerignore
├── .editorconfig
├── .golangci.yml
├── go.mod
├── Makefile
├── README.md
├── LICENSE
└── CHANGELOG.md
```

# Makefile — Build automation for my-application

.PHONY: build run test lint migrate-up migrate-down docker-up docker-down tidy clean hasura-console hasura-metadata-apply hasura-metadata-export hasura-metadata-reload

APP_NAME := my-application
BINARY_API := bin/api
BINARY_MIGRATE := bin/migrate

## Build

build:
	go build -o $(BINARY_API) ./cmd/api
	go build -o $(BINARY_MIGRATE) ./cmd/migration

run: build
	APP_ENV=dev ./$(BINARY_API)

clean:
	rm -rf bin/

## Dependencies

tidy:
	go mod tidy

## Testing

test:
	go test ./... -v -count=1

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

## Linting

lint:
	golangci-lint run ./...

## Database

migrate-up:
	go run ./cmd/migration -direction=up

migrate-down:
	go run ./cmd/migration -direction=down

## Docker

docker-up:
	docker compose -f deployments/docker/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker/docker-compose.yml down

docker-build:
	docker compose -f deployments/docker/docker-compose.yml build

## Hasura

hasura-console:
	cd hasura && hasura console --admin-secret hasura-dev-admin-secret

hasura-metadata-apply:
	cd hasura && hasura metadata apply --admin-secret hasura-dev-admin-secret

hasura-metadata-export:
	cd hasura && hasura metadata export --admin-secret hasura-dev-admin-secret

hasura-metadata-reload:
	cd hasura && hasura metadata reload --admin-secret hasura-dev-admin-secret

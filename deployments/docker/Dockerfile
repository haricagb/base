# deployments/docker/Dockerfile — Multi-stage production build

# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/migrate ./cmd/migration

# Run stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
COPY config/ ./config/
COPY migrations/ ./migrations/

EXPOSE 3000

CMD ["./api"]

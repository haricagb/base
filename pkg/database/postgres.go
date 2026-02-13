// pkg/database/postgres.go
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresConfig holds connection parameters.
type PostgresConfig struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// NewPostgresPool creates and validates a pgxpool connection pool.
func NewPostgresPool(ctx context.Context, cfg PostgresConfig, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	logger.Info("PostgreSQL connection pool established",
		slog.String("host", poolCfg.ConnConfig.Host),
		slog.Int("max_conns", int(cfg.MaxConns)),
	)

	return pool, nil
}

// healthCheckTimeout is the maximum time a health ping may take.
const healthCheckTimeout = 5 * time.Second

// HealthCheck pings the database and returns nil if healthy.
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()
	return pool.Ping(ctx)
}

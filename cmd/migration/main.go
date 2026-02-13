// cmd/migration/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"my-application/config"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	steps := flag.Int("steps", 0, "Number of steps (0 = all)")
	flag.Parse()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	cfg, err := config.Load("config", env)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// pgx/v5 driver requires "pgx5://" scheme instead of "postgres://"
	dsn := strings.Replace(cfg.Database.DSN(), "postgres://", "pgx5://", 1)
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	default:
		log.Fatalf("unknown direction: %s (use 'up' or 'down')", *direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("Migration %s completed successfully\n", *direction)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"opensesame/internal/config"
	"opensesame/internal/httpserver"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	dsn := "postgres://opensesame_user:supersecret@localhost:5432/opensesame?sslmode=disable"

	if err := migrateDB(dsn); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := httpserver.Start(cfg); err != nil {
		log.Fatalf("error starting HTTP server: %v", err)
	}
}

func migrateDB(databaseURL string) error {
	// figure out absolute path to ./migrations
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %w", err)
	}
	sourcePath := filepath.Join(wd, "migrations")
	sourceURL := "file://" + sourcePath

	log.Printf("migrate: looking for files in %s", sourcePath)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("migrate.New failed: %w", err)
	}

	log.Println("migrate: running Up()")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Up failed: %w", err)
	}
	log.Println("migrate: Up() done (or no change)")

	version, dirty, verr := m.Version()
	if verr != nil && verr != migrate.ErrNilVersion {
		log.Printf("migrate: Version() failed: %v", verr)
	} else {
		log.Printf("migrate: current version=%d, dirty=%v", version, dirty)
	}
	return nil
}

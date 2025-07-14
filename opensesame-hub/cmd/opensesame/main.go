package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"opensesame/internal/config"
	"opensesame/internal/httpserver"
	"opensesame/internal/model"
)

func main() {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// find the sqlite app.db
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working dir: %v", err)
	}
	dbFile := filepath.Join(wd, "app.db")
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		f, err := os.Create(dbFile)
		if err != nil {
			log.Fatalf("could not create sqlite file: %v", err)
		}
		f.Close()
	}

	dsn := fmt.Sprintf("%s?_foreign_keys=1", dbFile)
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite via GORM: %v", err)
	}

	// auto migrate models on startup
	if err := gdb.AutoMigrate(
		&model.Entry{},
		&model.ControlClient{},
		&model.ControlClientEntry{},
		&model.SystemConfig{},
		&model.Session{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("AutoMigrate completed")

	// create http server
	mux := httpserver.AddHttpRoutes(gdb)
	if err := httpserver.Start(cfg, mux); err != nil {
		log.Fatalf("error starting HTTP server: %v", err)
	}
}

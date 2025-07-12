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
	"opensesame/internal/router"
)

func main() {
	// 1) load your existing app config
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// 2) ensure app.db exists in project root
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

	// 3) open with GORM + SQLite
	dsn := fmt.Sprintf("%s?_foreign_keys=1", dbFile)
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite via GORM: %v", err)
	}

	// 4) auto‐migrate all your models
	if err := gdb.AutoMigrate(
		&model.Entry{},
		&model.ControlClient{},
		&model.ControlClientEntry{},
		&model.SystemInfo{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("GORM AutoMigrate completed")

	// 5) extract *sql.DB if your HTTP server expects it, or pass *gorm.DB
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB from GORM: %v", err)
	}
	defer sqlDB.Close()

	mux := router.AddRoutes(gdb)

	// 6) start your server (update httpserver.Start to accept *gorm.DB or *sql.DB)
	if err := httpserver.Start(cfg, mux); err != nil {
		log.Fatalf("error starting HTTP server: %v", err)
	}
}

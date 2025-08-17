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
	"opensesame/internal/models/db"
	"opensesame/internal/models/types"
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

func main() {
	cfg := loadAppConfig()
	gdb := setupDatabase("os_data.db")
	repos := createRepositories(gdb)
	svcs := createServices(repos)
	mux := httpserver.AddHTTPRoutes(svcs)

	if err := httpserver.Start(cfg, mux); err != nil {
		log.Fatalf("error starting HTTP server: %v", err)
	}
}

func loadAppConfig() *config.Config {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return cfg
}

func setupDatabase(filename string) *gorm.DB {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working dir: %v", err)
	}
	dbFile := filepath.Join(wd, filename)
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

	if err := gdb.AutoMigrate(
		&db.SystemConfig{},
		&db.EntryDevice{},
		&db.EntryCommand{},
		&db.HttpCommand{},
		// &db.UdpCommand{}, // not supported yet
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return gdb
}

func createRepositories(gdb *gorm.DB) *types.Repositories {
	return &types.Repositories{
		Config: repository.NewConfigRepository(gdb),
		Entry:  repository.NewEntryRepository(gdb),
	}
}

func createServices(repos *types.Repositories) *types.Services {
	configSvc := service.NewConfigService(repos.Config)

	return &types.Services{
		Auth:   service.NewAuthService(configSvc),
		Config: configSvc,
		Entry:  service.NewEntryService(repos.Entry),
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	// Start discovery service in background
	if err := svcs.Discovery.Start(context.Background()); err != nil {
		log.Fatalf("failed to start discovery service: %v", err)
	}
	defer svcs.Discovery.Stop()

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

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // log slow queries
			LogLevel:                  logger.Warn, // log only warnings or errors
			IgnoreRecordNotFoundError: true,        // 👈 don't log ErrRecordNotFound
			Colorful:                  true,
		},
	)

	// Open DB with custom logger
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to open sqlite via GORM: %v", err)
	}

	// Run migrations
	if err := gdb.AutoMigrate(
		&db.SystemConfig{},
		&db.EntryDevice{},
		&db.EntryCommand{},
		&db.DiscoveredDevice{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return gdb
}

func createRepositories(gdb *gorm.DB) *types.Repositories {
	return &types.Repositories{
		Config:           repository.NewConfigRepository(gdb),
		Entry:            repository.NewEntryRepository(gdb),
		DiscoveredDevice: repository.NewDiscoveredDeviceRepository(gdb),
	}
}

func createServices(repos *types.Repositories) *types.Services {
	configSvc := service.NewConfigService(repos.Config)

	return &types.Services{
		Auth:      service.NewAuthService(configSvc),
		Config:    configSvc,
		Entry:     service.NewEntryService(repos.Entry),
		Discovery: service.NewDiscoveryService(repos.DiscoveredDevice),
		Status:    service.NewStatusService(repos.Config, repos.Entry, repos.DiscoveredDevice),
	}
}

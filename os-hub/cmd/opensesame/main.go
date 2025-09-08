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
	"opensesame/internal/etag"
	"opensesame/internal/httpserver"
	"opensesame/internal/models/db"
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

func main() {
	cfg := loadAppConfig()
	gdb := setupDatabase("os_data.db")
	repos := createRepositories(gdb)
	svcs := createServices(repos)

	etag.Init()

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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
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
		&db.Device{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return gdb
}

func createRepositories(gdb *gorm.DB) *repository.RepositoriesType {
	return &repository.RepositoriesType{
		Config: repository.NewConfigRepository(gdb),
		Device: repository.NewDeviceRepository(gdb),
	}
}

func createServices(repos *repository.RepositoriesType) *service.ServicesType {
	configSvc := service.NewConfigService(repos.Config)
	entrySvc := service.NewEntryService(repos.Device)

	return &service.ServicesType{
		Auth:      service.NewAuthService(configSvc),
		Config:    configSvc,
		Entry:     entrySvc,
		Discovery: service.NewDiscoveryService(repos.Device, entrySvc),
		Status:    service.NewStatusService(repos.Config, repos.Device),
	}
}

package main

import (
	"context"
	"log"
	"opensesame/internal/config"
	"opensesame/internal/httpserver"
)

func main() {
	// Load configuration (with context)
	cfg, err := config.LoadConfig(context.Background()) // Pass context and an empty path
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Delegate server startup (applies logging middleware)
	if err := httpserver.Start(cfg); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

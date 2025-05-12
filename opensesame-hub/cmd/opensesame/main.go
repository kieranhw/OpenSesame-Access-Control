package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"opensesame/internal/config"
	"opensesame/internal/router"
)

func main() {
	// Load configuration (with context)
	cfg, err := config.LoadConfig(context.Background()) // Pass context and an empty path
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize router and routes
	r := router.AddRoutes()

	// Start the HTTP server
	address := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	fmt.Printf("Starting server on %s\n", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

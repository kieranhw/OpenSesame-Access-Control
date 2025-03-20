// cmd/opensesame/main.go
package main

import (
	"fmt"
	"log"

	"opensesame/internal/config"
	"opensesame/internal/httpserver"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("Starting Security Hub...")

	// Start the HTTP server (this will use your HTTP server package)
	if err := httpserver.Start(cfg); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

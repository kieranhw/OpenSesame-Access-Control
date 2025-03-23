package main

import (
	"context"
	"fmt"
	"log"

	"opensesame/internal/config"
	"opensesame/internal/httpserver"
)

func main() {
	// Load configuration from environment variables (no config file)
	cfg, err := config.LoadConfig(context.Background(), "")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Print out configuration details
	fmt.Println("Starting Security Hub...")
	fmt.Printf("HTTP Listener Port: %s\n", cfg.HttpListenerPort)
	fmt.Printf("Management Port: %s\n", cfg.ManagementPort)
	fmt.Printf("TCP Listener Port: %s\n", cfg.TcpListenerPort)

	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		fmt.Println("TLS enabled")
	} else {
		fmt.Println("TLS disabled")
	}

	// Start the HTTP server (this will use your HTTP server package)
	if err := httpserver.Start(cfg); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

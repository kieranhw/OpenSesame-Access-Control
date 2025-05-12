package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/router"
)

func Start(cfg *config.Config) error {
	r := router.AddRoutes()

	address := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	fmt.Printf("Starting server on %s\n", address)

	// Start the HTTP server
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
		return err
	}
	return nil
}

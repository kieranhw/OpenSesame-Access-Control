package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/middleware" // Import the middleware package
	"opensesame/internal/router"
)

func Start(cfg *config.Config) error {
	mux := router.AddRoutes()
	loggedMux := middleware.Logger(mux)

	address := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	fmt.Printf("Starting HTTP server on %s\n", address)

	if err := http.ListenAndServe(address, loggedMux); err != nil {
		log.Fatalf("Error starting server: %v", err)
		return err
	}
	return nil
}

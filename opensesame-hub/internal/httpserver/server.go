package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/httpserver/handlers"
)

// Start initializes and starts the HTTP server.
func Start(cfg *config.Config) error {
	mux := http.NewServeMux()

	// Register endpoints for different API categories.
	handlers.RegisterDoorRoutes(mux)
	handlers.RegisterManagementRoutes(mux)

	// Optional catch-all route for unmatched endpoints.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Start HTTPS if TLS config is provided.
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		log.Printf("Starting HTTPS server on :%s", cfg.HttpListenerPort)
		return http.ListenAndServeTLS(fmt.Sprintf(":%s", cfg.HttpListenerPort), cfg.TLSCert, cfg.TLSKey, mux)
	}

	log.Printf("Starting HTTP server on :%s", cfg.HttpListenerPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpListenerPort), mux)
}

package httpserver

import (
	"log"
	"net/http"

	"opensesame/internal/config"
)

// Start starts the HTTP server.
func Start(cfg *config.Config) error {
	// Set up your routes.
	mux := http.NewServeMux()
	// Example routes:
	mux.HandleFunc("/doors/", DoorHandler)

	// For simplicity, we'll run without TLS if TLS config isn't provided.
	if cfg.TLSCert != "" && cfg.TLSKey != "" {
		log.Printf("Starting HTTPS server on %s", cfg.Port)
		return http.ListenAndServeTLS(cfg.Port, cfg.TLSCert, cfg.TLSKey, mux)
	}
	log.Printf("Starting HTTP server on %s", cfg.Port)
	return http.ListenAndServe(cfg.Port, mux)
}

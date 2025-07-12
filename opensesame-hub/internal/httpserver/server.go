package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/middleware"
)

func Start(cfg *config.Config, handler http.Handler) error {
	// wrap with your logger middleware
	logged := middleware.Logger(handler)

	addr := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	log.Printf("starting HTTP server on %s", addr)

	return http.ListenAndServe(addr, logged)
}

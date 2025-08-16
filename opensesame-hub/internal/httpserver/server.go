package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/handlers/management"
	"opensesame/internal/middleware"
	"opensesame/internal/models/types"

	"github.com/gorilla/mux"
)

func Start(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	log.Printf("starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func AddHttpRoutes(svcs *types.Services) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.HttpLogger)
	r.Use(middleware.ValidateJSONBody)

	management.MountRoutes(r, svcs)

	return r
}

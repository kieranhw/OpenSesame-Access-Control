package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	handlers "opensesame/internal/handlers/management"
	"opensesame/internal/middleware"
	"opensesame/internal/models/types"

	"github.com/gorilla/mux"
)

func Start(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func AddHTTPRoutes(svcs *types.Services) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.HTTPLogger)
	r.Use(middleware.ValidateJSONBody)

	handlers.MountManagementRoutes(r, svcs)

	return r
}

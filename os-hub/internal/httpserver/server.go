package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	handlers "opensesame/internal/handlers"
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func Start(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	log.Printf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func AddHTTPRoutes(svcs *service.ServicesType) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.HTTPLogger)
	r.Use(middleware.ValidateJSONBody)

	handlers.MountRoutes(r, svcs)

	return r
}

package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"opensesame/internal/config"
	"opensesame/internal/handlers/management"
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func Start(cfg *config.Config, handler http.Handler) error {
	addr := fmt.Sprintf(":%s", cfg.HttpListenerPort)
	log.Printf("starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func AddHttpRoutes(db *gorm.DB) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.HttpLogger)
	r.Use(middleware.ValidateJSONBody)

	configSvc := service.NewConfigService(db)
	authSvc := service.NewAuthService(db)

	management.MountRoutes(r, configSvc, authSvc)

	return r
}

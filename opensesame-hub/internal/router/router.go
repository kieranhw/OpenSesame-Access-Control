package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"opensesame/internal/handlers/management"
	"opensesame/internal/service"
)

func AddRoutes(db *gorm.DB) http.Handler {
	r := mux.NewRouter()

	// instantiate the SetupService
	setupSvc := service.NewSetupService(db)

	// register management endpoints
	management.MountRoutes(r, setupSvc)

	return r
}

package management

import (
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MountRoutes(parent *mux.Router, setupSvc *service.SetupService) {
	mgmt := parent.PathPrefix("/management").Subrouter()
	mgmt.Use(middleware.MgmtAuthValidator())

	mgmt.HandleFunc("/setup", SystemSetupHandler(setupSvc)).Methods("GET", "POST")
	mgmt.HandleFunc("/access", GetAccessHandler).Methods("GET")
	mgmt.HandleFunc("/access", PostAccessHandler).Methods("POST")
}

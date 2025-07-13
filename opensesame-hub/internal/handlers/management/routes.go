package management

import (
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MountRoutes(parent *mux.Router, setupSvc *service.ConfigService) {
	mgmt := parent.PathPrefix("/management").Subrouter()
	mgmt.Use(middleware.MgmtAuthValidator())

	mgmt.HandleFunc("/config", GetSystemConfig(setupSvc)).Methods("GET")
	mgmt.HandleFunc("/config", PostSystemConfig(setupSvc)).Methods("POST")
	mgmt.HandleFunc("/config", PatchSystemConfig(setupSvc)).Methods("PATCH")

	mgmt.HandleFunc("/access", GetAccessHandler).Methods("GET")
	mgmt.HandleFunc("/access", PostAccessHandler).Methods("POST")
}

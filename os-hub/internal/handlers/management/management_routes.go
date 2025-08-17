package handlers

import (
	"net/http"
	"opensesame/internal/middleware"
	"opensesame/internal/models/types"

	"github.com/gorilla/mux"
)

func MountManagementRoutes(parent *mux.Router, svcs *types.Services) {
	// TODO: this should be changed in production to 127.0.0.1:443, where the management app is hosted
	parent.Use(middleware.CORSMiddleware("http://localhost:3000"))

	parent.PathPrefix("/management/").Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

	// add public routes to parent
	parent.HandleFunc("/management/config", GetSystemConfig(svcs.Config)).Methods("GET")
	parent.HandleFunc("/management/config", CreateSystemConfig(svcs.Config)).Methods("POST")

	parent.HandleFunc("/management/session", LoginHandler(svcs.Auth)).Methods("POST")
	parent.HandleFunc("/management/session", ValidateSessionHandler(svcs.Config, svcs.Auth)).Methods("GET")
	parent.HandleFunc("/management/session", LogoutHandler()).Methods("DELETE")

	// add protected routes to subrouter
	mgmt := parent.PathPrefix("/management").Subrouter()
	mgmt.Use(middleware.MgmtSessionValidator(svcs.Config, svcs.Auth))
	mgmt.HandleFunc("/config", UpdateSystemConfig(svcs.Config)).Methods("PATCH")

	// entry
	mgmt.HandleFunc("/entry_devices", ListEntryDevices(svcs.Entry)).Methods("GET")
}

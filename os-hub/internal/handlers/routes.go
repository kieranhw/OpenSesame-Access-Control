package handlers

import (
	"net/http"
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MountRoutes(parent *mux.Router, svcs *service.ServicesType) {
	// TODO: this should be changed in production to 127.0.0.1:443, where the management app is hosted
	parent.Use(middleware.CORSMiddleware("http://localhost:3000"))
	parent.PathPrefix("/").Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

	// public routes
	parent.HandleFunc("/config", GetSystemConfig(svcs.Config)).Methods("GET")
	parent.HandleFunc("/config", CreateSystemConfig(svcs.Config)).Methods("POST")

	// admin auth routes, but not protected by session
	parent.HandleFunc("/admin/session", LoginHandler(svcs.Auth)).Methods("POST")
	parent.HandleFunc("/admin/session", ValidateSessionHandler(svcs.Config, svcs.Auth)).Methods("GET")
	parent.HandleFunc("/admin/session", LogoutHandler()).Methods("DELETE")

	// protected routes
	admin := parent.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.MgmtSessionValidator(svcs.Config, svcs.Auth))

	// system
	admin.HandleFunc("/config", UpdateSystemConfig(svcs.Config)).Methods("PATCH")
	admin.HandleFunc("/status", GetStatus(svcs)).Methods("GET")

	// entry
	admin.HandleFunc("/entry_devices", ListEntryDevices(svcs.Entry)).Methods("GET")
	admin.HandleFunc("/entry_devices", CreateEntryDevice(svcs.Entry)).Methods("POST")

}

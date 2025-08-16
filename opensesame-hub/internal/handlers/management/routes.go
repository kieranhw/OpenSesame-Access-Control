package management

import (
	"net/http"
	"opensesame/internal/middleware"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MountRoutes(
	parent *mux.Router,
	configSvc *service.ConfigService,
	authSvc *service.AuthService,
) {
	parent.Use(middleware.CORSMiddleware("http://localhost:3000"))

	parent.PathPrefix("/management/").Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

	// public routes
	parent.HandleFunc("/management/config", GetSystemConfig(configSvc)).Methods("GET")
	parent.HandleFunc("/management/config", CreateSystemConfig(configSvc)).Methods("POST")

	parent.HandleFunc("/management/session", LoginHandler(authSvc)).Methods("POST")
	parent.HandleFunc("/management/session", ValidateSessionHandler(configSvc, authSvc)).Methods("GET")
	parent.HandleFunc("/management/session", LogoutHandler()).Methods("DELETE")

	// protected routes
	mgmt := parent.PathPrefix("/management").Subrouter()
	mgmt.Use(middleware.MgmtSessionValidator(configSvc, authSvc))

	mgmt.HandleFunc("/config", UpdateSystemConfig(configSvc)).Methods("PATCH")
}

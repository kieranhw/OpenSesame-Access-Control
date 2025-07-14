// management/router.go
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

	parent.PathPrefix("/management/").Methods("OPTIONS").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

	// public routes
	parent.HandleFunc("/management/config", GetSystemConfig(configSvc)).Methods("GET")
	parent.HandleFunc("/management/config", PostSystemConfig(configSvc)).Methods("POST")
	parent.HandleFunc("/management/login", LoginHandler(configSvc, authSvc)).Methods("POST")
	parent.HandleFunc("/management/validate_session", SessionHandler(configSvc)).Methods("GET")

	// protected routes
	mgmt := parent.PathPrefix("/management").Subrouter()
	mgmt.Use(middleware.MgmtSessionValidator(authSvc))
	mgmt.HandleFunc("/config", PatchSystemConfig(configSvc)).Methods("PATCH")
	mgmt.HandleFunc("/access", GetAccessHandler).Methods("GET")
	mgmt.HandleFunc("/access", PostAccessHandler).Methods("POST")
}

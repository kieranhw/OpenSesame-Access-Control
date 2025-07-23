package middleware

import (
	"net/http"

	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MgmtSessionValidator(configSvc *service.ConfigService, authSvc *service.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg, err := configSvc.GetSystemConfig(r.Context())
			if err != nil || cfg == nil {
				http.Error(w, "system configuration required", http.StatusPreconditionRequired)
				return
			}
			systemSecret := cfg.SystemSecret

			c, err := r.Cookie("os_session")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			isValid, err := authSvc.ValidateSession(r.Context(), c.Value, systemSecret)

			if !isValid || err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

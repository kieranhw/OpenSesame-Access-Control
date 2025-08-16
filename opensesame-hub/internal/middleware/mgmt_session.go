package middleware

import (
	"net/http"

	"opensesame/internal/constants"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MgmtSessionValidator(configSvc *service.ConfigService, authSvc *service.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(constants.SessionCookieName)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			isValid, err := authSvc.ValidateSession(r.Context(), cookie.Value)
			if !isValid || err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

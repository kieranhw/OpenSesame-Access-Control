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

			// Refresh the session if valid
			newCookie, err := authSvc.RefreshSession(r.Context(), cookie)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Set the refreshed cookie in the response
			http.SetCookie(w, newCookie)
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"
	"opensesame/internal/service"

	"github.com/gorilla/mux"
)

func MgmtSessionValidator(authSvc *service.AuthService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("os_session")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ok, err := authSvc.ValidateSession(r.Context(), c.Value)
			if err != nil || !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

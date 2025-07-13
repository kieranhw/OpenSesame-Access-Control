package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MgmtAuthValidator() mux.MiddlewareFunc {
	var tokenHeader = "Management-Key"
	var validKey = "test"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: replace with real validation logic, we need to check the presented key
			// is valid and not expired, etc. Should probably use session cookies here.

			if r.Header.Get(tokenHeader) != validKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

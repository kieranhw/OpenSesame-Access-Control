package middleware

import (
	"log"
	"net/http"
	"strings"
)

// StatusRecorder wraps http.ResponseWriter to record the status code.
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *StatusRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// HttpLogger logs each request and its response status code.
func HttpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(rec, r)
		log.Printf("Response: %d", rec.StatusCode)
	})
}

func ValidateJSONBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			// Check Content-Type
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				http.Error(
					w,
					"Content-Type must be application/json",
					http.StatusUnsupportedMediaType,
				)
				return
			}

			// Check non-empty body
			if r.ContentLength == 0 {
				http.Error(
					w,
					"request body is required",
					http.StatusBadRequest,
				)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func JSONOnly(next http.Handler) http.Handler {
	return ValidateJSONBody(next)
}

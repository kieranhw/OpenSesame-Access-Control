package middleware

import (
	"log"
	"net/http"
)

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *StatusRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// A middleware to log HTTP requests and responses.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &StatusRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(rec, r)
		log.Printf("Response: %d", rec.StatusCode)
	})
}

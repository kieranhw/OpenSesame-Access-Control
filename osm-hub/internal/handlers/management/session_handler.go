package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"opensesame/internal/models/dto"
	"opensesame/internal/service"
)

func LoginHandler(authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req dto.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := "invalid JSON payload"
			json.NewEncoder(w).Encode(dto.SessionResponse{
				Message: &msg,
			})
			return
		}

		resp, cookie, err := authSvc.Login(r.Context(), req)
		if cookie != nil {
			http.SetCookie(w, cookie)
		}

		switch {
		case errors.Is(err, service.ErrNotConfigured):
			w.WriteHeader(http.StatusPreconditionRequired)
		case err != nil:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func ValidateSessionHandler(configSvc *service.ConfigService, authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		configured, err := configSvc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !configured {
			w.WriteHeader(http.StatusPreconditionRequired)
			json.NewEncoder(w).Encode(dto.ConfigResponse{Configured: false})
			return
		}

		cookie, err := r.Cookie("os_session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "login required"
			json.NewEncoder(w).Encode(dto.SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		isValid, err := authSvc.ValidateSession(r.Context(), cookie.Value)
		if err == service.ErrNotConfigured {
			msg := service.ErrNotConfigured.Error()
			json.NewEncoder(w).Encode(dto.SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    false,
			})
			return
		} else if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "session validation error"
			json.NewEncoder(w).Encode(dto.SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		} else if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "invalid or expired session"
			json.NewEncoder(w).Encode(dto.SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		json.NewEncoder(w).Encode(dto.SessionResponse{
			Authenticated: true,
			Configured:    true,
		})
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		http.SetCookie(w, &http.Cookie{
			Name:     "os_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})

		json.NewEncoder(w).Encode(dto.LogoutResponse{Success: true})
	}
}

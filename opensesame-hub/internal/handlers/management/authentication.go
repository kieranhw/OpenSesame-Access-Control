package management

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"opensesame/internal/service"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type SessionResponse struct {
	Message       *string `json:"message,omitempty"`
	Authenticated bool    `json:"authenticated"`
	Configured    bool    `json:"configured"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}

func LoginHandler(configSvc *service.ConfigService, authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			msg := "method not allowed"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			msg := "invalid JSON payload"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil || cfg == nil {
			w.WriteHeader(http.StatusPreconditionRequired)
			msg := "system configuration required"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    false,
			})
			return
		}

		// validate password
		if err := bcrypt.CompareHashAndPassword(
			[]byte(cfg.AdminPasswordHash),
			[]byte(req.Password),
		); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "invalid credentials"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		// create jwt signed with the system secret
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
			"sub": "admin",
		})

		signed, err := token.SignedString([]byte(cfg.SystemSecret))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := "could not sign token"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "os_session",
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		json.NewEncoder(w).Encode(SessionResponse{
			Authenticated: true,
			Configured:    true,
		})
	}
}

func ValidateSessionHandler(configSvc *service.ConfigService, authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		c, err := r.Cookie("os_session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "login required"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil || cfg == nil {
			w.WriteHeader(http.StatusPreconditionRequired)
			msg := "system configuration required"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    false,
			})
			return
		}
		systemSecret := cfg.SystemSecret

		isValid, err := authSvc.ValidateSession(r.Context(), c.Value, systemSecret)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "session validation error"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			msg := "invalid or expired session"
			json.NewEncoder(w).Encode(SessionResponse{
				Message:       &msg,
				Authenticated: false,
				Configured:    true,
			})
			return
		}

		json.NewEncoder(w).Encode(SessionResponse{
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

		json.NewEncoder(w).Encode(LogoutResponse{Success: true})
	}
}

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

type LoginResponse struct {
	Success bool `json:"success"`
}

type SessionResponse struct {
	Authenticated bool `json:"authenticated"`
}

// LoginHandler issues an os_session JWT cookie on successful login.
func LoginHandler(
	configSvc *service.ConfigService,
	authSvc *service.AuthService, // not used here but kept for parity
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Ensure system is configured
		configured, err := configSvc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if !configured {
			http.Error(w, "system not configured", http.StatusPreconditionFailed)
			return
		}

		// Load config to get the password hash & secret
		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword(
			[]byte(cfg.AdminPasswordHash),
			[]byte(req.Password),
		); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
			"sub": "admin",
		})
		signed, err := token.SignedString([]byte(cfg.SystemSecret))
		if err != nil {
			http.Error(w, "could not sign token", http.StatusInternalServerError)
			return
		}

		// Set the cookie (host-only, HttpOnly)
		http.SetCookie(w, &http.Cookie{
			Name:     "os_session",
			Value:    signed,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		json.NewEncoder(w).Encode(LoginResponse{Success: true})
	}
}

// SessionHandler lets the client verify they have a valid os_session.
func SessionHandler(configSvc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 1) grab the cookie
		c, err := r.Cookie("os_session")
		if err != nil {
			http.Error(w, "no session cookie", http.StatusUnauthorized)
			return
		}

		// 2) load system secret
		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if cfg == nil {
			http.Error(w, "system not configured", http.StatusPreconditionFailed)
			return
		}

		// 3) parse & validate JWT
		token, err := jwt.ParseWithClaims(
			c.Value,
			&jwt.RegisteredClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.SystemSecret), nil
			},
		)
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}

		// 4) success
		json.NewEncoder(w).Encode(SessionResponse{Authenticated: true})
	}
}

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

type LogoutResponse struct {
	Success bool `json:"success"`
}

func LoginHandler(configSvc *service.ConfigService, authSvc *service.AuthService) http.HandlerFunc {
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

		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		// validate password
		if err := bcrypt.CompareHashAndPassword(
			[]byte(cfg.AdminPasswordHash),
			[]byte(req.Password),
		); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// create jwt signed with system secret
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

func SessionHandler(configSvc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		c, err := r.Cookie("os_session")
		if err != nil {
			http.Error(w, "no session cookie", http.StatusUnauthorized)
			return
		}

		// get system secret
		cfg, err := configSvc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if cfg == nil {
			http.Error(w, "system not configured", http.StatusPreconditionFailed)
			return
		}

		// validate jwt
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

		json.NewEncoder(w).Encode(SessionResponse{Authenticated: true})
	}
}

func LogoutHandler(authSvc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("os_session"); err == nil {
			_ = authSvc.DeleteSession(r.Context(), c.Value)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "os_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})

		json.NewEncoder(w).Encode(LogoutResponse{
			Success: true,
		})
	}
}

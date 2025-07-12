package management

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"opensesame/internal/model"
	"opensesame/internal/service"
)

// DTOs
type ConfigRequest struct {
	SystemName    string `json:"system_name"`
	AdminPassword string `json:"admin_password"`
	BackupCode    string `json:"backup_code"`
}

type ConfigStatus struct {
	Configured bool `json:"configured"`
}

type ConfigResponse struct {
	ID         int    `json:"id"`
	SystemName string `json:"system_name"`
}

func GetSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		isConfigured, err := svc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := ConfigStatus{Configured: isConfigured}
		json.NewEncoder(w).Encode(resp)
	}
}

func PostSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		isConfigured, err := svc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isConfigured {
			http.Error(w, "System already configured", http.StatusConflict)
			return
		}

		// decode request
		var req ConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// hash the password (TODO: swap for real bcrypt cost)
		pwdHash, err := bcrypt.GenerateFromPassword(
			[]byte(req.AdminPassword), bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// generate a system secret
		secret := uuid.New().String()

		// map to entity
		entity := &model.SystemConfig{
			SystemName:        req.SystemName,
			AdminPasswordHash: string(pwdHash),
			BackupCodeHash:    req.BackupCode,
			SystemSecret:      secret,
		}

		// persist
		if err := svc.Create(r.Context(), entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// respond
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ConfigResponse{
			ID:         entity.ID,
			SystemName: entity.SystemName,
		})
	}
}

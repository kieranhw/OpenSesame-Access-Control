package management

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"opensesame/internal/model"
	"opensesame/internal/service"
)

type CreateConfigRequest struct {
	SystemName    string `json:"system_name"`
	AdminPassword string `json:"admin_password"`
}

type UpdateConfigRequest struct {
	SystemName    *string `json:"system_name,omitempty"`
	AdminPassword *string `json:"admin_password,omitempty"`
}

type ConfigResponse struct {
	Configured bool    `json:"configured"`
	SystemName *string `json:"system_name,omitempty"`
	BackupCode *string `json:"backup_code,omitempty"`
}

func GetSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		configured, err := svc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !configured {
			json.NewEncoder(w).Encode(ConfigResponse{
				Configured: false,
			})
			return
		}

		cfg, err := svc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ConfigResponse{
			Configured: true,
			SystemName: &cfg.SystemName,
		})
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

		var req CreateConfigRequest
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
		backupCode := uuid.New().String()

		entity := &model.SystemConfig{
			SystemName:        req.SystemName,
			AdminPasswordHash: string(pwdHash),
			SystemSecret:      secret,
			BackupCode:        backupCode,
		}

		if err := svc.CreateConfig(r.Context(), entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ConfigResponse{
			Configured: true,
			SystemName: &entity.SystemName,
			BackupCode: &entity.BackupCode,
		})
	}
}

func PatchSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ok, err := svc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "system not configured", http.StatusPreconditionFailed)
			return
		}

		var req UpdateConfigRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.SystemName == nil && req.AdminPassword == nil {
			http.Error(w, "nothing to update", http.StatusBadRequest)
			return
		}

		cfg, err := svc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if req.SystemName != nil {
			cfg.SystemName = *req.SystemName
		}

		if req.AdminPassword != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(*req.AdminPassword), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "error hashing password", http.StatusInternalServerError)
				return
			}
			cfg.AdminPasswordHash = string(hash)
		}

		if err := svc.UpdateConfig(r.Context(), cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ConfigResponse{
			Configured: true,
			SystemName: &cfg.SystemName,
		})
	}
}

package management

import (
	"encoding/json"
	"errors"
	"net/http"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/service"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		configured, err := svc.IsSystemConfigured(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !configured {
			json.NewEncoder(w).Encode(dto.ConfigResponse{Configured: false})
			return
		}

		cfg, err := svc.GetSystemConfig(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(dto.ConfigResponse{
			Configured:        true,
			SystemName:        cfg.SystemName,
			SessionTimeoutSec: cfg.SessionTimeoutSec,
		})
	}
}

func CreateSystemConfig(svc *service.ConfigService) http.HandlerFunc {
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

		var req dto.CreateConfigPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pwdHash, err := bcrypt.GenerateFromPassword(
			[]byte(req.AdminPassword), bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(w, "internal error during password hashing", http.StatusInternalServerError)
			return
		}

		secret := uuid.NewString()
		backupCode := uuid.NewString()

		entity := &db.SystemConfig{
			SystemName:        req.SystemName,
			SessionTimeoutSec: 86400,
			AdminPasswordHash: string(pwdHash),
			SystemSecret:      secret,
			BackupCode:        backupCode,
		}

		if err := svc.CreateConfig(r.Context(), entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.ConfigResponse{
			Configured: true,
			SystemName: &entity.SystemName,
			BackupCode: &entity.BackupCode,
		})
	}
}

func UpdateSystemConfig(svc *service.ConfigService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var reqPayload dto.UpdateConfigPayload
		if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedCfg, err := svc.UpdateConfig(r.Context(), &reqPayload)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrNotConfigured):
				http.Error(w, "system not configured", http.StatusPreconditionFailed)
			case errors.Is(err, service.ErrNoUpdateFields):
				http.Error(w, "nothing to update", http.StatusBadRequest)
			case errors.Is(err, service.ErrPasswordHashingFailed):
				http.Error(w, "error hashing password", http.StatusInternalServerError)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(dto.ConfigResponse{
			Configured:        true,
			SystemName:        &updatedCfg.SystemName,
			SessionTimeoutSec: &updatedCfg.SessionTimeoutSec,
		})
	}
}

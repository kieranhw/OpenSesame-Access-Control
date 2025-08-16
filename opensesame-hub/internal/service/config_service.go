// service/config_service.go
package service

import (
	"context"
	"fmt"

	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type ConfigService struct {
	repo repository.ConfigRepository
}

func NewConfigService(repo repository.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

func (s *ConfigService) IsSystemConfigured(ctx context.Context) (bool, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return false, fmt.Errorf("checking system configuration: %w", err)
	}
	return count > 0, nil
}

func (s *ConfigService) GetSystemConfig(ctx context.Context) (*dto.ConfigResponse, error) {
	cfg, err := s.repo.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting system config: %w", err)
	}
	if cfg == nil {
		return nil, nil
	}

	return s.toConfigResponse(cfg), nil
}

// GetSystemConfigEntity returns the full entity including sensitive fields
// This is used internally by other services like AuthService
func (s *ConfigService) GetSystemConfigEntity(ctx context.Context) (*db.SystemConfig, error) {
	cfg, err := s.repo.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting system config entity: %w", err)
	}
	return cfg, nil
}

func (s *ConfigService) CreateConfig(ctx context.Context, cfg *db.SystemConfig) error {
	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return fmt.Errorf("error checking configuration status before create: %w", err)
	}
	if configured {
		return ErrAlreadyConfigured
	}

	if err := s.repo.CreateSystemConfig(ctx, cfg); err != nil {
		return fmt.Errorf("creating system config: %w", err)
	}
	return nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, payload *dto.UpdateConfigPayload) (*db.SystemConfig, error) {
	if payload.SystemName == nil && payload.AdminPassword == nil && payload.SessionTimeoutSec == nil {
		return nil, ErrNoUpdateFields
	}

	systemConfig, err := s.repo.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching current config for update: %w", err)
	}
	if systemConfig == nil {
		return nil, ErrNotConfigured
	}

	// Apply updates
	if err := s.applyConfigUpdates(systemConfig, payload); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateSystemConfig(ctx, systemConfig); err != nil {
		return nil, fmt.Errorf("error saving updated config: %w", err)
	}

	return systemConfig, nil
}

// Helper methods
func (s *ConfigService) toConfigResponse(cfg *db.SystemConfig) *dto.ConfigResponse {
	return &dto.ConfigResponse{
		Configured:        true,
		SystemName:        &cfg.SystemName,
		SessionTimeoutSec: &cfg.SessionTimeoutSec,
	}
}

func (s *ConfigService) applyConfigUpdates(config *db.SystemConfig, payload *dto.UpdateConfigPayload) error {
	if payload.SystemName != nil {
		config.SystemName = *payload.SystemName
	}

	if payload.AdminPassword != nil {
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(*payload.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
		}
		config.AdminPasswordHash = string(newPasswordHash)
	}

	if payload.SessionTimeoutSec != nil {
		if *payload.SessionTimeoutSec <= 0 {
			return fmt.Errorf("invalid session timeout: must be positive")
		}
		config.SessionTimeoutSec = *payload.SessionTimeoutSec
	}

	return nil
}

package service

import (
	"context"
	"fmt"
	"strings"

	"opensesame/internal/etag"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/models/types"
	"opensesame/internal/repository"

	"github.com/google/uuid"
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

func (s *ConfigService) GetSystemConfigEntity(ctx context.Context) (*db.SystemConfig, error) {
	cfg, err := s.repo.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting system config entity: %w", err)
	}
	return cfg, nil
}

func (s *ConfigService) CreateConfig(ctx context.Context, req dto.CreateConfigRequest) (*dto.ConfigResponse, error) {
	if len(strings.TrimSpace(req.SystemName)) <= 1 ||
		len(strings.TrimSpace(req.AdminPassword)) <= 1 {
		return nil, fmt.Errorf("invalid config: system name and admin password must be longer than 1 character")
	}

	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve system configuration: %w", err)
	}
	if configured {
		return nil, types.ErrAlreadyConfigured
	}

	adminPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing admin password: %w", err)
	}

	backupCode := uuid.NewString()
	backupCodeHash, err := bcrypt.GenerateFromPassword([]byte(backupCode), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing backup code: %w", err)
	}

	sysCfg := &db.SystemConfig{
		SystemName:        req.SystemName,
		SessionTimeoutSec: req.SessionTimeoutSec,
		AdminPasswordHash: string(adminPasswordHash),
		BackupCodeHash:    string(backupCodeHash),
		SystemSecret:      uuid.NewString(),
	}

	if err := s.repo.CreateSystemConfig(ctx, sysCfg); err != nil {
		return nil, fmt.Errorf("creating system config: %w", err)
	}

	etag.Bump()
	return &dto.ConfigResponse{
		Configured: true,
		SystemName: &sysCfg.SystemName,
		// include backup code once only in create response
		BackupCode:        &backupCode,
		SessionTimeoutSec: &sysCfg.SessionTimeoutSec,
	}, nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, payload *dto.UpdateConfigRequest) (*dto.ConfigResponse, error) {
	if payload.SystemName == nil && payload.AdminPassword == nil && payload.SessionTimeoutSec == nil {
		return nil, types.ErrNoUpdateFields
	}

	sysCfg, err := s.repo.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching current config for update: %w", err)
	}
	if sysCfg == nil {
		return nil, types.ErrNotConfigured
	}

	// Apply updates
	if err := s.applyConfigUpdates(sysCfg, payload); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateSystemConfig(ctx, sysCfg); err != nil {
		return nil, fmt.Errorf("error saving updated config: %w", err)
	}

	etag.Bump()
	return s.toConfigResponse(sysCfg), nil
}

func (s *ConfigService) toConfigResponse(cfg *db.SystemConfig) *dto.ConfigResponse {
	return &dto.ConfigResponse{
		Configured:        true,
		SystemName:        &cfg.SystemName,
		SessionTimeoutSec: &cfg.SessionTimeoutSec,
	}
}

func (s *ConfigService) applyConfigUpdates(config *db.SystemConfig, payload *dto.UpdateConfigRequest) error {
	if payload.SystemName != nil {
		config.SystemName = *payload.SystemName
	}

	if payload.AdminPassword != nil {
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(*payload.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("%w: %v", types.ErrPasswordHashingFailed, err)
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

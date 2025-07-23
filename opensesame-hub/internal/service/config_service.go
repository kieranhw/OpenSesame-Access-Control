package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

func (s *ConfigService) IsSystemConfigured(ctx context.Context) (bool, error) {
	var rows int64
	if err := s.db.WithContext(ctx).Model(&db.SystemConfig{}).Count(&rows).Error; err != nil {
		return false, fmt.Errorf("counting system config entries: %w", err)
	}
	return rows > 0, nil
}

// GetSystemConfig retrieves the system configuration.
func (s *ConfigService) GetSystemConfig(ctx context.Context) (*db.SystemConfig, error) {
	var config db.SystemConfig
	// Assumes there's only one entry or the query implicitly selects one (e.g., by PK 1)
	// If multiple configs can exist and you need a specific one, adjust this query.
	if err := s.db.WithContext(ctx).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no configuration exists, not an error
		}
		return nil, fmt.Errorf("retrieving system config: %w", err)
	}
	return &config, nil
}

// CreateConfig initializes the system configuration.
func (s *ConfigService) CreateConfig(ctx context.Context, cfg *db.SystemConfig) error {
	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return fmt.Errorf("checking configuration status before create: %w", err)
	}
	if configured {
		return ErrAlreadyConfigured
	}

	if err := s.db.WithContext(ctx).Create(cfg).Error; err != nil {
		return fmt.Errorf("creating system config record: %w", err)
	}
	return nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, payload *dto.UpdateConfigPayload) (*db.SystemConfig, error) {
	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking configuration status for update: %w", err)
	}
	if !configured {
		return nil, ErrNotConfigured
	}

	if payload.SystemName == nil && payload.AdminPassword == nil {
		return nil, ErrNoUpdateFields
	}

	currentConfig, err := s.GetSystemConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching current config for update: %w", err)
	}
	if currentConfig == nil {
		// This case should ideally be prevented by IsSystemConfigured, but acts as a safeguard.
		return nil, ErrNotConfigured
	}

	// 4. Apply updates to the fetched configuration object.
	if payload.SystemName != nil {
		currentConfig.SystemName = *payload.SystemName
	}

	if payload.AdminPassword != nil {
		// Hash the new password.
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(*payload.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrPasswordHashingFailed, err)
		}
		currentConfig.AdminPasswordHash = string(newPasswordHash)
	}

	// 5. Save the modified configuration back to the database.
	if err := s.db.WithContext(ctx).Save(currentConfig).Error; err != nil {
		return nil, fmt.Errorf("saving updated config: %w", err)
	}

	// 6. Return the updated configuration.
	return currentConfig, nil
}

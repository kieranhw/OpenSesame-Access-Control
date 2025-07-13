package service

import (
	"context"

	"opensesame/internal/model"

	"gorm.io/gorm"
)

type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

func (s *ConfigService) IsSystemConfigured(ctx context.Context) (bool, error) {
	var rows int64
	if err := s.db.WithContext(ctx).Model(&model.SystemConfig{}).Count(&rows).Error; err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *ConfigService) GetSystemConfig(ctx context.Context) (*model.SystemConfig, error) {
	var config model.SystemConfig
	if err := s.db.WithContext(ctx).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (s *ConfigService) CreateConfig(ctx context.Context, cfg *model.SystemConfig) error {
	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return err
	}
	if configured {
		return ErrAlreadyConfigured
	}

	return s.db.WithContext(ctx).Create(cfg).Error
}

func (s *ConfigService) UpdateConfig(ctx context.Context, cfg *model.SystemConfig) error {
	configured, err := s.IsSystemConfigured(ctx)
	if err != nil {
		return err
	}
	if !configured {
		return ErrNotConfigured
	}

	return s.db.WithContext(ctx).Save(cfg).Error
}

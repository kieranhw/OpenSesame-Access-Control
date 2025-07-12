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

func (s *ConfigService) Create(ctx context.Context, si *model.SystemConfig) error {
	return s.db.WithContext(ctx).Create(si).Error
}

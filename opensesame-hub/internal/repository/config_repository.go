// repository/config_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"opensesame/internal/models/db"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	Count(ctx context.Context) (int64, error)
	GetSystemConfig(ctx context.Context) (*db.SystemConfig, error)
	CreateSystemConfig(ctx context.Context, cfg *db.SystemConfig) error
	UpdateSystemConfig(ctx context.Context, cfg *db.SystemConfig) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&db.SystemConfig{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("counting system config entries: %w", err)
	}
	return count, nil
}

func (r *configRepository) GetSystemConfig(ctx context.Context) (*db.SystemConfig, error) {
	var config db.SystemConfig
	if err := r.db.WithContext(ctx).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving system config: %w", err)
	}
	return &config, nil
}

func (r *configRepository) CreateSystemConfig(ctx context.Context, cfg *db.SystemConfig) error {
	if err := r.db.WithContext(ctx).Create(cfg).Error; err != nil {
		return fmt.Errorf("creating system config record: %w", err)
	}
	return nil
}

func (r *configRepository) UpdateSystemConfig(ctx context.Context, cfg *db.SystemConfig) error {
	if err := r.db.WithContext(ctx).Save(cfg).Error; err != nil {
		return fmt.Errorf("saving system config: %w", err)
	}
	return nil
}

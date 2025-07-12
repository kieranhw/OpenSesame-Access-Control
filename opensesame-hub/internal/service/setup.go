package service

import (
	"context"

	"opensesame/internal/model"

	"gorm.io/gorm"
)

// SetupService manages the singleton SystemInfo.
type SetupService struct {
	db *gorm.DB
}

func NewSetupService(db *gorm.DB) *SetupService {
	return &SetupService{db: db}
}

// Exists returns true if a row already exists.
func (s *SetupService) Exists(ctx context.Context) (bool, error) {
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&model.SystemInfo{}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create inserts the initial SystemInfo row.
func (s *SetupService) Create(ctx context.Context, si *model.SystemInfo) error {
	return s.db.WithContext(ctx).Create(si).Error
}

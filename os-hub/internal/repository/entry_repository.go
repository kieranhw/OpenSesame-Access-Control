package repository

import (
	"context"
	"errors"
	"fmt"

	"opensesame/internal/etag"
	"opensesame/internal/models/db"

	"gorm.io/gorm"
)

type EntryRepository interface {
	List(ctx context.Context) ([]*db.EntryDevice, error)
	GetEntryDeviceById(ctx context.Context, id uint) (*db.EntryDevice, error)
	CreateEntryDevice(ctx context.Context, entry *db.EntryDevice) error
	UpdateEntryDevice(ctx context.Context, entry *db.EntryDevice) error
}

type entryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) EntryRepository {
	return &entryRepository{db: db}
}

func (r *entryRepository) List(ctx context.Context) ([]*db.EntryDevice, error) {
	var devices []*db.EntryDevice
	if err := r.db.WithContext(ctx).
		Preload("Commands").
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("listing entry devices: %w", err)
	}
	return devices, nil
}

func (r *entryRepository) GetEntryDeviceById(ctx context.Context, id uint) (*db.EntryDevice, error) {
	var device db.EntryDevice
	if err := r.db.WithContext(ctx).
		Preload("Commands").
		First(&device, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving entry device %d: %w", id, err)
	}
	return &device, nil
}

func (r *entryRepository) CreateEntryDevice(ctx context.Context, entry *db.EntryDevice) error {
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return fmt.Errorf("creating entry device: %w", err)
	}
	etag.Bump()
	return nil
}

func (r *entryRepository) UpdateEntryDevice(ctx context.Context, entry *db.EntryDevice) error {
	if err := r.db.WithContext(ctx).Save(entry).Error; err != nil {
		return fmt.Errorf("updating entry device %d: %w", entry.EntryID, err)
	}
	etag.Bump()
	return nil
}

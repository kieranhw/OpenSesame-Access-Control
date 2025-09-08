package repository

import (
	"context"
	"errors"
	"fmt"

	"opensesame/internal/etag"
	"opensesame/internal/models/db"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	ListEntryDevices(ctx context.Context) ([]*db.EntryDevice, error)
	ListAllDevices(ctx context.Context) ([]*db.Device, error)
	GetEntryDeviceById(ctx context.Context, id uint) (*db.EntryDevice, error)
	GetEntryDeviceByMac(ctx context.Context, mac string) (*db.EntryDevice, error)
	CreateEntryDevice(ctx context.Context, entry *db.EntryDevice) error
	UpsertEntryDevice(ctx context.Context, entry *db.EntryDevice) error
	UpdateEntryDevice(ctx context.Context, id uint, fields map[string]interface{}) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) ListEntryDevices(ctx context.Context) ([]*db.EntryDevice, error) {
	var devices []*db.EntryDevice
	if err := r.db.WithContext(ctx).
		Preload("Device").
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("listing entry devices: %w", err)
	}
	return devices, nil
}

func (r *deviceRepository) ListAllDevices(ctx context.Context) ([]*db.Device, error) {
	var devices []*db.Device
	if err := r.db.WithContext(ctx).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("listing all devices: %w", err)
	}
	return devices, nil
}

func (r *deviceRepository) GetEntryDeviceById(ctx context.Context, id uint) (*db.EntryDevice, error) {
	var device db.EntryDevice
	if err := r.db.WithContext(ctx).
		Preload("Device").
		First(&device, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving entry device %d: %w", id, err)
	}
	return &device, nil
}
func (r *deviceRepository) GetEntryDeviceByMac(ctx context.Context, mac string) (*db.EntryDevice, error) {
	var device db.EntryDevice
	if err := r.db.WithContext(ctx).
		Preload("Device").
		Joins("Device").
		Where("Device.mac_address = ?", mac).
		First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving entry device by mac %s: %w", mac, err)
	}
	return &device, nil
}

func (r *deviceRepository) CreateEntryDevice(ctx context.Context, entry *db.EntryDevice) error {
	if err := r.db.WithContext(ctx).Create(entry).Error; err != nil {
		return fmt.Errorf("creating entry device: %w", err)
	}
	etag.Bump()
	return nil
}

func (r *deviceRepository) UpsertEntryDevice(ctx context.Context, entry *db.EntryDevice) error {
	if err := r.db.WithContext(ctx).Save(entry).Error; err != nil {
		return fmt.Errorf("upserting entry device %d: %w", entry.DeviceID, err)
	}
	etag.Bump()
	return nil
}
func (r *deviceRepository) UpdateEntryDevice(ctx context.Context, id uint, fields map[string]interface{}) error {
	if err := r.db.WithContext(ctx).
		Model(&db.Device{}).
		Where("id = ?", id).
		Updates(fields).Error; err != nil {
		return fmt.Errorf("updating entry device %d fields: %w", id, err)
	}

	etag.Bump()
	return nil
}

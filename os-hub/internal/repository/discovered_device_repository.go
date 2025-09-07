package repository

import (
	"context"
	"time"

	"opensesame/internal/models/db"

	"gorm.io/gorm"
)

type DiscoveredDeviceRepository interface {
	Upsert(ctx context.Context, device *db.DiscoveredDevice) error
	List(ctx context.Context) ([]db.DiscoveredDevice, error)
}

type discoveredDeviceRepository struct {
	db *gorm.DB
}

func NewDiscoveredDeviceRepository(db *gorm.DB) DiscoveredDeviceRepository {
	return &discoveredDeviceRepository{db: db}
}

func (r *discoveredDeviceRepository) Upsert(ctx context.Context, device *db.DiscoveredDevice) error {
	var existing db.DiscoveredDevice
	err := r.db.WithContext(ctx).
		Where("mac_address = ?", device.MacAddress).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		device.CreatedAt = time.Now()
		device.UpdatedAt = time.Now()
		// device is newly discovered, so insert and return
		return r.db.WithContext(ctx).Create(device).Error
	} else if err != nil {
		return err
	}

	// device already exists so update
	existing.InstanceName = device.InstanceName
	existing.DeviceType = device.DeviceType
	existing.IPAddress = device.IPAddress
	existing.Port = device.Port
	existing.ServiceType = device.ServiceType
	existing.LastSeen = time.Now()
	existing.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *discoveredDeviceRepository) List(ctx context.Context) ([]db.DiscoveredDevice, error) {
	var devices []db.DiscoveredDevice
	if err := r.db.WithContext(ctx).Find(&devices).Error; err != nil {
		return nil, err
	}

	return devices, nil
}

package db

import (
	"time"
)

type DiscoveredDevice struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress  string `gorm:"size:12;not null;index"`
	IPAddress   string `gorm:"size:45"`
	Port        int
	ServiceType string `gorm:"size:100"`

	DeviceType   string `gorm:"size:50"`
	InstanceType string `gorm:"size:32"`
	InstanceName string `gorm:"size:255"`

	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

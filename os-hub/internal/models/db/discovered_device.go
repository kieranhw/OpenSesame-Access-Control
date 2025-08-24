package db

import (
	"time"
)

type DiscoveredDevice struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress  string `gorm:"size:12;not null;index"`
	Hostname    string `gorm:"size:255"`
	Instance    string `gorm:"size:255"`
	DeviceType  string `gorm:"size:50"`
	IPv4        string `gorm:"size:45"`
	Port        int
	ServiceType string    `gorm:"size:100"`
	LastSeen    time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

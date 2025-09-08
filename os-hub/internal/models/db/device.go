package db

import (
	"time"
)

type Device struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	MacAddress   string    `gorm:"not null;uniqueIndex"`
	IPAddress    string    `gorm:"not null"`
	Port         int       `gorm:"not null"`
	Name         string    `gorm:"not null"`
	Description  *string   `gorm:"type:text"`
	ServiceType  *string   `gorm:"type:text"` // mDNS service type, e.g., "_http._tcp", will be null if not discovered via mDNS
	DeviceType   string    `gorm:"not null"`
	InstanceType string    `gorm:"not null"`
	InstanceName string    `gorm:"not null"`
	LastSeen     time.Time `gorm:"autoCreateTime"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

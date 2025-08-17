package db

import "time"

type SystemConfig struct {
	ID                int    `gorm:"primaryKey;autoIncrement"`
	SystemName        string `gorm:"not null"`
	SessionTimeoutSec int    `gorm:"not null;default:86400"` // default 24 hours
	AdminPasswordHash string `gorm:"not null"`
	BackupCodeHash    string `gorm:"not null"`
	SystemSecret      string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

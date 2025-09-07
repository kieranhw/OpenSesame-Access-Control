package db

import (
	"opensesame/internal/models/types"
	"time"
)

type EntryDevice struct {
	EntryID     uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress  string `gorm:"not null;uniqueIndex"`
	IPAddress   string `gorm:"not null"`
	Port        int    `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description *string

	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	DeviceType types.DeviceType `gorm:"not null"` // e.g. relay_lock
	LockStatus types.LockStatus `gorm:"type:text;not null;default:UNKNOWN"`
	Commands   []EntryCommand   `gorm:"foreignKey:EntryID;references:EntryID;constraint:OnDelete:CASCADE"`
}

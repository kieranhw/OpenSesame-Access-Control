package db

import (
	"time"
)

type EntryDevice struct {
	EntryID     uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress  string `gorm:"not null;uniqueIndex"`
	IPAddress   string `gorm:"not null"`
	Port        int    `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description *string
	LockStatus  LockStatus `gorm:"type:text;not null;default:UNKNOWN"`

	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Commands []EntryCommand `gorm:"foreignKey:EntryID;references:EntryID;constraint:OnDelete:CASCADE"`
}

type LockStatus string

const (
	LockStatusUnknown  LockStatus = "UNKNOWN"
	LockStatusLocked   LockStatus = "LOCKED"
	LockStatusUnlocked LockStatus = "UNLOCKED"
)

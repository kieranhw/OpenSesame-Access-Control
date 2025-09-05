package db

import "time"

type EntryDevice struct {
	EntryID    uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress string `gorm:"not null;uniqueIndex"`
	IPAddress  string `gorm:"not null"`
	Port       int    `gorm:"not null"`

	Name        string `gorm:"not null"`
	Description *string
	LockStatus  *string // e.g. "LOCKED", "UNLOCKED", "UNKNOWN"

	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// one-to-many device to commands
	Commands []EntryCommand `gorm:"foreignKey:EntryID;references:EntryID;constraint:OnDelete:CASCADE"`
}

type LockStatus string

const (
	Locked   LockStatus = "LOCKED"
	Unlocked LockStatus = "UNLOCKED"
	Unknown  LockStatus = "UNKNOWN"
)

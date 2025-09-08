package db

import (
	"opensesame/internal/models/types"
)

type EntryDevice struct {
	DeviceID   uint             `gorm:"primaryKey"`
	Device     Device           `gorm:"constraint:OnDelete:CASCADE"`
	LockStatus types.LockStatus `gorm:"type:text;not null;default:unknown"`
}

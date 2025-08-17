package db

import "time"

type Session struct {
	Token     string    `gorm:"primaryKey;size:36"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

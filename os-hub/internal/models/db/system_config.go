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

type EntryCommandDTO struct {
	ID        uint            `json:"id"`
	Type      string          `json:"type"`
	Status    string          `json:"status,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Http      *HttpCommandDTO `json:"http,omitempty"`
	// Udp      *UdpCommandDTO  `json:"udp,omitempty"` // TODO: future extension
}

type HttpCommandDTO struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

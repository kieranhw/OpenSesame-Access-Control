package db

import "time"

type EntryProtocol string

const (
	ProtocolHTTP EntryProtocol = "HTTP"
	ProtocolUDP  EntryProtocol = "UDP"
)

type EntryDevice struct {
	EntryID     uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	IP          string
	Port        int
	Description string
	Protocol    EntryProtocol `gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`

	// Relations
	Commands []EntryCommand `gorm:"foreignKey:EntryID"`
}

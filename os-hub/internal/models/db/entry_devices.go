package db

import "time"

type EntryProtocol string

const (
	HTTPWebhook EntryProtocol = "HTTP_WEBHOOK"
	UDPMessage  EntryProtocol = "UDP_MESSAGE"
)

type Entry struct {
	EntryID     uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	EntryType   EntryProtocol `gorm:"type:varchar(50);not null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`

	// Relations (optional, if you want eager loading)
	HttpWebhookEntry *HttpWebhookEntry `gorm:"foreignKey:EntryID"`
	UdpMessageEntry  *UdpMessageEntry  `gorm:"foreignKey:EntryID"`
}

type HttpWebhookEntry struct {
	EntryID uint   `gorm:"primaryKey"`
	URL     string `gorm:"not null"`
	Method  string `gorm:"not null"`
	Headers string // JSON or key-value
	Body    string

	// Relation back to Entry
	Entry Entry `gorm:"constraint:OnDelete:CASCADE;foreignKey:EntryID"`
}
type UdpMessageEntry struct {
	EntryID uint   `gorm:"primaryKey"`
	IP      string `gorm:"not null"`
	Port    int    `gorm:"not null"`
	Payload string

	// Relation back to Entry
	Entry Entry `gorm:"constraint:OnDelete:CASCADE;foreignKey:EntryID"`
}

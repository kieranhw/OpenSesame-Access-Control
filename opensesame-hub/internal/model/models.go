package model

import "time"

type Entry struct {
	EntryID     int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:255;not null"`
	Description string
	EntryType   string `gorm:"size:50;not null"` // e.g. "door","garage",etc.s
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ControlClient struct {
	ClientID        int       `gorm:"primaryKey;autoIncrement"`
	Name            string    `gorm:"not null"`
	RegistrationPin int       // 6-digit int
	GeneratedAt     time.Time // optional timestamp
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// many-to-many relationship via control_client_entries
	Entries []Entry `gorm:"many2many:control_client_entries;joinForeignKey:ClientID;JoinReferences:EntryID"`
}

type ControlClientEntry struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	ClientID int `gorm:"index;not null"`
	EntryID  int `gorm:"index;not null"`
}

type SystemConfig struct {
	ID                int    `gorm:"primaryKey;autoIncrement"`
	SystemName        string `gorm:"not null"`
	AdminPasswordHash string `gorm:"not null"`
	BackupCode        string `gorm:"not null"`
	SystemSecret      string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Session struct {
	Token     string    `gorm:"primaryKey;size:36"` // UUID v4
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}

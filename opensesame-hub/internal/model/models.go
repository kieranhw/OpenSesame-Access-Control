package model

import "time"

// Entry corresponds to your entries table
type Entry struct {
	EntryID     int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:255;not null"`
	Description string
	EntryType   string    `gorm:"size:50;not null"` // e.g. "door","garage",etc.
	CreatedAt   time.Time // GORM populates these automatically
	UpdatedAt   time.Time
}

// ControlClient is your clients table
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

// join table; GORM will create this automatically but
// we declare it to control the PK & references.
type ControlClientEntry struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	ClientID int `gorm:"index;not null"`
	EntryID  int `gorm:"index;not null"`
}

// SystemInfo is your singleton settings table
type SystemInfo struct {
	ID                int    `gorm:"primaryKey;autoIncrement"`
	SystemName        string `gorm:"not null"`
	AdminPasswordHash string `gorm:"not null"`
	BackupCodeHash    string
	SystemSecret      string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

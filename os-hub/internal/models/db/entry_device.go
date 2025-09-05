package db

import "time"

type EntryDevice struct {
	EntryID    uint   `gorm:"primaryKey;autoIncrement"`
	MacAddress string `gorm:"not null;uniqueIndex"`
	IPAddress  string `gorm:"not null"`
	Port       int    `gorm:"not null"`

	Name        string `gorm:"not null"`
	Description *string

	LastSeen  *time.Time `gorm:"index"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`

	// one-to-many device to commands
	Commands []EntryCommand `gorm:"foreignKey:EntryID;references:EntryID;constraint:OnDelete:CASCADE"`
}

type CommandType string

const (
	CommandLock   CommandType = "LOCK"
	CommandUnlock CommandType = "UNLOCK"
	CommandFail   CommandType = "CMD_FAIL"
)

type CommandStatus string

const (
	StatusPending CommandStatus = "PENDING"
	StatusSent    CommandStatus = "SENT"
	StatusSuccess CommandStatus = "SUCCESS"
	StatusFailed  CommandStatus = "FAILED"
)

type EntryCommand struct {
	CommandID   uint          `gorm:"primaryKey;autoIncrement"`
	EntryID     uint          `gorm:"not null;index"` // FK to EntryDevice
	CommandType CommandType   `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	Status      CommandStatus `gorm:"type:varchar(20)"`

	// HTTP fields
	URL     string `gorm:"not null"`
	Method  string `gorm:"not null"`
	Headers string
	Body    string
}

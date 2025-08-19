package db

import "time"

// --------------------
// EntryDevice (Parent)
// --------------------
type EntryDevice struct {
	EntryID     uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	IP          string
	Port        int
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// One-to-many: Device → Commands
	Commands []EntryCommand `gorm:"foreignKey:EntryID;references:EntryID;constraint:OnDelete:CASCADE"`
}

// --------------------
// EntryCommand (Child)
// --------------------
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

	// HTTP fields (inlined)
	URL     string `gorm:"not null"`
	Method  string `gorm:"not null"`
	Headers string
	Body    string
}

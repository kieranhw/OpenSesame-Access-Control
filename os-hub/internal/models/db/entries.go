package db

import "time"

type EntryProtocol string

const (
	ProtocolHTTP EntryProtocol = "HTTP"
	ProtocolUDP  EntryProtocol = "UDP"
)

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

type EntryCommand struct {
	CommandID   uint          `gorm:"primaryKey;autoIncrement"`
	EntryID     uint          `gorm:"not null"`
	CommandType CommandType   `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	Status      CommandStatus `gorm:"type:varchar(20)"`

	// Relations
	EntryDevice EntryDevice  `gorm:"foreignKey:EntryID"`
	HttpCommand *HttpCommand `gorm:"foreignKey:CommandID"`
	//UdpCommand  *UdpCommand  `gorm:"foreignKey:CommandID"`
}

type HttpCommand struct {
	CommandID uint   `gorm:"primaryKey"`
	URL       string `gorm:"not null"`
	Method    string `gorm:"not null"`
	Headers   string
	Body      string

	// Relation back to EntryCommand
	EntryCommand EntryCommand `gorm:"foreignKey:CommandID;constraint:OnDelete:CASCADE"`
}

// Not used for now
// type UdpCommand struct {
// 	CommandID       uint   `gorm:"primaryKey"`
// 	DestinationIP   string `gorm:"not null"`
// 	DestinationPort int    `gorm:"not null"`
// 	Payload         string

// 	// Relation back to EntryCommand
// 	EntryCommand EntryCommand `gorm:"foreignKey:CommandID;constraint:OnDelete:CASCADE"`
// }

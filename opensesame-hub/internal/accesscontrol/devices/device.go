package devices

// Command represents a generic command to control a device.
type Command struct {
	DoorID string
	Action string // e.g., "lock" or "unlock"
}

// DeviceHandler defines behavior for processing a device command.
type DeviceHandler interface {
	ProcessCommand(cmd Command) error
}

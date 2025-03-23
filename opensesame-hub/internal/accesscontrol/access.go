package accesscontrol

import "fmt"

// Command represents a door command.
type Command struct {
	DoorID string
	Action string // "lock" or "unlock"
}

// ProcessCommand processes a given door command.
func ProcessCommand(cmd Command) error {
	fmt.Printf("Processing command: Door %s, Action: %s\n", cmd.DoorID, cmd.Action)
	// Return nil for now (simulate success)
	return nil
}

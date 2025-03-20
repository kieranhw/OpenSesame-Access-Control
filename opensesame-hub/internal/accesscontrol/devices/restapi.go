// internal/accesscontrol/devices/restapi.go
package devices

import "fmt"

// RestAPIHandler handles commands received via REST API.
type RestAPIHandler struct{}

// ProcessCommand processes a command for REST API.
func (r RestAPIHandler) ProcessCommand(cmd Command) error {
	// Validate and forward the command to the core access control logic.
	fmt.Printf("REST API processing command for door %s: %s\n", cmd.DoorID, cmd.Action)
	return nil
}

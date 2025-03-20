package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"opensesame/internal/accesscontrol"
)

// DoorCommand represents the JSON structure of a door command.
type DoorCommand struct {
	DoorID string `json:"door_id"`
	Action string `json:"action"` // "lock" or "unlock"
}

// DoorHandler handles REST API requests for door commands.
func DoorHandler(w http.ResponseWriter, r *http.Request) {
	// Expect URLs like /doors/{door_id}/unlock or /doors/{door_id}/lock.
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(segments) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	doorID := segments[1]
	action := segments[2]

	cmd := accesscontrol.Command{
		DoorID: doorID,
		Action: action,
	}

	// In a real application, you would add authentication and authorization here.

	if err := accesscontrol.ProcessCommand(cmd); err != nil {
		http.Error(w, "Command processing failed", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status": "success",
		"door":   doorID,
		"action": action,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

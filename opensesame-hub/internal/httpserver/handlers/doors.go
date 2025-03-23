package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"opensesame/internal/accesscontrol"
)

// RegisterDoorRoutes sets up the routes for door commands.
func RegisterDoorRoutes(mux *http.ServeMux) {
	// Example: endpoints like /v1/doors/{lock_id}/lock or /v1/doors/{lock_id}/unlock
	mux.HandleFunc("/v1/doors/", DoorHandler)
}

// DoorHandler handles requests for door commands.
func DoorHandler(w http.ResponseWriter, r *http.Request) {
	// Simple parsing: expect URL like /v1/doors/{lock_id}/{action}
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	lockID := segments[1]
	action := segments[2]

	// Build the command and process it.
	cmd := accesscontrol.Command{
		DoorID: lockID,
		Action: action,
	}
	if err := accesscontrol.ProcessCommand(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status": "success",
		"door":   lockID,
		"action": action,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package handler

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/service"
)

// StartPairingHandler starts the pairing process
func StartPairingHandler(w http.ResponseWriter, r *http.Request) {
	var pairingData struct {
		DeviceID string `json:"device_id"`
	}

	// Decode the incoming request body
	if err := json.NewDecoder(r.Body).Decode(&pairingData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service layer to handle pairing
	result := service.StartPairing(pairingData.DeviceID)

	// Respond based on pairing result
	if result {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pairing started successfully")
	} else {
		http.Error(w, "Failed to start pairing", http.StatusInternalServerError)
	}
}

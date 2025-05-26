package handlers

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/service"
)

func StartPairingHandler(w http.ResponseWriter, r *http.Request) {
	var pairingData struct {
		DeviceID string `json:"device_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pairingData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result := service.StartPairing(pairingData.DeviceID)

	if result {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pairing started successfully")
	} else {
		http.Error(w, "Failed to start pairing", http.StatusInternalServerError)
	}
}

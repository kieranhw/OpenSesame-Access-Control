package handler

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/service"
)

// ValidatePinHandler validates the pin
func ValidatePinHandler(w http.ResponseWriter, r *http.Request) {
	var pin struct {
		Pin string `json:"pin"`
	}

	// Decode the incoming request body
	if err := json.NewDecoder(r.Body).Decode(&pin); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service layer to validate the pin
	isValid := service.ValidatePin(pin.Pin)

	// Respond based on validation result
	if isValid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pin validated successfully")
	} else {
		http.Error(w, "Invalid pin", http.StatusUnauthorized)
	}
}

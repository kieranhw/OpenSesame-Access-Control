package handlers

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/service"
)

func GetAccessInfo(w http.ResponseWriter, r *http.Request) {

}

func ValidatePinHandler(w http.ResponseWriter, r *http.Request) {
	var pin struct {
		Pin string `json:"pin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&pin); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	isValid := service.ValidatePin(pin.Pin)

	if isValid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pin validated successfully")
	} else {
		http.Error(w, "Invalid pin", http.StatusUnauthorized)
	}
}

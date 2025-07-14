package management

import (
	"encoding/json"
	"net/http"
)

type Access struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
	Enabled  bool   `json:"enabled"`
}

func GetAccessHandler(w http.ResponseWriter, r *http.Request) {
	dummy := []Access{
		{ID: 1, Location: "FrontDoor", Enabled: true},
		{ID: 2, Location: "BackDoor", Enabled: false},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dummy)
}

func PostAccessHandler(w http.ResponseWriter, r *http.Request) {
	var req Access
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.ID = 42
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

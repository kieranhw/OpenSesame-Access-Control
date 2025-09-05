package handlers

import (
	"encoding/json"
	"net/http"
	"opensesame/internal/models/dto"
	"opensesame/internal/service"
)

func ListEntryDevices(svc *service.EntryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		devices, err := svc.ListEntryDevices(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if devices == nil {
			devices = []dto.EntryDevice{}
		}

		if err := json.NewEncoder(w).Encode(devices); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateEntryDevice(svc *service.EntryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req dto.CreateEntryDeviceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device, err := svc.CreateEntryDevice(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(device); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

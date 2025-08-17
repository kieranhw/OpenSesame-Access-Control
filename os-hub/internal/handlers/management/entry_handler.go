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
			devices = []dto.EntryDeviceDTO{}
		}

		if err := json.NewEncoder(w).Encode(devices); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

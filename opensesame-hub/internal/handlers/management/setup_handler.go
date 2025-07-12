package management

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/model"
	"opensesame/internal/service"
)

// SystemSetupHandler returns an http.HandlerFunc that branches on GET vs POST
func SystemSetupHandler(svc *service.SetupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// query whether a row exists
			exists, err := svc.Exists(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"configured": exists})

		case http.MethodPost:
			// read the SystemInfo payload
			var in model.SystemInfo
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// TODO: hash in.AdminPasswordHash, in.BackupCodeHash,
			// generate in.SystemSecret if blank, validate etc.
			if err := svc.Create(r.Context(), &in); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(in)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

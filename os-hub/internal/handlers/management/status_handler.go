package handlers

import (
	"encoding/json"
	"net/http"

	"opensesame/internal/models/types"
)

func GetStatus(svcs *types.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		status, err := svcs.Status.GetStatus(ctx)
		if err != nil {
			http.Error(w, "failed to get status: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	}
}

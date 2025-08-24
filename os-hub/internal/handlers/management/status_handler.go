package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"opensesame/internal/models/types"
)

func GetStatus(svcs *types.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse ETag from query param (string -> uint64)
		var clientETag uint64
		if etagStr := r.URL.Query().Get("etag"); etagStr != "" {
			if v, err := strconv.ParseUint(etagStr, 10, 64); err == nil {
				clientETag = v
			}
		}

		// Parse timeout from query param
		var timeout time.Duration
		if t := r.URL.Query().Get("timeout"); t != "" {
			// Try parsing with units first (e.g. "30s", "2m")
			if dur, err := time.ParseDuration(t); err == nil {
				timeout = dur
			} else {
				// If no unit, assume seconds
				if secs, err2 := strconv.Atoi(t); err2 == nil {
					timeout = time.Duration(secs) * time.Second
				}
			}
		} else {
			timeout = 0 // no timeout param means return instantly
		}

		status, newETag, changed, err := svcs.Status.WaitForStatus(ctx, clientETag, timeout)
		if err != nil {
			http.Error(w, "failed to get status: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !changed {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("ETag", fmt.Sprintf("%d", newETag))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	}
}

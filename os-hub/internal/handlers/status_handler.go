package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"opensesame/internal/service"
	"strconv"
	"time"
)

func GetStatus(svcs *service.ServicesType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// parse ETag from query param
		var clientETag uint64
		if etagStr := r.URL.Query().Get("etag"); etagStr != "" {
			if v, err := strconv.ParseUint(etagStr, 10, 64); err == nil {
				clientETag = v
			}
		}

		// parse timeout from query param
		var timeout time.Duration
		if t := r.URL.Query().Get("timeout"); t != "" {
			if secs, err := strconv.Atoi(t); err == nil {
				timeout = time.Duration(secs) * time.Second
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

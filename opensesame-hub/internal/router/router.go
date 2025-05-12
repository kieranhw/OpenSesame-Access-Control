package router

import (
	"net/http"

	"opensesame/internal/handler"
)

func AddRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/access/validate_pin", handler.ValidatePinHandler)
	mux.HandleFunc("/pairing/start", handler.StartPairingHandler)

	return mux
}

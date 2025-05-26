package router

import (
	"net/http"

	"opensesame/internal/handlers"

	"github.com/gorilla/mux"
)

func AddRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/access/{access_id}/validate_pin", handlers.ValidatePinHandler).Methods("POST")
	r.HandleFunc("/pairing/start", handlers.StartPairingHandler).Methods("POST")

	return r
}

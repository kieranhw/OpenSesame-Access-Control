package router

import (
	"net/http"

	"opensesame/internal/handlers/management"

	"github.com/gorilla/mux"
)

func AddRoutes() http.Handler {
	r := mux.NewRouter()

	management.MountRoutes(r)

	return r
}

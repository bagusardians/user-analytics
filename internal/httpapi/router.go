package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/v1/logins", h.IngestLogin)
	return r
}

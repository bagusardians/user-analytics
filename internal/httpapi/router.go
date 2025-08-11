package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handlers) http.Handler {
	r := chi.NewRouter()
	r.Post("/v1/logins", h.IngestLogin)
	r.Get("/v1/user/uniques/daily", h.GetDailyUniqueUsers)
	r.Get("/v1/user/uniques/monthly", h.GetMonthlyUniqueUsers)
	return r
}

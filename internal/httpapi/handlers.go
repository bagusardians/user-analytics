package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
	"user-analytics/internal/service"

	"github.com/google/uuid"
)

type Handlers struct {
	svc *service.UserService
}

func NewHandlers(svc *service.UserService) *Handlers { return &Handlers{svc: svc} }

type ingestReq struct {
	UserID    string `json:"user_id"`
	LoginTime string `json:"login_time"`
	TZ        string `json:"tz,omitempty"`
}

func (h *Handlers) IngestLogin(w http.ResponseWriter, r *http.Request) {
	var req ingestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	t, err := time.Parse(time.RFC3339, req.LoginTime)
	if err != nil {
		http.Error(w, "invalid login_time", http.StatusBadRequest)
		return
	}
	ts := t.UTC()
	tz := req.TZ
	if tz == "" {
		tz = "UTC"
	}

	if err := h.svc.IngestLogin(r.Context(), uid, ts, tz); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) GetDailyUniqueUsers(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date required", http.StatusBadRequest)
		return
	}
	tz := r.URL.Query().Get("tz")
	if tz == "" {
		tz = "UTC"
	}

	day, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	n, err := h.svc.GetDailyUniqueUsers(r.Context(), day, tz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"date": dateStr, "tz": tz, "unique_users": n})
}

package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	IngestLogin(ctx context.Context, uid uuid.UUID, at time.Time) error
	GetDailyUniqueUsers(ctx context.Context, day time.Time) (int, error)
	GetMonthlyUniqueUsers(ctx context.Context, month time.Time) (int, error)
}

type Handlers struct {
	svc UserService
}

func NewHandlers(svc UserService) *Handlers { return &Handlers{svc: svc} }

type ingestReq struct {
	UserID    string `json:"user_id"`
	LoginTime string `json:"login_time"`
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

	if err := h.svc.IngestLogin(r.Context(), uid, ts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) GetDailyUniqueUsers(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date required", http.StatusBadRequest)
		return
	}

	day, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	n, err := h.svc.GetDailyUniqueUsers(r.Context(), day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"date": dateStr, "unique_users": n})
}

func (h *Handlers) GetMonthlyUniqueUsers(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	if monthStr == "" {
		http.Error(w, "month required", http.StatusBadRequest)
		return
	}

	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		http.Error(w, "invalid month", http.StatusBadRequest)
		return
	}

	n, err := h.svc.GetMonthlyUniqueUsers(r.Context(), month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"month": monthStr, "unique_users": n})
}

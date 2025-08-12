package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockSvc struct {
	ingestCalled bool
	gotUID       uuid.UUID
	gotAt        time.Time
	ingestErr    error
	daily        int
	dErr         error
	monthly      int
	mErr         error
}

func (m *mockSvc) IngestLogin(ctx context.Context, uid uuid.UUID, at time.Time) error {
	m.ingestCalled = true
	m.gotUID = uid
	m.gotAt = at
	return m.ingestErr
}
func (m *mockSvc) GetDailyUniqueUsers(ctx context.Context, day time.Time) (int, error) {
	return m.daily, m.dErr
}
func (m *mockSvc) GetMonthlyUniqueUsers(ctx context.Context, month time.Time) (int, error) {
	return m.monthly, m.mErr
}

func TestIngestLogin_SuccessUTCConversion(t *testing.T) {
	ms := &mockSvc{}
	h := NewHandlers(ms)

	body := []byte(`{"user_id":"05378ca8-961d-49e7-a903-8026dad78bb7","login_time":"2025-08-12T15:04:05+07:00"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.IngestLogin(w, req)

	if got, want := w.Code, http.StatusCreated; got != want {
		t.Fatalf("status: got %d want %d body=%s", got, want, w.Body.String())
	}
	if !ms.ingestCalled {
		t.Fatal("expected IngestLogin to be called")
	}
	wantUTC := time.Date(2025, 8, 12, 8, 4, 5, 0, time.UTC) // +07:00 -> UTC
	if !ms.gotAt.Equal(wantUTC) {
		t.Fatalf("UTC conversion: got %v want %v", ms.gotAt, wantUTC)
	}
}

func TestIngestLogin_InvalidJSON(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewReader([]byte(`{`)))
	w := httptest.NewRecorder()
	h.IngestLogin(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestIngestLogin_InvalidUserID(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	body := []byte(`{"user_id":"notuuid","login_time":"2025-08-12T15:04:05+07:00"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.IngestLogin(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestIngestLogin_InvalidLoginTime(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	body := []byte(`{"user_id":"05378ca8-961d-49e7-a903-8026dad78bb7","login_time":"notatimestamp"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.IngestLogin(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestIngestLogin_ServiceError(t *testing.T) {
	ms := &mockSvc{ingestErr: errors.New("error")}
	h := NewHandlers(ms)
	body := []byte(`{"user_id":"05378ca8-961d-49e7-a903-8026dad78bb7","login_time":"2025-08-12T15:04:05+07:00"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.IngestLogin(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d want 500", w.Code)
	}
}

func TestGetDailyUniqueUsers_MissingDate(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/daily", nil)
	w := httptest.NewRecorder()
	h.GetDailyUniqueUsers(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestGetDailyUniqueUsers_InvalidDateFormat(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/daily?date=2025/08/12", nil)
	w := httptest.NewRecorder()
	h.GetDailyUniqueUsers(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestGetDailyUniqueUsers_Success(t *testing.T) {
	ms := &mockSvc{daily: 100}
	h := NewHandlers(ms)
	q := url.Values{}
	q.Set("date", "2025-08-12")
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/daily?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	h.GetDailyUniqueUsers(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d want 200 body=%s", w.Code, w.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got["date"] != "2025-08-12" || int(got["unique_users"].(float64)) != 100 {
		t.Fatalf("payload: got %v", got)
	}
}

func TestGetDailyUniqueUsers_ServiceError(t *testing.T) {
	ms := &mockSvc{dErr: errors.New("error")}
	h := NewHandlers(ms)
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/daily?date=2025-08-12", nil)
	w := httptest.NewRecorder()
	h.GetDailyUniqueUsers(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d want 500", w.Code)
	}
}

func TestGetMonthlyUniqueUsers_MissingMonth(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/monthly", nil)
	w := httptest.NewRecorder()
	h.GetMonthlyUniqueUsers(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestGetMonthlyUniqueUsers_InvalidMonthFormat(t *testing.T) {
	h := NewHandlers(&mockSvc{})
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/monthly?month=2025-8", nil)
	w := httptest.NewRecorder()
	h.GetMonthlyUniqueUsers(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want 400", w.Code)
	}
}

func TestGetMonthlyUniqueUsers_Success(t *testing.T) {
	ms := &mockSvc{monthly: 5}
	h := NewHandlers(ms)
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/monthly?month=2025-08", nil)
	w := httptest.NewRecorder()
	h.GetMonthlyUniqueUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d want 200 body=%s", w.Code, w.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got["month"] != "2025-08" || int(got["unique_users"].(float64)) != 5 {
		t.Fatalf("payload: got %v", got)
	}
}

func TestGetMonthlyUniqueUsers_ServiceError(t *testing.T) {
	ms := &mockSvc{mErr: errors.New("error")}
	h := NewHandlers(ms)
	req := httptest.NewRequest(http.MethodGet, "/v1/user/uniques/monthly?month=2025-08", nil)
	w := httptest.NewRecorder()
	h.GetMonthlyUniqueUsers(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d want 500", w.Code)
	}
}

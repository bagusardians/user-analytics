package repo

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
)

func qm(s string) string { return regexp.QuoteMeta(s) }

func TestIngestLogin_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	defer mockPool.Close()
	r := NewUserPgRepo(mockPool)
	uid := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	ts := time.Date(2025, 8, 12, 8, 4, 5, 0, time.UTC)
	mockPool.ExpectBegin()
	mockPool.ExpectExec(qm(`
    INSERT INTO user_logins (user_id, login_time)
    VALUES ($1, $2)
    ON CONFLICT (user_id, login_time) DO NOTHING`)).
		WithArgs(uid, ts).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockPool.ExpectExec(qm(`
    INSERT INTO daily_unique_users (day, user_id)
    VALUES ( ( $1 )::date , $2 )
    ON CONFLICT DO NOTHING`)).
		WithArgs(ts, uid).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockPool.ExpectExec(qm(`
    INSERT INTO monthly_unique_users (month, user_id)
    VALUES ( date_trunc('month', $1::timestamptz )::date , $2 )
    ON CONFLICT DO NOTHING`)).
		WithArgs(ts, uid).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockPool.ExpectCommit()
	if err := r.IngestLogin(context.Background(), uid, ts); err != nil {
		t.Fatalf("IngestLogin err: %v", err)
	}
	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("doesn't meet expectations: %v", err)
	}
}

func TestGetDailyUniqueUsers_Success(t *testing.T) {
	mockPool, _ := pgxmock.NewPool()
	defer mockPool.Close()
	r := NewUserPgRepo(mockPool)
	day := time.Date(2025, 8, 12, 0, 0, 0, 0, time.UTC)
	mockPool.ExpectQuery(qm(`
    SELECT COUNT(*)::int FROM daily_unique_users
    WHERE day = ($1)::date`)).
		WithArgs(day).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(42))
	n, err := r.GetDailyUniqueUsers(context.Background(), day)
	if err != nil {
		t.Fatalf("GetDailyUniqueUsers err: %v", err)
	}
	if n != 42 {
		t.Fatalf("want 42, got %d", n)
	}
	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("doesn't meet expectations: %v", err)
	}
}

func TestGetMonthlyUniqueUsers_Success(t *testing.T) {
	mockPool, _ := pgxmock.NewPool()
	defer mockPool.Close()
	r := NewUserPgRepo(mockPool)
	month := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	mockPool.ExpectQuery(qm(`
    SELECT COUNT(*)::int FROM monthly_unique_users
    WHERE month = date_trunc('month', $1::timestamp)::date`)).
		WithArgs(month).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(7))
	n, err := r.GetMonthlyUniqueUsers(context.Background(), month)
	if err != nil {
		t.Fatalf("GetMonthlyUniqueUsers err: %v", err)
	}
	if n != 7 {
		t.Fatalf("want 7, got %d", n)
	}
	if err := mockPool.ExpectationsWereMet(); err != nil {
		t.Fatalf("doesn't meet expectations: %v", err)
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

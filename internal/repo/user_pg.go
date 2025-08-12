package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPool interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
type UserPgRepo struct{ db PgxPool }

func NewUserPgRepo(db PgxPool) *UserPgRepo { return &UserPgRepo{db: db} }

func (r *UserPgRepo) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time) error {

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err = r.db.Exec(ctx, `
    INSERT INTO user_logins (user_id, login_time)
    VALUES ($1, $2)
    ON CONFLICT (user_id, login_time) DO NOTHING
  `, userID, tsUTC); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, `
    INSERT INTO daily_unique_users (day, user_id)
    VALUES ( ( $1 )::date , $2 )
    ON CONFLICT DO NOTHING
  `, tsUTC, userID); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, `
    INSERT INTO monthly_unique_users (month, user_id)
    VALUES ( date_trunc('month', $1::timestamptz )::date , $2 )
    ON CONFLICT DO NOTHING
  `, tsUTC, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserPgRepo) GetDailyUniqueUsers(ctx context.Context, day time.Time) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
    SELECT COUNT(*)::int FROM daily_unique_users
    WHERE day = ($1)::date
  `, day).Scan(&n)
	return n, err
}

func (r *UserPgRepo) GetMonthlyUniqueUsers(ctx context.Context, month time.Time) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
    SELECT COUNT(*)::int FROM monthly_unique_users
    WHERE month = date_trunc('month', $1::timestamp)::date
  `, month).Scan(&n)
	return n, err
}

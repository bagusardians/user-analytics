package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPgRepo struct{ db *pgxpool.Pool }

func NewUserPgRepo(db *pgxpool.Pool) *UserPgRepo { return &UserPgRepo{db: db} }

func (r *UserPgRepo) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time, tz string) error {

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
    VALUES ( ( $1 AT TIME ZONE $2 )::date , $3 )
    ON CONFLICT DO NOTHING
  `, tsUTC, tz, userID); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, `
    INSERT INTO monthly_unique_users (month, user_id)
    VALUES ( date_trunc('month', $1 AT TIME ZONE $2 )::date , $3 )
    ON CONFLICT DO NOTHING
  `, tsUTC, tz, userID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserPgRepo) GetDailyUniqueUsers(ctx context.Context, day time.Time, tz string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
    SELECT COUNT(*)::int FROM daily_unique_users
    WHERE day = ($1 AT TIME ZONE $2)::date
  `, day, tz).Scan(&n)
	return n, err
}

func (r *UserPgRepo) GetMonthlyUniqueUsers(ctx context.Context, month time.Time, tz string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
    SELECT COUNT(*)::int FROM monthly_unique_users
    WHERE month = date_trunc('month', $1 AT TIME ZONE $2)::date
  `, month, tz).Scan(&n)
	return n, err
}

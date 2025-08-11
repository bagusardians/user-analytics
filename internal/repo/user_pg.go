package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPgRepo struct{ db *pgxpool.Pool }

func NewUserPgRepo(db *pgxpool.Pool) *UserPgRepo { return &UserPgRepo{db: db} }

func (r *UserPgRepo) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time, tz string) error {
	sql := `
    INSERT INTO user_logins (user_id, login_time)
    VALUES ($1, $2)
    ON CONFLICT (user_id, login_time) DO NOTHING
  `

	_, err := r.db.Exec(ctx, sql, userID, tsUTC)
	if err != nil {
		return fmt.Errorf("error inserting to db: %w", err)
	}
	return nil
}

func (r *UserPgRepo) GetDailyUniqueUsers(ctx context.Context, day time.Time, tz string) (int, error) {
	// TODO: retrieve from db
	fmt.Println("GetDailyUniqueUsers : day: ", day, ", tz: ", tz)
	return 100, nil
}

func (r *UserPgRepo) GetMonthlyUniqueUsers(ctx context.Context, month time.Time, tz string) (int, error) {
	// TODO: retrieve from db
	fmt.Println("GetMonthlyUniqueUsers : month: ", month, ", tz: ", tz)
	return 2000, nil
}

package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserPgRepo struct{}

func NewUserPgRepo() *UserPgRepo { return &UserPgRepo{} }

func (r *UserPgRepo) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time, tz string) error {
	// TODO: insert to db
	fmt.Println("login ingested with: userId: ", userID, ", tsUTC: ", tsUTC, ", tz: ", tz)
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

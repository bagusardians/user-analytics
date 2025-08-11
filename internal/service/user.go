package service

import (
	"context"
	"time"
	"user-analytics/internal/repo"

	"github.com/google/uuid"
)

type UserService struct{ repo *repo.UserPgRepo }

func NewUserService(r *repo.UserPgRepo) *UserService { return &UserService{repo: r} }

func (s *UserService) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time, tz string) error {
	return s.repo.IngestLogin(ctx, userID, tsUTC, tz)
}

func (s *UserService) GetDailyUniqueUsers(ctx context.Context, day time.Time, tz string) (int, error) {
	return s.repo.GetDailyUniqueUsers(ctx, day, tz)
}

func (s *UserService) GetMonthlyUniqueUsers(ctx context.Context, month time.Time, tz string) (int, error) {
	return s.repo.GetMonthlyUniqueUsers(ctx, month, tz)
}

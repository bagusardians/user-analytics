package service

import (
	"context"
	"time"
	"user-analytics/internal/repo"

	"github.com/google/uuid"
)

type UserService struct{ repo *repo.UserPgRepo }

func NewUserService(r *repo.UserPgRepo) *UserService { return &UserService{repo: r} }

func (s *UserService) IngestLogin(ctx context.Context, userID uuid.UUID, tsUTC time.Time) error {
	return s.repo.IngestLogin(ctx, userID, tsUTC)
}

func (s *UserService) GetDailyUniqueUsers(ctx context.Context, day time.Time) (int, error) {
	return s.repo.GetDailyUniqueUsers(ctx, day)
}

func (s *UserService) GetMonthlyUniqueUsers(ctx context.Context, month time.Time) (int, error) {
	return s.repo.GetMonthlyUniqueUsers(ctx, month)
}

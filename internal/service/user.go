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

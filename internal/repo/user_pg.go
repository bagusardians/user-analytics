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

package repository

import (
	"context"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

type MembershipRepository interface {
	Create(ctx context.Context, m *domain.Membership) error
}

// Postgres implementation

type membershipRepo struct{}

func NewMembershipRepository() MembershipRepository {
	return &membershipRepo{}
}

func (r *membershipRepo) Create(ctx context.Context, m *domain.Membership) error {
	d := db.GetPool()
	query := `INSERT INTO team_members (user_id, team_id, joined_at) VALUES ($1, $2, $3)`
	_, err := d.Exec(ctx, query, m.UserID, m.TeamID, m.JoinedAt)
	return err
}

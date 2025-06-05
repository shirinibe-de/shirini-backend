package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

type ClaimRepository interface {
	Create(ctx context.Context, c *domain.Claim) error
	GetByID(ctx context.Context, id string) (*domain.Claim, error)
	UpdateStatus(ctx context.Context, id string, status domain.ClaimStatus) error
}

// Postgres implementation

type claimRepo struct{}

func NewClaimRepository() ClaimRepository {
	return &claimRepo{}
}

func (r *claimRepo) Create(ctx context.Context, c *domain.Claim) error {
	d := db.GetPool()
	query := `INSERT INTO achievement_claims (id, team_id, claimed_by, claimed_for, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := d.Exec(ctx, query, c.ID, c.TeamID, c.ClaimedBy, c.ClaimedFor, c.Message, c.Status, c.CreatedAt)
	return err
}

func (r *claimRepo) GetByID(ctx context.Context, id string) (*domain.Claim, error) {
	d := db.GetPool()
	query := `SELECT id, team_id, claimed_by, claimed_for, message, status, created_at FROM achievement_claims WHERE id=$1`
	row := d.QueryRow(ctx, query, id)
	claim := &domain.Claim{}
	var status string
	err := row.Scan(&claim.ID, &claim.TeamID, &claim.ClaimedBy, &claim.ClaimedFor, &claim.Message, &status, &claim.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	claim.Status = domain.ClaimStatus(status)
	return claim, err
}

func (r *claimRepo) UpdateStatus(ctx context.Context, id string, status domain.ClaimStatus) error {
	d := db.GetPool()
	query := `UPDATE achievement_claims SET status=$1 WHERE id=$2`
	_, err := d.Exec(ctx, query, status, id)
	return err
}

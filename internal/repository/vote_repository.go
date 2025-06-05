package repository

import (
	"context"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

type VoteRepository interface {
	Create(ctx context.Context, v *domain.Vote) error
	CountVotes(ctx context.Context, claimID string) (up int, down int, err error)
	DeleteByClaimID(ctx context.Context, claimID string) error
}

// Postgres implementation

type voteRepo struct{}

func NewVoteRepository() VoteRepository {
	return &voteRepo{}
}

func (r *voteRepo) Create(ctx context.Context, v *domain.Vote) error {
	d := db.GetPool()
	query := `INSERT INTO votes (id, claim_id, voted_by, vote, voted_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := d.Exec(ctx, query, v.ID, v.ClaimID, v.VotedBy, v.Vote, v.VotedAt)
	return err
}

func (r *voteRepo) CountVotes(ctx context.Context, claimID string) (int, int, error) {
	d := db.GetPool()
	queryUp := `SELECT COUNT(*) FROM votes WHERE claim_id=$1 AND vote=true`
	row := d.QueryRow(ctx, queryUp, claimID)
	var up int
	err := row.Scan(&up)
	if err != nil {
		return 0, 0, err
	}
	queryDown := `SELECT COUNT(*) FROM votes WHERE claim_id=$1 AND vote=false`
	row2 := d.QueryRow(ctx, queryDown, claimID)
	var down int
	err = row2.Scan(&down)
	if err != nil {
		return up, 0, err
	}
	return up, down, nil
}

func (r *voteRepo) DeleteByClaimID(ctx context.Context, claimID string) error {
	d := db.GetPool()
	query := `DELETE FROM votes WHERE claim_id=$1`
	_, err := d.Exec(ctx, query, claimID)
	return err
}

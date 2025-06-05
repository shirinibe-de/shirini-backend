package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByJoinToken(ctx context.Context, token string) (*domain.Team, error)
}

// Postgres implementation

type teamRepo struct{}

func NewTeamRepository() TeamRepository {
	return &teamRepo{}
}

func (r *teamRepo) Create(ctx context.Context, team *domain.Team) error {
	d := db.GetPool()
	query := `INSERT INTO teams (id, name, join_token, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := d.Exec(ctx, query, team.ID, team.Name, team.JoinToken, team.CreatedBy, team.CreatedAt)
	return err
}

func (r *teamRepo) GetByJoinToken(ctx context.Context, token string) (*domain.Team, error) {
	d := db.GetPool()
	query := `SELECT id, name, join_token, created_by, created_at FROM teams WHERE join_token=$1`
	row := d.QueryRow(ctx, query, token)
	team := &domain.Team{}
	err := row.Scan(&team.ID, &team.Name, &team.JoinToken, &team.CreatedBy, &team.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return team, err
}

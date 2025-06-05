package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// Postgres implementation

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	d := db.GetPool()
	query := `INSERT INTO users (id, email, name, avatar_url, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := d.Exec(ctx, query, user.ID, user.Email, user.Name, user.AvatarURL, user.CreatedAt)
	return err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	d := db.GetPool()
	query := `SELECT id, email, name, avatar_url, created_at FROM users WHERE email=$1`
	row := d.QueryRow(ctx, query, email)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

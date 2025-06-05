package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shirinibe-de/shirini-backend/config"
)

var pool *pgxpool.Pool

func Init(cfg *config.Config) error {
	var err error
	pool, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func GetPool() *pgxpool.Pool {
	return pool
}

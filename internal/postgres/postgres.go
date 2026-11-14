package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

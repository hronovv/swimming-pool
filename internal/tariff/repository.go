package tariff

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]Tariff, error) {
	sqlQuery := `
		SELECT
		t.name, t.format, t.price, t.visit_limit, t.duration_days
		FROM tariffs t;
	`

	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	tariffs, err := pgx.CollectRows(rows, pgx.RowToStructByName[Tariff])
	if err != nil {
		return nil, err
	}

	return tariffs, nil
}

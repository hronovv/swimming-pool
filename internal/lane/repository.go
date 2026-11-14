package lane

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

func (r *Repository) GetAll(ctx context.Context) ([]Lane, error) {
	sqlQuery := `
		SELECT
		l.id, l.name, l.length_m, l.depth_m, l.status
		FROM lanes l;
	`

	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	lanes, err := pgx.CollectRows(rows, pgx.RowToStructByName[Lane])
	if err != nil {
		return nil, err
	}

	return lanes, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Lane, error) {
	sqlQuery := `
		SELECT
		l.id, l.name, l.length_m, l.depth_m, l.status
		FROM lanes l
		WHERE l.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Lane{}, err
	}

	lane, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Lane])
	if err != nil {
		return Lane{}, err
	}

	return lane, nil
}

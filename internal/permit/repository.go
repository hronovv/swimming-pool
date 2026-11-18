package permit

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Permit, error) {
	sqlQuery := `
		SELECT p.trainer_id, p.work_format_id, p.issued_on, p.expires_on
		FROM trainer_permits p
		WHERE ($1 = 0 OR p.trainer_id = $1)
		ORDER BY p.trainer_id, p.work_format_id;
	`
	rows, err := r.db.Query(ctx, sqlQuery, filter.TrainerID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Permit])
}

package trainer

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const trainerColumns = `
	t.id, t.user_id, t.qualification, t.rank, t.experience_years, t.specialization
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Trainer, error) {
	sqlQuery := `
		SELECT` + trainerColumns + `
		FROM trainers t
		WHERE ($1 = 0 OR t.user_id = $1)
		ORDER BY t.id;
	`
	rows, err := r.db.Query(ctx, sqlQuery, filter.UserID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Trainer])
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Trainer, error) {
	sqlQuery := `
		SELECT` + trainerColumns + `
		FROM trainers t
		WHERE t.id = $1;
	`
	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Trainer{}, err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Trainer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Trainer{}, ErrNotFound
		}
		return Trainer{}, err
	}
	return item, nil
}

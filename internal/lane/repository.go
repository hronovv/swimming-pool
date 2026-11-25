package lane

import (
	"context"
	"errors"
	"fmt"

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
		FROM lanes l
		ORDER BY l.id;
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
		if errors.Is(err, pgx.ErrNoRows) {
			return Lane{}, ErrNotFound
		}
		return Lane{}, err
	}

	return lane, nil
}

func (r *Repository) Create(ctx context.Context, lane Lane) (Lane, error) {
	sqlQuery := `
		INSERT INTO lanes (name, length_m, depth_m, status)
		VALUES
			($1, $2, $3, $4)
		RETURNING id;
	`
	err := r.db.QueryRow(ctx, sqlQuery, lane.Name, lane.LengthM, lane.DepthM, lane.Status).Scan(&lane.ID)
	if err != nil {
		return Lane{}, fmt.Errorf("failed to create lane: %w", err)
	}
	return lane, nil
}

func (r *Repository) Update(ctx context.Context, lane Lane) error {
	sqlQuery := `
		UPDATE lanes
		SET
			name = @name,
			length_m = @length_m,
			depth_m = @depth_m,
			status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":       lane.ID,
		"name":     lane.Name,
		"length_m": lane.LengthM,
		"depth_m":  lane.DepthM,
		"status":   lane.Status,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return fmt.Errorf("update lane: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	sqlQuery := `
		UPDATE lanes
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, status, id)
	if err != nil {
		return fmt.Errorf("update lane status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM lanes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete lane: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

package standard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const standardColumns = `
	s.id, s.client_id, s.trainer_id, s.style, s.distance_m, s.held_on,
	s.result_time, s.outcome, s.attempt_no, s.comment
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Standard, error) {
	sqlQuery := `
		SELECT` + standardColumns + `
		FROM standards s
		WHERE ($1 = 0 OR s.client_id = $1)
		  AND ($2 = 0 OR s.trainer_id = $2)
		ORDER BY s.held_on DESC, s.id DESC;
	`
	rows, err := r.db.Query(ctx, sqlQuery, filter.ClientID, filter.TrainerID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Standard])
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Standard, error) {
	sqlQuery := `
		SELECT` + standardColumns + `
		FROM standards s
		WHERE s.id = $1;
	`
	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Standard{}, err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Standard])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Standard{}, ErrNotFound
		}
		return Standard{}, err
	}
	return item, nil
}

func (r *Repository) Create(ctx context.Context, item Standard) (Standard, error) {
	sqlQuery := `
		INSERT INTO standards (
			client_id, trainer_id, style, distance_m, held_on,
			result_time, outcome, attempt_no, comment
		)
		VALUES ($1, $2, $3, $4, $5, $6::time, $7, $8, $9)
		RETURNING id;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		item.ClientID,
		item.TrainerID,
		item.Style,
		item.DistanceM,
		item.HeldOn,
		clockArg(item.ResultTime),
		item.Outcome,
		item.AttemptNo,
		item.Comment,
	).Scan(&item.ID)
	if err != nil {
		return Standard{}, mapRepoErr(err)
	}
	return item, nil
}

func mapRepoErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		switch pgErr.ConstraintName {
		case "fk_standards_client":
			return ErrClientNotFound
		case "fk_standards_trainer":
			return ErrTrainerNotFound
		}
	}
	return fmt.Errorf("failed to create standard: %w", err)
}

func clockArg(value *time.Time) any {
	if value == nil || value.IsZero() {
		return nil
	}
	return value.Format("15:04:05")
}

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
	s.result_time, s.outcome, s.attempt_no, s.comment,
	` + personNameSQL + ` AS client_name,
	trim(tu.last_name || ' ' || left(tu.first_name, 1) || '.' || CASE WHEN coalesce(tu.patronymic, '') = '' THEN '' ELSE ' ' || left(tu.patronymic, 1) || '.' END) AS trainer_name
`

type Repository struct {
	db *pgxpool.Pool
}

const personNameSQL = `trim(u.last_name || ' ' || left(u.first_name, 1) || '.' || CASE WHEN coalesce(u.patronymic, '') = '' THEN '' ELSE ' ' || left(u.patronymic, 1) || '.' END)`

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Standard, error) {
	sqlQuery := `
		SELECT` + standardColumns + `
		FROM standards s
		JOIN clients c ON c.id = s.client_id
		JOIN users u ON u.id = c.user_id
		JOIN trainers t ON t.id = s.trainer_id
		JOIN users tu ON tu.id = t.user_id
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

func (r *Repository) ClientIDByUser(ctx context.Context, userID int) (int, error) {
	return cardIDByUser(ctx, r.db, `SELECT id FROM clients WHERE user_id = $1`, userID)
}

func (r *Repository) TrainerIDByUser(ctx context.Context, userID int) (int, error) {
	return cardIDByUser(ctx, r.db, `SELECT id FROM trainers WHERE user_id = $1`, userID)
}

func cardIDByUser(ctx context.Context, db *pgxpool.Pool, sqlQuery string, userID int) (int, error) {
	var id int
	err := db.QueryRow(ctx, sqlQuery, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Standard, error) {
	sqlQuery := `
		SELECT` + standardColumns + `
		FROM standards s
		JOIN clients c ON c.id = s.client_id
		JOIN users u ON u.id = c.user_id
		JOIN trainers t ON t.id = s.trainer_id
		JOIN users tu ON tu.id = t.user_id
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

func (r *Repository) Update(ctx context.Context, item Standard) error {
	sqlQuery := `
		UPDATE standards
		SET client_id = $2, trainer_id = $3, style = $4, distance_m = $5, held_on = $6,
			result_time = $7::time, outcome = $8, attempt_no = $9, comment = $10
		WHERE id = $1
	`
	tag, err := r.db.Exec(
		ctx,
		sqlQuery,
		item.ID,
		item.ClientID,
		item.TrainerID,
		item.Style,
		item.DistanceM,
		item.HeldOn,
		clockArg(item.ResultTime),
		item.Outcome,
		item.AttemptNo,
		item.Comment,
	)
	if err != nil {
		return mapRepoErr(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM standards WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete standard: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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

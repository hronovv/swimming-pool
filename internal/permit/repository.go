package permit

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) TrainerIDByUser(ctx context.Context, userID int) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM trainers WHERE user_id = $1`, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
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

func (r *Repository) Create(ctx context.Context, item Permit) (Permit, error) {
	sqlQuery := `
		INSERT INTO trainer_permits (trainer_id, work_format_id, issued_on, expires_on)
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.Exec(ctx, sqlQuery, item.TrainerID, item.WorkFormatID, item.IssuedOn, item.ExpiresOn)
	if err != nil {
		return Permit{}, mapRepoErr(err)
	}
	return item, nil
}

func (r *Repository) Update(ctx context.Context, trainerID, workFormatID int, item Permit) error {
	sqlQuery := `
		UPDATE trainer_permits
		SET trainer_id = $3, work_format_id = $4, issued_on = $5, expires_on = $6
		WHERE trainer_id = $1 AND work_format_id = $2
	`
	tag, err := r.db.Exec(ctx, sqlQuery, trainerID, workFormatID, item.TrainerID, item.WorkFormatID, item.IssuedOn, item.ExpiresOn)
	if err != nil {
		return mapRepoErr(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, trainerID, workFormatID int) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM trainer_permits
		WHERE trainer_id = $1 AND work_format_id = $2
	`, trainerID, workFormatID)
	if err != nil {
		return fmt.Errorf("delete permit: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func mapRepoErr(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fmt.Errorf("failed to create permit: %w", err)
	}
	switch pgErr.Code {
	case "23505":
		return ErrConflict
	case "23503":
		if pgErr.ConstraintName == "fk_permits_trainer" {
			return ErrTrainerNotFound
		}
		return fmt.Errorf("%w: work format not found", ErrInvalid)
	case "23514":
		return fmt.Errorf("%w: permit is not allowed", ErrInvalid)
	default:
		return fmt.Errorf("failed to create permit: %w", err)
	}
}

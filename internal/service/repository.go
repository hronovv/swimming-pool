package service

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
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Model, error) {
	sqlQuery := `
		SELECT
		s.id, s.client_id, s.tariff_id, s.group_id, s.start_date, s.end_date, s.status
		FROM services s
		WHERE ($1 = 0 OR s.client_id = $1)
		  AND ($2 = 0 OR s.tariff_id = $2)
		  AND ($3 = 0 OR s.group_id = $3)
		  AND ($4 = '' OR s.status = $4)
		ORDER BY s.id;
	`

	rows, err := r.db.Query(ctx, sqlQuery, filter.ClientID, filter.TariffID, filter.GroupID, filter.Status)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[Model])
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Model, error) {
	sqlQuery := `
		SELECT
		s.id, s.client_id, s.tariff_id, s.group_id, s.start_date, s.end_date, s.status
		FROM services s
		WHERE s.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Model{}, err
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Model])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Model{}, ErrNotFound
		}
		return Model{}, err
	}

	return model, nil
}

func (r *Repository) Create(ctx context.Context, model Model) (Model, error) {
	sqlQuery := `
		INSERT INTO services (client_id, tariff_id, group_id, start_date, end_date, status)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		model.ClientID,
		model.TariffID,
		model.GroupID,
		model.StartDate,
		model.EndDate,
		model.Status,
	).Scan(&model.ID)
	if err != nil {
		return Model{}, mapRepoErr(err, "failed to create service")
	}
	return model, nil
}

func (r *Repository) Update(ctx context.Context, model Model) error {
	sqlQuery := `
		UPDATE services
		SET
			client_id = @client_id,
			tariff_id = @tariff_id,
			group_id = @group_id,
			start_date = @start_date,
			end_date = @end_date,
			status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":         model.ID,
		"client_id":  model.ClientID,
		"tariff_id":  model.TariffID,
		"group_id":   model.GroupID,
		"start_date": model.StartDate,
		"end_date":   model.EndDate,
		"status":     model.Status,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return mapRepoErr(err, "update service")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	sqlQuery := `
		UPDATE services
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, status, id)
	if err != nil {
		return fmt.Errorf("update service status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) TariffInfo(ctx context.Context, tariffID int) (TariffInfo, error) {
	sqlQuery := `
		SELECT t.format, t.work_format_id, t.duration_days
		FROM tariffs t
		WHERE t.id = $1;
	`

	var info TariffInfo
	err := r.db.QueryRow(ctx, sqlQuery, tariffID).Scan(&info.Format, &info.WorkFormatID, &info.DurationDays)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TariffInfo{}, ErrTariffNotFound
		}
		return TariffInfo{}, err
	}
	return info, nil
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		switch pgErr.ConstraintName {
		case "fk_services_client":
			return ErrClientNotFound
		case "fk_services_tariff":
			return ErrTariffNotFound
		case "fk_services_group":
			return ErrGroupNotFound
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

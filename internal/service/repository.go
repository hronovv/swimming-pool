package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
		v.service_id AS id, v.client_id, v.tariff_id, v.group_id, v.start_date, v.end_date, v.status,
		` + viewPersonNameSQL + ` AS client_name, v.tariff_name, v.work_format_id, v.group_name
		FROM v_services v
		WHERE ($1 = 0 OR v.client_id = $1)
		  AND ($2 = 0 OR v.tariff_id = $2)
		  AND ($3 = 0 OR v.group_id = $3)
		  AND ($4 = '' OR v.status = $4)
		ORDER BY v.service_id;
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
		v.service_id AS id, v.client_id, v.tariff_id, v.group_id, v.start_date, v.end_date, v.status,
		` + viewPersonNameSQL + ` AS client_name, v.tariff_name, v.work_format_id, v.group_name
		FROM v_services v
		WHERE v.service_id = $1;
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
	_, _ = r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_tariff_stats`)
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

func (r *Repository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM services WHERE id = $1`, id)
	if err != nil {
		return mapRepoErr(err, "delete service")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	_, _ = r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_tariff_stats`)
	return nil
}

const viewPersonNameSQL = `trim(v.last_name || ' ' || left(v.first_name, 1) || '.' || CASE WHEN coalesce(v.patronymic, '') = '' THEN '' ELSE ' ' || left(v.patronymic, 1) || '.' END)`

func (r *Repository) MedicalCertificate(ctx context.Context, clientID int) (bool, error) {
	var allowed bool
	err := r.db.QueryRow(ctx, `SELECT medical_certificate FROM clients WHERE id = $1`, clientID).Scan(&allowed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrClientNotFound
		}
		return false, err
	}
	return allowed, nil
}

func (r *Repository) ClientIDByUser(ctx context.Context, userID int) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM clients WHERE user_id = $1`, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
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
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			switch pgErr.ConstraintName {
			case "fk_services_client":
				return ErrClientNotFound
			case "fk_services_tariff":
				return ErrTariffNotFound
			case "fk_services_group":
				return ErrGroupNotFound
			default:
				if strings.HasPrefix(op, "delete") {
					return ErrInUse
				}
			}
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

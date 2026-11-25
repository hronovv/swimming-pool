package tariff

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

func (r *Repository) GetAll(ctx context.Context) ([]Tariff, error) {
	sqlQuery := `
		SELECT
		t.id, t.name, t.format, t.work_format_id, t.price, t.visit_limit, t.duration_days,
		coalesce(m.services_cnt, 0) AS services_cnt,
		coalesce(m.completed_payments_cnt, 0) AS completed_payments_cnt,
		coalesce(m.revenue, 0) AS revenue,
		coalesce(m.avg_check, 0) AS avg_check
		FROM tariffs t
		LEFT JOIN mv_tariff_stats m ON m.tariff_id = t.id
		ORDER BY t.id;
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

func (r *Repository) GetByID(ctx context.Context, ID int) (Tariff, error) {
	sqlQuery := `
		SELECT
		t.id, t.name, t.format, t.work_format_id, t.price, t.visit_limit, t.duration_days,
		coalesce(m.services_cnt, 0) AS services_cnt,
		coalesce(m.completed_payments_cnt, 0) AS completed_payments_cnt,
		coalesce(m.revenue, 0) AS revenue,
		coalesce(m.avg_check, 0) AS avg_check
		FROM tariffs t
		LEFT JOIN mv_tariff_stats m ON m.tariff_id = t.id
		WHERE t.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Tariff{}, err
	}

	tariff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Tariff])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tariff{}, ErrNotFound
		}
		return Tariff{}, err
	}

	return tariff, nil
}

func (r *Repository) Create(ctx context.Context, tariff Tariff) (Tariff, error) {
	sqlQuery := `
		INSERT INTO tariffs (name, format, work_format_id, price, visit_limit, duration_days)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		tariff.Name,
		tariff.Format,
		tariff.WorkFormatID,
		tariff.Price,
		tariff.VisitLimit,
		tariff.DurationDays,
	).Scan(&tariff.ID)
	if err != nil {
		return Tariff{}, mapRepoErr(err, "failed to create tariff")
	}
	_ = r.RefreshStats(ctx)
	return tariff, nil
}

func (r *Repository) Update(ctx context.Context, tariff Tariff) error {
	sqlQuery := `
		UPDATE tariffs
		SET
			name = @name,
			format = @format,
			work_format_id = @work_format_id,
			price = @price,
			visit_limit = @visit_limit,
			duration_days = @duration_days
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":             tariff.ID,
		"name":           tariff.Name,
		"format":         tariff.Format,
		"work_format_id": tariff.WorkFormatID,
		"price":          tariff.Price,
		"visit_limit":    tariff.VisitLimit,
		"duration_days":  tariff.DurationDays,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return mapRepoErr(err, "update tariff")
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	_ = r.RefreshStats(ctx)
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM tariffs WHERE id = $1`, id)
	if err != nil {
		return mapRepoErr(err, "delete tariff")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	_ = r.RefreshStats(ctx)
	return nil
}

func (r *Repository) RefreshStats(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_tariff_stats`)
	return err
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrConflict
		case "23503":
			if strings.HasPrefix(op, "delete") {
				return ErrInUse
			}
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

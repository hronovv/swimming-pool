package payment

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

const viewPersonNameSQL = `trim(v.last_name || ' ' || left(v.first_name, 1) || '.' || CASE WHEN coalesce(v.patronymic, '') = '' THEN '' ELSE ' ' || left(v.patronymic, 1) || '.' END)`

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Payment, error) {
	sqlQuery := `
		SELECT
		v.payment_id AS id, v.service_id, v.amount, v.paid_on, v.method, v.purpose, v.status,
		v.client_id, ` + viewPersonNameSQL + ` AS client_name
		FROM v_payments v
		WHERE ($1 = 0 OR v.service_id = $1)
		  AND ($2 = '' OR v.status = $2)
		  AND ($3 = '' OR v.method = $3)
		  AND ($4 = 0 OR v.client_id = $4)
		ORDER BY v.payment_id;
	`

	rows, err := r.db.Query(ctx, sqlQuery, filter.ServiceID, filter.Status, filter.Method, filter.ClientID)
	if err != nil {
		return nil, err
	}

	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payment])
	if err != nil {
		return nil, err
	}

	return payments, nil
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

func (r *Repository) GetByID(ctx context.Context, ID int) (Payment, error) {
	sqlQuery := `
		SELECT
		v.payment_id AS id, v.service_id, v.amount, v.paid_on, v.method, v.purpose, v.status,
		v.client_id, ` + viewPersonNameSQL + ` AS client_name
		FROM v_payments v
		WHERE v.payment_id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Payment{}, err
	}

	payment, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Payment])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Payment{}, ErrNotFound
		}
		return Payment{}, err
	}

	return payment, nil
}

func (r *Repository) Create(ctx context.Context, payment Payment) (Payment, error) {
	sqlQuery := `
		INSERT INTO payments (service_id, amount, paid_on, method, purpose, status)
		VALUES
			($1, $2, COALESCE($3::date, CURRENT_DATE), $4, $5, $6)
		RETURNING id, paid_on;
	`

	var paidOn any
	if !payment.PaidOn.IsZero() {
		paidOn = payment.PaidOn
	}

	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		payment.ServiceID,
		payment.Amount,
		paidOn,
		payment.Method,
		payment.Purpose,
		payment.Status,
	).Scan(&payment.ID, &payment.PaidOn)
	if err != nil {
		return Payment{}, mapRepoErr(err, "failed to create payment")
	}
	_ = r.refreshTariffStats(ctx)
	return payment, nil
}

func (r *Repository) Update(ctx context.Context, payment Payment) error {
	sqlQuery := `
		UPDATE payments
		SET
			amount = @amount,
			paid_on = @paid_on,
			method = @method,
			purpose = @purpose,
			status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":      payment.ID,
		"amount":  payment.Amount,
		"paid_on": payment.PaidOn,
		"method":  payment.Method,
		"purpose": payment.Purpose,
		"status":  payment.Status,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return mapRepoErr(err, "update payment")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	_ = r.refreshTariffStats(ctx)
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	sqlQuery := `
		UPDATE payments
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, status, id)
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	_ = r.refreshTariffStats(ctx)
	return nil
}

func (r *Repository) refreshTariffStats(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_tariff_stats`)
	return err
}

func (r *Repository) ServiceFormat(ctx context.Context, serviceID int) (string, error) {
	sqlQuery := `
		SELECT t.format
		FROM services s
		JOIN tariffs t ON t.id = s.tariff_id
		WHERE s.id = $1;
	`

	var format string
	err := r.db.QueryRow(ctx, sqlQuery, serviceID).Scan(&format)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrServiceNotFound
		}
		return "", err
	}
	return format, nil
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrServiceNotFound
	}
	return fmt.Errorf("%s: %w", op, err)
}

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

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Payment, error) {
	sqlQuery := `
		SELECT
		p.id, p.service_id, p.amount, p.paid_on, p.method, p.purpose, p.status
		FROM payments p
		WHERE ($1 = 0 OR p.service_id = $1)
		  AND ($2 = '' OR p.status = $2)
		  AND ($3 = '' OR p.method = $3)
		ORDER BY p.id;
	`

	rows, err := r.db.Query(ctx, sqlQuery, filter.ServiceID, filter.Status, filter.Method)
	if err != nil {
		return nil, err
	}

	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payment])
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Payment, error) {
	sqlQuery := `
		SELECT
		p.id, p.service_id, p.amount, p.paid_on, p.method, p.purpose, p.status
		FROM payments p
		WHERE p.id = $1;
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
	return nil
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

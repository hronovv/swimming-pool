package client

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const clientColumns = `
	c.id, c.user_id, c.passport, c.birth_date, c.address, c.medical_certificate,
	trim(u.last_name || ' ' || left(u.first_name, 1) || '.' || CASE WHEN coalesce(u.patronymic, '') = '' THEN '' ELSE ' ' || left(u.patronymic, 1) || '.' END) AS name
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Client, error) {
	sqlQuery := `
		SELECT` + clientColumns + `
		FROM clients c
		JOIN users u ON u.id = c.user_id
		WHERE ($1 = 0 OR c.user_id = $1)
		ORDER BY c.id;
	`
	rows, err := r.db.Query(ctx, sqlQuery, filter.UserID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Client])
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Client, error) {
	sqlQuery := `
		SELECT` + clientColumns + `
		FROM clients c
		JOIN users u ON u.id = c.user_id
		WHERE c.id = $1;
	`
	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Client{}, err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Client])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Client{}, ErrNotFound
		}
		return Client{}, err
	}
	return item, nil
}

func (r *Repository) UpdateMedical(ctx context.Context, id int, allowed bool) error {
	tag, err := r.db.Exec(ctx, `UPDATE clients SET medical_certificate = $1 WHERE id = $2`, allowed, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

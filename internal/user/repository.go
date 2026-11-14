package user

import (
	"context"

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

func (r *Repository) GetByID(ctx context.Context, ID int) (User, error) {
	sqlQuery := `
		SELECT
		u.id, u.last_name, u.first_name, u.patronymic, u.role, u.status, u.registered_at, u.phone, u.email
		FROM users u 
		WHERE id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return User{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, err
	}

	return user, nil
}

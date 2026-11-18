package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const userColumns = `
	u.id, u.login, u.password_hash, u.last_name, u.first_name, u.patronymic,
	u.role, u.status, u.registered_at, u.phone, u.email
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]User, error) {
	sqlQuery := `
		SELECT` + userColumns + `
		FROM users u
		WHERE ($1 = '' OR u.role = $1)
		  AND ($2 = '' OR u.status = $2)
		ORDER BY u.id;
	`

	rows, err := r.db.Query(ctx, sqlQuery, filter.Role, filter.Status)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (User, error) {
	sqlQuery := `
		SELECT` + userColumns + `
		FROM users u
		WHERE u.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return User{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (User, error) {
	sqlQuery := `
		SELECT` + userColumns + `
		FROM users u
		WHERE u.login = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, login)
	if err != nil {
		return User{}, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, user User) (User, error) {
	sqlQuery := `
		INSERT INTO users (
			login, password_hash, last_name, first_name, patronymic,
			role, status, phone, email
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, registered_at;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		user.Login,
		user.PasswordHash,
		user.LastName,
		user.FirstName,
		user.Patronymic,
		user.Role,
		user.Status,
		user.Phone,
		user.Email,
	).Scan(&user.ID, &user.RegisteredAt)
	if err != nil {
		return User{}, mapRepoErr(err, "failed to create user")
	}
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user User) error {
	sqlQuery := `
		UPDATE users
		SET
			login = @login,
			last_name = @last_name,
			first_name = @first_name,
			patronymic = @patronymic,
			role = @role,
			phone = @phone,
			email = @email
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":         user.ID,
		"login":      user.Login,
		"last_name":  user.LastName,
		"first_name": user.FirstName,
		"patronymic": user.Patronymic,
		"role":       user.Role,
		"phone":      user.Phone,
		"email":      user.Email,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return mapRepoErr(err, "update user")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	sqlQuery := `
		UPDATE users
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, status, id)
	if err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	sqlQuery := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, passwordHash, id)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return fmt.Errorf("%s: %w", op, err)
}

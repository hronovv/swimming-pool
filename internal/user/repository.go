package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (r *Repository) Create(ctx context.Context, user User, profile Profile) (User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return User{}, fmt.Errorf("begin user create: %w", err)
	}
	defer tx.Rollback(ctx)

	sqlQuery := `
		INSERT INTO users (
			login, password_hash, last_name, first_name, patronymic,
			role, status, phone, email
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, registered_at;
	`
	err = tx.QueryRow(
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

	switch user.Role {
	case RoleClient:
		var address *string
		if profile.Address != "" {
			address = &profile.Address
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO clients (user_id, passport, birth_date, address, medical_certificate)
			VALUES ($1, $2, $3, $4, $5)
		`, user.ID, profile.Passport, profile.BirthDate, address, profile.MedicalCertificate)
	case RoleTrainer:
		_, err = tx.Exec(ctx, `
			INSERT INTO trainers (user_id, qualification, rank, experience_years, specialization)
			VALUES ($1, $2, $3, $4, $5)
		`, user.ID, profile.Qualification, profile.Rank, profile.ExperienceYears, profile.Specialization)
	}
	if err != nil {
		return User{}, mapRepoErr(err, "failed to create user profile")
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, fmt.Errorf("commit user create: %w", err)
	}
	return user, nil
}

func (r *Repository) Card(ctx context.Context, userID int) (int, int, error) {
	var clientID, trainerID int
	err := r.db.QueryRow(ctx, `
		SELECT
			COALESCE((SELECT id FROM clients WHERE user_id = $1), 0),
			COALESCE((SELECT id FROM trainers WHERE user_id = $1), 0)
	`, userID).Scan(&clientID, &trainerID)
	if err != nil {
		return 0, 0, fmt.Errorf("user card: %w", err)
	}
	return clientID, trainerID, nil
}

func (r *Repository) ChangeRole(ctx context.Context, user User, profile Profile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin role change: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldRole string
	err = tx.QueryRow(ctx, `SELECT role FROM users WHERE id = $1`, user.ID).Scan(&oldRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("read user role: %w", err)
	}
	if err := updateUser(ctx, tx, user); err != nil {
		return err
	}
	if oldRole == RoleClient && user.Role != RoleClient {
		if _, err := tx.Exec(ctx, `DELETE FROM clients WHERE user_id = $1`, user.ID); err != nil {
			return cardInUse(err, "client card is still in use")
		}
	}
	if oldRole == RoleTrainer && user.Role != RoleTrainer {
		if _, err := tx.Exec(ctx, `DELETE FROM trainers WHERE user_id = $1`, user.ID); err != nil {
			return cardInUse(err, "trainer card is still in use")
		}
	}
	switch user.Role {
	case RoleClient:
		var address *string
		if profile.Address != "" {
			address = &profile.Address
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO clients (user_id, passport, birth_date, address, medical_certificate)
			VALUES ($1, $2, $3, $4, $5)
		`, user.ID, profile.Passport, profile.BirthDate, address, profile.MedicalCertificate)
	case RoleTrainer:
		_, err = tx.Exec(ctx, `
			INSERT INTO trainers (user_id, qualification, rank, experience_years, specialization)
			VALUES ($1, $2, $3, $4, $5)
		`, user.ID, profile.Qualification, profile.Rank, profile.ExperienceYears, profile.Specialization)
	}
	if err != nil {
		return mapRepoErr(err, "failed to create user profile")
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit role change: %w", err)
	}
	return nil
}

func updateUser(ctx context.Context, tx pgx.Tx, user User) error {
	tag, err := tx.Exec(ctx, `
		UPDATE users
		SET login = $1, last_name = $2, first_name = $3, patronymic = $4,
			role = $5, phone = $6, email = $7
		WHERE id = $8
	`, user.Login, user.LastName, user.FirstName, user.Patronymic, user.Role, user.Phone, user.Email, user.ID)
	if err != nil {
		return mapRepoErr(err, "failed to update user")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func cardInUse(err error, message string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23001") {
		return fmt.Errorf("%w: %s", ErrInvalid, message)
	}
	return fmt.Errorf("delete profile card: %w", err)
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
		if strings.Contains(pgErr.ConstraintName, "passport") {
			return ErrPassportTaken
		}
		return ErrConflict
	}
	return fmt.Errorf("%s: %w", op, err)
}

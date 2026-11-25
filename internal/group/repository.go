package group

import (
	"context"
	"errors"
	"fmt"

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

func (r *Repository) Members(ctx context.Context, groupID int) ([]Member, error) {
	sqlQuery := `
		SELECT
		v.client_id, v.login, v.last_name, v.first_name, v.patronymic,
		v.service_id, v.service_status
		FROM v_client_groups v
		WHERE v.group_id = $1
		ORDER BY v.last_name, v.first_name, v.service_id;
	`
	rows, err := r.db.Query(ctx, sqlQuery, groupID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Member])
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Group, error) {
	sqlQuery := `
		SELECT
		g.id, g.name, g.start_date, g.end_date, g.capacity, g.status
		FROM training_groups g
		WHERE ($1 = '' OR g.status = $1)
		ORDER BY g.id;
	`

	rows, err := r.db.Query(ctx, sqlQuery, filter.Status)
	if err != nil {
		return nil, err
	}

	groups, err := pgx.CollectRows(rows, pgx.RowToStructByName[Group])
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Group, error) {
	sqlQuery := `
		SELECT
		g.id, g.name, g.start_date, g.end_date, g.capacity, g.status
		FROM training_groups g
		WHERE g.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Group{}, err
	}

	group, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Group])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Group{}, ErrNotFound
		}
		return Group{}, err
	}

	return group, nil
}

func (r *Repository) Create(ctx context.Context, group Group) (Group, error) {
	sqlQuery := `
		INSERT INTO training_groups (name, start_date, end_date, capacity, status)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		group.Name,
		group.StartDate,
		group.EndDate,
		group.Capacity,
		group.Status,
	).Scan(&group.ID)
	if err != nil {
		return Group{}, fmt.Errorf("failed to create group: %w", err)
	}
	return group, nil
}

func (r *Repository) Update(ctx context.Context, group Group) error {
	sqlQuery := `
		UPDATE training_groups
		SET
			name = @name,
			start_date = @start_date,
			end_date = @end_date,
			capacity = @capacity,
			status = @status
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":         group.ID,
		"name":       group.Name,
		"start_date": group.StartDate,
		"end_date":   group.EndDate,
		"capacity":   group.Capacity,
		"status":     group.Status,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return fmt.Errorf("update group: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	sqlQuery := `
		UPDATE training_groups
		SET status = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, status, id)
	if err != nil {
		return fmt.Errorf("update group status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

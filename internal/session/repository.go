package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const sessionColumns = `
	s.id, s.work_format_id, s.group_id, s.lane_id, s.session_date,
	s.start_time, s.end_time, s.topic, s.was_held
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Session, error) {
	sqlQuery := `
		SELECT` + sessionColumns + `
		FROM sessions s
		WHERE ($1 = 0 OR s.work_format_id = $1)
		  AND ($2 = 0 OR s.group_id = $2)
		  AND ($3 = 0 OR s.lane_id = $3)
		  AND ($4::boolean IS NULL OR s.was_held = $4)
		  AND ($5::date IS NULL OR s.session_date >= $5)
		  AND ($6::date IS NULL OR s.session_date <= $6)
		ORDER BY s.session_date, s.start_time, s.id;
	`

	rows, err := r.db.Query(
		ctx,
		sqlQuery,
		filter.WorkFormatID,
		filter.GroupID,
		filter.LaneID,
		filter.WasHeld,
		nullableTime(filter.From),
		nullableTime(filter.To),
	)
	if err != nil {
		return nil, err
	}

	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Session, error) {
	sqlQuery := `
		SELECT` + sessionColumns + `
		FROM sessions s
		WHERE s.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Session{}, err
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Session])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Session{}, ErrNotFound
		}
		return Session{}, err
	}

	return item, nil
}

func (r *Repository) Create(ctx context.Context, item Session) (Session, error) {
	sqlQuery := `
		INSERT INTO sessions (
			work_format_id, group_id, lane_id, session_date,
			start_time, end_time, topic, was_held
		)
		VALUES
			($1, $2, $3, $4, $5::time, $6::time, $7, $8)
		RETURNING id;
	`
	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		item.WorkFormatID,
		item.GroupID,
		item.LaneID,
		item.SessionDate,
		clockArg(item.StartTime),
		clockArg(item.EndTime),
		item.Topic,
		item.WasHeld,
	).Scan(&item.ID)
	if err != nil {
		return Session{}, mapRepoErr(err, "failed to create session")
	}
	return item, nil
}

func (r *Repository) Update(ctx context.Context, item Session) error {
	sqlQuery := `
		UPDATE sessions
		SET
			work_format_id = @work_format_id,
			group_id = @group_id,
			lane_id = @lane_id,
			session_date = @session_date,
			start_time = @start_time::time,
			end_time = @end_time::time,
			topic = @topic,
			was_held = @was_held
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":             item.ID,
		"work_format_id": item.WorkFormatID,
		"group_id":       item.GroupID,
		"lane_id":        item.LaneID,
		"session_date":   item.SessionDate,
		"start_time":     clockArg(item.StartTime),
		"end_time":       clockArg(item.EndTime),
		"topic":          item.Topic,
		"was_held":       item.WasHeld,
	}

	tag, err := r.db.Exec(ctx, sqlQuery, args)
	if err != nil {
		return mapRepoErr(err, "update session")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateHeld(ctx context.Context, id int, wasHeld bool) error {
	sqlQuery := `
		UPDATE sessions
		SET was_held = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(ctx, sqlQuery, wasHeld, id)
	if err != nil {
		return fmt.Errorf("update session held: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		switch pgErr.ConstraintName {
		case "fk_sessions_work_format":
			return ErrWorkFormatNotFound
		case "fk_sessions_group":
			return ErrGroupNotFound
		case "fk_sessions_lane":
			return ErrLaneNotFound
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func clockArg(value time.Time) string {
	return value.Format("15:04:05")
}

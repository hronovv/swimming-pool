package journal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const entryColumns = `
	al.id, al.user_id, u.login, al.occurred_at, al.action,
	al.object, al.description, al.result, al.ip_address
`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Entry, error) {
	sqlQuery := `
		SELECT` + entryColumns + `
		FROM action_logs al
		JOIN users u ON u.id = al.user_id
		WHERE ($1 = 0 OR al.user_id = $1)
		  AND ($2 = '' OR al.action = $2)
		  AND ($3 = '' OR al.result = $3)
		  AND ($4::timestamp IS NULL OR al.occurred_at >= $4)
		  AND ($5::timestamp IS NULL OR al.occurred_at < $5)
		  AND ($6 = '' OR al.description ILIKE '%' || $6 || '%')
		ORDER BY al.occurred_at DESC;
	`

	rows, err := r.db.Query(
		ctx,
		sqlQuery,
		filter.UserID,
		filter.Action,
		filter.Result,
		nullableTime(filter.From),
		nullableTime(filter.To),
		filter.Query,
	)
	if err != nil {
		return nil, err
	}

	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entry])
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *Repository) GetByID(ctx context.Context, ID int) (Entry, error) {
	sqlQuery := `
		SELECT` + entryColumns + `
		FROM action_logs al
		JOIN users u ON u.id = al.user_id
		WHERE al.id = $1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, ID)
	if err != nil {
		return Entry{}, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Entry])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Entry{}, ErrNotFound
		}
		return Entry{}, err
	}

	return entry, nil
}

func (r *Repository) Create(ctx context.Context, entry Entry) (Entry, error) {
	procedure, err := procedureName(entry.Action)
	if err != nil {
		return Entry{}, err
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Entry{}, fmt.Errorf("begin journal entry: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `CALL `+procedure+`($1, $2, $3, $4, $5)`,
		entry.UserID,
		entry.Object,
		entry.Description,
		entry.Result,
		entry.IPAddress,
	)
	if err != nil {
		return Entry{}, mapRepoErr(err)
	}
	err = tx.QueryRow(ctx, `
		SELECT id, occurred_at
		FROM action_logs
		WHERE id = currval(pg_get_serial_sequence('action_logs', 'id'))
	`).Scan(&entry.ID, &entry.OccurredAt)
	if err != nil {
		return Entry{}, fmt.Errorf("read journal entry: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return Entry{}, fmt.Errorf("commit journal entry: %w", err)
	}
	return entry, nil
}

func procedureName(action string) (string, error) {
	switch action {
	case ActionLogin:
		return "log_authorization", nil
	case ActionCreate:
		return "log_create", nil
	case ActionUpdate:
		return "log_update", nil
	case ActionDelete:
		return "log_delete", nil
	default:
		return "", fmt.Errorf("%w: action must be login, create, update or delete", ErrInvalid)
	}
}

func mapRepoErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrUserNotFound
	}
	return fmt.Errorf("failed to create journal entry: %w", err)
}

func nullableTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

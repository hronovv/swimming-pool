package enrollment

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const personNameSQL = `trim(u.last_name || ' ' || left(u.first_name, 1) || '.' || CASE WHEN coalesce(u.patronymic, '') = '' THEN '' ELSE ' ' || left(u.patronymic, 1) || '.' END)`

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, filter Filter) ([]Enrollment, error) {
	sqlQuery := `
		SELECT e.service_id, e.session_id, e.attended, e.comment, s.client_id,
			` + personNameSQL + ` AS client_name, t.name AS tariff_name
		FROM session_enrollments e
		JOIN services s ON s.id = e.service_id
		JOIN clients c ON c.id = s.client_id
		JOIN users u ON u.id = c.user_id
		JOIN tariffs t ON t.id = s.tariff_id
		WHERE ($1 = 0 OR e.service_id = $1)
		  AND ($2 = 0 OR e.session_id = $2)
		  AND ($3 = 0 OR c.user_id = $3)
		ORDER BY e.session_id, e.service_id;
	`
	rows, err := r.db.Query(ctx, sqlQuery, filter.ServiceID, filter.SessionID, filter.ClientUserID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Enrollment])
}

func (r *Repository) Create(ctx context.Context, item Enrollment) (Enrollment, error) {
	sqlQuery := `
		INSERT INTO session_enrollments (service_id, session_id, attended, comment)
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.Exec(ctx, sqlQuery, item.ServiceID, item.SessionID, item.Attended, item.Comment)
	if err != nil {
		return Enrollment{}, mapRepoErr(err, "failed to create enrollment")
	}
	return item, nil
}

func (r *Repository) Delete(ctx context.Context, serviceID, sessionID int) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM session_enrollments
		WHERE service_id = $1 AND session_id = $2
	`, serviceID, sessionID)
	if err != nil {
		return fmt.Errorf("delete enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateAttended(ctx context.Context, serviceID, sessionID int, attended bool) error {
	sqlQuery := `
		UPDATE session_enrollments
		SET attended = $1
		WHERE service_id = $2 AND session_id = $3
	`
	tag, err := r.db.Exec(ctx, sqlQuery, attended, serviceID, sessionID)
	if err != nil {
		return fmt.Errorf("update enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ServiceLink(ctx context.Context, serviceID int) (Link, error) {
	sqlQuery := `
		SELECT t.work_format_id, s.group_id, s.status, s.start_date, s.end_date, t.visit_limit,
			(SELECT count(*)::int FROM session_enrollments e WHERE e.service_id = s.id),
			c.user_id
		FROM services s
		JOIN tariffs t ON t.id = s.tariff_id
		JOIN clients c ON c.id = s.client_id
		WHERE s.id = $1;
	`
	var link Link
	err := r.db.QueryRow(ctx, sqlQuery, serviceID).Scan(
		&link.WorkFormatID,
		&link.GroupID,
		&link.Status,
		&link.StartDate,
		&link.EndDate,
		&link.VisitLimit,
		&link.UsedVisits,
		&link.ClientUserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Link{}, ErrServiceNotFound
		}
		return Link{}, err
	}
	return link, nil
}

func (r *Repository) SessionLink(ctx context.Context, sessionID int) (Link, error) {
	sqlQuery := `
		SELECT s.work_format_id, s.group_id, s.session_date
		FROM sessions s
		WHERE s.id = $1;
	`
	var link Link
	err := r.db.QueryRow(ctx, sqlQuery, sessionID).Scan(&link.WorkFormatID, &link.GroupID, &link.SessionDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Link{}, ErrSessionNotFound
		}
		return Link{}, err
	}
	return link, nil
}

func mapRepoErr(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrConflict
		case "23503":
			switch pgErr.ConstraintName {
			case "fk_enrollments_service":
				return ErrServiceNotFound
			case "fk_enrollments_session":
				return ErrSessionNotFound
			}
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}

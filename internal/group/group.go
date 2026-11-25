package group

import (
	"errors"
	"time"
)

const (
	StatusRecruiting = "recruiting"
	StatusTraining   = "training"
	StatusFinished   = "finished"
)

var (
	ErrNotFound     = errors.New("group not found")
	ErrInvalid      = errors.New("invalid group")
	ErrInUse        = errors.New("group is still in use")
	allowedStatuses = map[string]struct{}{
		StatusRecruiting: {},
		StatusTraining:   {},
		StatusFinished:   {},
	}
)

type Group struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Capacity  int       `db:"capacity"`
	Status    string    `db:"status"`
}

type Member struct {
	ClientID      int    `db:"client_id"`
	Login         string `db:"login"`
	LastName      string `db:"last_name"`
	FirstName     string `db:"first_name"`
	Patronymic    *string `db:"patronymic"`
	ServiceID     int    `db:"service_id"`
	ServiceStatus string `db:"service_status"`
}

type Filter struct {
	Status string
}

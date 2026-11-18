package service

import (
	"errors"
	"time"
)

const (
	StatusActive = "active"
	StatusFrozen = "frozen"
	StatusClosed = "closed"

	workFormatFreeSwim   = 1
	workFormatIndividual = 2
	workFormatGroup      = 3
	workFormatSport      = 4
)

var (
	ErrNotFound       = errors.New("service not found")
	ErrInvalid        = errors.New("invalid service")
	ErrClientNotFound = errors.New("client not found")
	ErrTariffNotFound = errors.New("tariff not found")
	ErrGroupNotFound  = errors.New("group not found")
	allowedStatuses   = map[string]struct{}{
		StatusActive: {}, StatusFrozen: {}, StatusClosed: {},
	}
)

type Model struct {
	ID        int       `db:"id"`
	ClientID  int       `db:"client_id"`
	TariffID  int       `db:"tariff_id"`
	GroupID   *int      `db:"group_id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Status    string    `db:"status"`
}

type Filter struct {
	ClientID int
	TariffID int
	GroupID  int
	Status   string
}

type TariffInfo struct {
	Format       string
	WorkFormatID int
	DurationDays int
}

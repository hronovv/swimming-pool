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
	ErrInUse          = errors.New("service is still in use")
	allowedStatuses   = map[string]struct{}{
		StatusActive: {}, StatusFrozen: {}, StatusClosed: {},
	}
)

type Model struct {
	ID           int       `db:"id"`
	ClientID     int       `db:"client_id"`
	TariffID     int       `db:"tariff_id"`
	GroupID      *int      `db:"group_id"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	Status       string    `db:"status"`
	ClientName   string    `db:"client_name"`
	TariffName   string    `db:"tariff_name"`
	WorkFormatID int       `db:"work_format_id"`
	GroupName    *string   `db:"group_name"`
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

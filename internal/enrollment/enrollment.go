package enrollment

import (
	"errors"
	"time"
)

const maxCommentLen = 300

var (
	ErrNotFound        = errors.New("enrollment not found")
	ErrInvalid         = errors.New("invalid enrollment")
	ErrConflict        = errors.New("enrollment already exists")
	ErrServiceNotFound = errors.New("service not found")
	ErrSessionNotFound = errors.New("session not found")
)

type Enrollment struct {
	ServiceID   int     `db:"service_id"`
	SessionID   int     `db:"session_id"`
	Attended    bool    `db:"attended"`
	Comment     *string `db:"comment"`
	ClientName  string  `db:"client_name"`
	TariffName  string  `db:"tariff_name"`
	ActorUserID int     `db:"-"`
}

type Link struct {
	WorkFormatID int
	GroupID      *int
	Status       string
	StartDate    time.Time
	EndDate      time.Time
	SessionDate  time.Time
	VisitLimit   *int
	UsedVisits   int
	ClientUserID int
}

type Filter struct {
	ServiceID    int
	SessionID    int
	ClientUserID int
}

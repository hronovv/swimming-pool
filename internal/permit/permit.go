package permit

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("permit not found")
	ErrInvalid  = errors.New("invalid permit")
)

type Permit struct {
	TrainerID    int        `db:"trainer_id"`
	WorkFormatID int        `db:"work_format_id"`
	IssuedOn     time.Time  `db:"issued_on"`
	ExpiresOn    *time.Time `db:"expires_on"`
}

type Filter struct {
	TrainerID int
}

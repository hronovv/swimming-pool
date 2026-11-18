package standard

import (
	"errors"
	"time"
)

const (
	StyleFreestyle    = "freestyle"
	StyleBreaststroke = "breaststroke"
	StyleButterfly    = "butterfly"
	StyleBackstroke   = "backstroke"
	StyleMedley       = "medley"

	OutcomePassed = "passed"
	OutcomeFailed = "failed"

	maxCommentLen = 300
)

var (
	ErrNotFound        = errors.New("standard not found")
	ErrInvalid         = errors.New("invalid standard")
	ErrClientNotFound  = errors.New("client not found")
	ErrTrainerNotFound = errors.New("trainer not found")
	allowedStyles      = map[string]struct{}{
		StyleFreestyle: {}, StyleBreaststroke: {}, StyleButterfly: {}, StyleBackstroke: {}, StyleMedley: {},
	}
	allowedDistances = map[int]struct{}{25: {}, 50: {}, 100: {}, 200: {}, 400: {}}
	allowedOutcomes  = map[string]struct{}{OutcomePassed: {}, OutcomeFailed: {}}
)

type Standard struct {
	ID         int        `db:"id"`
	ClientID   int        `db:"client_id"`
	TrainerID  int        `db:"trainer_id"`
	Style      string     `db:"style"`
	DistanceM  int        `db:"distance_m"`
	HeldOn     time.Time  `db:"held_on"`
	ResultTime *time.Time `db:"result_time"`
	Outcome    string     `db:"outcome"`
	AttemptNo  int        `db:"attempt_no"`
	Comment    *string    `db:"comment"`
}

type Filter struct {
	ClientID  int
	TrainerID int
}

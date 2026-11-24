package session

import (
	"errors"
	"time"
)

const (
	workFormatFreeSwim   = 1
	workFormatIndividual = 2
	workFormatGroup      = 3
	workFormatSport      = 4

	maxTopicLen = 150
)

var (
	ErrNotFound           = errors.New("session not found")
	ErrInvalid            = errors.New("invalid session")
	ErrWorkFormatNotFound = errors.New("work format not found")
	ErrGroupNotFound      = errors.New("group not found")
	ErrLaneNotFound       = errors.New("lane not found")
	ErrTrainerNotFound    = errors.New("trainer not found")
	allowedWorkFormats    = map[int]struct{}{
		workFormatFreeSwim:   {},
		workFormatIndividual: {},
		workFormatGroup:      {},
		workFormatSport:      {},
	}
)

type Session struct {
	ID           int       `db:"id"`
	WorkFormatID int       `db:"work_format_id"`
	GroupID      *int      `db:"group_id"`
	LaneID       *int      `db:"lane_id"`
	SessionDate  time.Time `db:"session_date"`
	StartTime    time.Time `db:"start_time"`
	EndTime      time.Time `db:"end_time"`
	Topic        *string   `db:"topic"`
	WasHeld      bool      `db:"was_held"`
	GroupName    *string   `db:"group_name"`
	LaneName     *string   `db:"lane_name"`
	TrainerIDs   []int     `db:"-"`
	TrainerNames []string  `db:"-"`
}

type Filter struct {
	WorkFormatID int
	GroupID      int
	LaneID       int
	WasHeld      *bool
	From         time.Time
	To           time.Time
}

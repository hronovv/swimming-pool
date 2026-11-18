package journal

import (
	"errors"
	"time"
)

const (
	ActionLogin  = "login"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"

	ResultSuccess = "success"
	ResultError   = "error"

	defaultIP = "127.0.0.1"
)

var (
	ErrNotFound     = errors.New("journal entry not found")
	ErrInvalid      = errors.New("invalid journal entry")
	ErrUserNotFound = errors.New("user not found")
	allowedActions  = map[string]struct{}{
		ActionLogin: {}, ActionCreate: {}, ActionUpdate: {}, ActionDelete: {},
	}
	allowedResults = map[string]struct{}{
		ResultSuccess: {}, ResultError: {},
	}
)

type Entry struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Login       string    `db:"login"`
	OccurredAt  time.Time `db:"occurred_at"`
	Action      string    `db:"action"`
	Object      string    `db:"object"`
	Description string    `db:"description"`
	Result      string    `db:"result"`
	IPAddress   string    `db:"ip_address"`
}

type Filter struct {
	UserID int
	Action string
	Result string
	Query  string
	From   time.Time
	To     time.Time
}

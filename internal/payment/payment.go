package payment

import (
	"errors"
	"time"
)

const (
	MethodCash     = "cash"
	MethodCard     = "card"
	MethodTransfer = "transfer"

	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"

	formatSportPass = "sport_pass"
)

var (
	ErrNotFound        = errors.New("payment not found")
	ErrInvalid         = errors.New("invalid payment")
	ErrServiceNotFound = errors.New("service not found")
	allowedMethods     = map[string]struct{}{MethodCash: {}, MethodCard: {}, MethodTransfer: {}}
	allowedStatuses    = map[string]struct{}{StatusPending: {}, StatusCompleted: {}, StatusCancelled: {}}
)

type Payment struct {
	ID        int       `db:"id"`
	ServiceID int       `db:"service_id"`
	Amount    float64   `db:"amount"`
	PaidOn    time.Time `db:"paid_on"`
	Method    string    `db:"method"`
	Purpose   *string   `db:"purpose"`
	Status    string    `db:"status"`
}

type Filter struct {
	ServiceID int
	Status    string
	Method    string
}

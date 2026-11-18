package lane

import "errors"

const (
	StatusAvailable = "available"
	StatusClosed    = "closed"
)

var (
	ErrNotFound = errors.New("lane not found")
	ErrInvalid  = errors.New("invalid lane")
)

type Lane struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	LengthM int     `db:"length_m"`
	DepthM  float64 `db:"depth_m"`
	Status  string  `db:"status"`
}

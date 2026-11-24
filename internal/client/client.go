package client

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("client not found")
	ErrInvalid  = errors.New("invalid client")
)

type Client struct {
	ID                 int       `db:"id"`
	UserID             int       `db:"user_id"`
	Passport           string    `db:"passport"`
	BirthDate          time.Time `db:"birth_date"`
	Address            *string   `db:"address"`
	MedicalCertificate bool      `db:"medical_certificate"`
	Name               string    `db:"name"`
}

type Filter struct {
	UserID int
}

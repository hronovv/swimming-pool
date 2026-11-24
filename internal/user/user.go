package user

import (
	"errors"
	"time"
)

const (
	RoleAdmin   = "admin"
	RoleTrainer = "trainer"
	RoleClient  = "client"

	StatusActive  = "active"
	StatusBlocked = "blocked"

	minPasswordLen = 8
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrInvalid       = errors.New("invalid user")
	ErrConflict      = errors.New("user already exists")
	ErrPassportTaken = errors.New("passport already exists")
	ErrUnauthorized  = errors.New("invalid credentials")
	ErrBlocked       = errors.New("user is blocked")
	allowedRoles     = map[string]struct{}{RoleAdmin: {}, RoleTrainer: {}, RoleClient: {}}
	allowedStatuses  = map[string]struct{}{StatusActive: {}, StatusBlocked: {}}
)

type User struct {
	ID           int       `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	LastName     string    `db:"last_name"`
	FirstName    string    `db:"first_name"`
	Patronymic   *string   `db:"patronymic"`
	Role         string    `db:"role"`
	Status       string    `db:"status"`
	RegisteredAt time.Time `db:"registered_at"`
	Phone        *string   `db:"phone"`
	Email        *string   `db:"email"`
}

type Filter struct {
	Role   string
	Status string
}

type Profile struct {
	Passport           string
	BirthDate          time.Time
	Address            string
	MedicalCertificate bool
	Qualification      string
	Rank               string
	ExperienceYears    int
	Specialization     string
}

type NewUser struct {
	Login      string
	Password   string
	LastName   string
	FirstName  string
	Patronymic string
	Role       string
	Phone      string
	Email      string
	Status     string
	Profile    Profile
}

type UserUpdate struct {
	Login      string
	LastName   string
	FirstName  string
	Patronymic string
	Role       string
	Phone      string
	Email      string
	Profile    Profile
}

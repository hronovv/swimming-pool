package user

import "time"

type User struct {
	ID           int       `db:"id"`
	LastName     string    `db:"last_name"`
	FirstName    string    `db:"first_name"`
	Patronymic   *string   `db:"patronymic"`
	Role         string    `db:"role"`
	Status       string    `db:"status"`
	RegisteredAt time.Time `db:"registered_at"`
	Phone        *string   `db:"phone"`
	Email        *string   `db:"email"`
}

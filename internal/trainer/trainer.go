package trainer

import "errors"

const (
	QualificationInstructor    = "instructor"
	QualificationTrainer       = "trainer"
	QualificationSeniorTrainer = "senior_trainer"

	RankNone = "none"
	Rank3    = "3"
	Rank2    = "2"
	Rank1    = "1"
	RankCMS  = "cms"
	RankMS   = "ms"
)

var (
	ErrNotFound = errors.New("trainer not found")
	ErrInvalid  = errors.New("invalid trainer")
)

type Trainer struct {
	ID              int    `db:"id"`
	UserID          int    `db:"user_id"`
	Qualification   string `db:"qualification"`
	Rank            string `db:"rank"`
	ExperienceYears int    `db:"experience_years"`
	Specialization  string `db:"specialization"`
}

type Filter struct {
	UserID int
}

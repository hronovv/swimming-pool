package tariff

import "errors"

const (
	FormatFreeSwimOneTime        = "free_swim_one_time"
	FormatFreeSwimSubscription   = "free_swim_subscription"
	FormatIndividualOneTime      = "individual_one_time"
	FormatIndividualSubscription = "individual_subscription"
	FormatGroupOneTime           = "group_one_time"
	FormatGroupSubscription      = "group_subscription"
	FormatSportPass              = "sport_pass"

	WorkFormatFreeSwim   = 1
	WorkFormatIndividual = 2
	WorkFormatGroup      = 3
	WorkFormatSport      = 4
)

var (
	ErrNotFound        = errors.New("tariff not found")
	ErrInvalid         = errors.New("invalid tariff")
	ErrConflict        = errors.New("tariff already exists")
	ErrInUse           = errors.New("tariff is still in use")
	workFormatByFormat = map[string]int{
		FormatFreeSwimOneTime:        WorkFormatFreeSwim,
		FormatFreeSwimSubscription:   WorkFormatFreeSwim,
		FormatIndividualOneTime:      WorkFormatIndividual,
		FormatIndividualSubscription: WorkFormatIndividual,
		FormatGroupOneTime:           WorkFormatGroup,
		FormatGroupSubscription:      WorkFormatGroup,
		FormatSportPass:              WorkFormatSport,
	}
)

type Tariff struct {
	ID                    int     `db:"id"`
	Name                  string  `db:"name"`
	Format                string  `db:"format"`
	WorkFormatID          int     `db:"work_format_id"`
	Price                 float64 `db:"price"`
	VisitLimit            *int    `db:"visit_limit"`
	DurationDays          int     `db:"duration_days"`
	ServicesCnt           int64   `db:"services_cnt"`
	CompletedPaymentsCnt  int64   `db:"completed_payments_cnt"`
	Revenue               float64 `db:"revenue"`
	AvgCheck              float64 `db:"avg_check"`
}

package tariff

type Tariff struct {
	Name         string `db:"name"`
	Format       string `db:"format"`
	Price        int    `db:"price"`
	VisitLimit   *int    `db:"visit_limit"`
	DurationDays int    `db:"duration_days"`
}

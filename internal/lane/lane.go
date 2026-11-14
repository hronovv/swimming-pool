package lane

type Lane struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	LengthM int     `db:"length_m"`
	DepthM  float64 `db:"depth_m"`
	Status  string  `db:"status"`
}

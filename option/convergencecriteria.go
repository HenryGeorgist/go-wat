package option

//this should probably live in go-statistics with inline histogram.
type ConvergenceCriteria struct {
	MinIterations int64   `json:"min_iterations"`
	MaxIterations int64   `json:"max_iterations"`
	ZAlpha        float64 `json:"z_alpha"`
	Tolerance     float64 `json:"tolerance"`
}

package types

// MetricsEnabled represents whether metrics collection is enabled for the server.
type MetricsEnabled struct {
	Enabled bool `json:"enabled"` // Enabled indicates if metrics are enabled (true) or disabled (false).
}

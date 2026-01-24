package types

// Limit represents a data transfer limit for an access key.
// The zero value indicates no limit.
type Limit struct {
	Bytes uint64 `json:"bytes"` // Bytes is the maximum number of bytes allowed for data transfer. A value of 0 means no limit is enforced.
}

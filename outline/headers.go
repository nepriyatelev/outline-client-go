package outline

// Headers represents a map of HTTP headers.
type Headers map[string]string

// DefaultHeaders returns the default HTTP headers used for API requests.
func DefaultHeaders() Headers {
	return Headers{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}

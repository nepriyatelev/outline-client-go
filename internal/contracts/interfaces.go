package contracts

import (
	"context"
)

// Request represents an HTTP request structure.
type Request struct {
	Method  string            // Method is the HTTP method (e.g., GET, POST).
	URL     string            // URL is the request URL.
	Headers map[string]string // Headers is a map of HTTP headers.
	Body    []byte            // Body is the request body.
}

// Response represents an HTTP response structure.
type Response struct {
	StatusCode int               // StatusCode is the HTTP status code.
	Headers    map[string]string // Headers is a map of HTTP response headers.
	Body       []byte            // Body is the response body.
}

// Doer defines an interface for executing HTTP requests.
type Doer interface {
	Do(ctx context.Context, req *Request) (*Response, error)
}

// Logger defines an interface for logging messages.
type Logger interface {
	// Debugf logs debug messages with formatting.
	Debugf(ctx context.Context, format string, args ...any)
	// Infof logs informational messages with formatting.
	Infof(ctx context.Context, format string, args ...any)
}

package outline

import (
	"fmt"
)

// ClientError is the root error type for the Outline client.
type ClientError struct {
	Code    int
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("outline client error [%d]: %s", e.Code, e.Message)
}

// Predefined errors for specific statuses.
var (
	errUnexpected = func(statusCode int) *ClientError {
		return &ClientError{
			Code:    statusCode,
			Message: "An unexpected error occurred.",
		}
	}
)
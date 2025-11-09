package outline

import (
	"fmt"
)

type ParseURLError struct {
	BaseURL string
	Err     error
}

func (e *ParseURLError) Error() string {
	return fmt.Sprintf("outline client error: invalid baseURL %q: %v", e.BaseURL, e.Err)
}

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
	errUnexpected = func(statusCode int, body []byte) *ClientError {
		return &ClientError{
			Code:    statusCode,
			Message: fmt.Sprintf("An unexpected error occurred: body=%s", string(body)),
		}
	}
	
	errParseBaseURL = func(baseURL string, err error) *ParseURLError {
		return &ParseURLError{
			BaseURL: baseURL,
			Err:     err,
		}
	}
)

// UnmarshalError содержит детали ошибки при распаковке JSON
type UnmarshalError struct {
	Data []byte
	Type string
	Err  error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal %s failed: %v", e.Type, e.Err)
}

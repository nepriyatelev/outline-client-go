package outline

import (
	"errors"
	"fmt"
)

const (
	clientOutlineErrStr        = "outline client error"
	invalidBaseURLErrStr       = "invalid baseURL"
	unmarshalFailedErrStr      = "unmarshal failed"
	unmarshalEmptyBodyErrStr   = "empty body"
	invalidHostnameErrStr      = "invalid hostname or IP address"
	internalHostNameErrStr     = "internal error occurred while validating hostname or IP address"
	invalidPortErrStr          = "requested port wasn't integer from 1 through 65535, or request had no port parameter"
	portAlreadyInUseErrStr     = "requested port was already in use by another service"
	invalidServerNameErrStr    = "invalid server name"
	invalidRequestErrStr       = "invalid request"
	invalidDataLimitErrStr     = "invalid data limit"
	accessKeyNotFoundErrStr    = "access key not found"
	unexpectedStatusCodeErrStr = "unexpected status code"
	doOperationErrStr          = "do operation error"
)

var (
	// ClientOutlineError is the base error for all client errors.
	ClientOutlineError = errors.New(clientOutlineErrStr)

	// InvalidBaseURLError indicates that the provided base URL is malformed or empty.
	InvalidBaseURLError = errors.New(invalidBaseURLErrStr)

	// UnmarshalFailedError indicates that JSON unmarshaling failed.
	UnmarshalFailedError = errors.New(unmarshalFailedErrStr)

	// UnmarshalEmptyBodyError indicates that the response body was empty when data was expected.
	UnmarshalEmptyBodyError = errors.New(unmarshalEmptyBodyErrStr)

	// InvalidHostnameError indicates that the provided hostname or IP address is invalid.
	InvalidHostnameError = errors.New(invalidHostnameErrStr)

	// InternalHostNameError indicates an internal server error during hostname validation.
	InternalHostNameError = errors.New(internalHostNameErrStr)

	// InvalidPortError indicates that the port number is outside the valid range 1-65535.
	InvalidPortError = errors.New(invalidPortErrStr)

	// PortAlreadyInUseError indicates that the requested port is already in use by another service.
	PortAlreadyInUseError = errors.New(portAlreadyInUseErrStr)

	// InvalidServerNameError indicates that the provided server name is invalid.
	InvalidServerNameError = errors.New(invalidServerNameErrStr)

	// InvalidRequestError indicates that the request parameters are invalid.
	InvalidRequestError = errors.New(invalidRequestErrStr)

	// InvalidDataLimitError indicates that the provided data limit value is invalid.
	InvalidDataLimitError = errors.New(invalidDataLimitErrStr)

	// AccessKeyNotFoundError indicates that the requested access key does not exist.
	AccessKeyNotFoundError = errors.New(accessKeyNotFoundErrStr)

	// UnexpectedStatusCodeError indicates that the server returned an unexpected HTTP status code.
	UnexpectedStatusCodeError = errors.New(unexpectedStatusCodeErrStr)

	// DoOperationError indicates that the HTTP request execution failed.
	DoOperationError = errors.New(doOperationErrStr)
)

// ClientError represents an error returned by the Outline server API.
// It contains the HTTP status code, response body, and a descriptive message.
type ClientError struct {
	statusCode int
	data       []byte
	message    string
	err        error
}

// Error returns a formatted error message including status code and response data.
func (e *ClientError) Error() string {
	msg := fmt.Sprintf("%s; status code: %d", e.message, e.statusCode)
	if len(e.data) > 0 {
		msg = fmt.Sprintf("%s; data: %s", msg, e.data)
	}
	return withLastError(msg, e.err)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *ClientError) Unwrap() error {
	return e.err
}

var (
	errInvalidHostname = func(statusCode int, hostnameOrIP string) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (host name or ip: %s)",
				ClientOutlineError.Error(),
				hostnameOrIP,
			),
			err: errors.Join(ClientOutlineError, InvalidHostnameError),
		}
	}
	errInternalHostname = func(statusCode int, hostnameOrIP string) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (host name or ip: %s)",
				ClientOutlineError.Error(),
				hostnameOrIP,
			),
			err: errors.Join(ClientOutlineError, InternalHostNameError),
		}
	}
	errInvalidPort = func(statusCode int, port uint16) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (port: %d)",
				ClientOutlineError.Error(),
				port,
			),
			err: errors.Join(ClientOutlineError, InvalidPortError),
		}
	}
	errPortAlreadyInUse = func(statusCode int, port uint16) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (port: %d)",
				ClientOutlineError.Error(),
				port,
			),
			err: errors.Join(ClientOutlineError, PortAlreadyInUseError),
		}
	}
	errInvalidServerName = func(statusCode int, name string) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (server name: %s)",
				ClientOutlineError.Error(),
				name,
			),
			err: errors.Join(ClientOutlineError, InvalidServerNameError),
		}
	}
	errInvalidRequest = func(statusCode int, body string) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (response body: %s)",
				ClientOutlineError.Error(),
				body,
			),
			err: errors.Join(ClientOutlineError, InvalidRequestError),
		}
	}
	errInvalidDataLimit = func(statusCode int, bytes uint64) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (data limit bytes: %d)",
				ClientOutlineError.Error(),
				bytes,
			),
			err: errors.Join(ClientOutlineError, InvalidDataLimitError),
		}
	}
	errAccessKeyNotFound = func(statusCode int, accessKeyID string) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			message: fmt.Sprintf("%s: (access key id: %s)",
				ClientOutlineError.Error(),
				accessKeyID,
			),
			err: errors.Join(ClientOutlineError, AccessKeyNotFoundError),
		}
	}
	errUnexpectedStatusCode = func(statusCode int, data []byte) *ClientError {
		return &ClientError{
			statusCode: statusCode,
			data:       data,
			message:    fmt.Sprintf("%s: %s", ClientOutlineError.Error(), UnexpectedStatusCodeError.Error()),
			err:        errors.Join(ClientOutlineError, UnexpectedStatusCodeError),
		}
	}
)

// ParseURLError represents an error that occurs when parsing the base URL.
// It wraps [InvalidBaseURLError] and contains the original URL that failed to parse.
type ParseURLError struct {
	baseURL string
	message string
	err     error
}

// Error returns a formatted error message including the problematic URL.
func (e *ParseURLError) Error() string {
	var msg string
	if e.baseURL == "" {
		msg = fmt.Sprintf("%s; baseUrl is empty", e.message)
	} else {
		msg = fmt.Sprintf("%s; (base url: %s)", e.message, e.baseURL)
	}
	return withLastError(msg, e.err)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *ParseURLError) Unwrap() error {
	return e.err
}

var errParseBaseURL = func(baseURL string, err error) *ParseURLError {
	return &ParseURLError{
		baseURL: baseURL,
		message: fmt.Sprintf("%s: %s", ClientOutlineError.Error(), InvalidBaseURLError.Error()),
		err:     errors.Join(ClientOutlineError, InvalidBaseURLError, err),
	}
}

// UnmarshalError represents an error that occurs when unmarshaling JSON response data.
// It wraps [UnmarshalFailedError] and contains the raw data that failed to unmarshal.
type UnmarshalError struct {
	data    []byte
	typeStr string
	message string
	err     error
}

// Error returns a formatted error message including the target type and raw data.
func (e *UnmarshalError) Error() string {
	msg := e.message
	if e.typeStr != "" {
		msg = fmt.Sprintf("%s; (type: %s)", msg, e.typeStr)
	}
	if len(e.data) > 0 {
		msg = fmt.Sprintf("%s; data: %s", msg, string(e.data))
	}
	return withLastError(msg, e.err)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *UnmarshalError) Unwrap() error {
	return e.err
}

var (
	errUnmarshal = func(data []byte, typeStr string, err error) *UnmarshalError {
		return &UnmarshalError{
			data:    data,
			typeStr: typeStr,
			message: fmt.Sprintf("%s: %s", ClientOutlineError.Error(), UnmarshalFailedError.Error()),
			err:     errors.Join(ClientOutlineError, UnmarshalFailedError, err),
		}
	}

	errUnmarshalEmptyBody = func(typeStr string) *UnmarshalError {
		return &UnmarshalError{
			data:    []byte{},
			typeStr: typeStr,
			message: fmt.Sprintf("%s: %s", ClientOutlineError.Error(), UnmarshalFailedError.Error()),
			err:     errors.Join(ClientOutlineError, UnmarshalFailedError, UnmarshalEmptyBodyError),
		}
	}
)

// DoError represents an error that occurs when executing an HTTP request.
// It wraps [DoOperationError] and contains the operation name that failed.
type DoError struct {
	operation string
	message   string
	err       error
}

// Error returns a formatted error message including the operation name.
func (e *DoError) Error() string {
	msg := fmt.Sprintf("%s; operation: %s", e.message, e.operation)
	return withLastError(msg, e.err)
}

// Unwrap returns the underlying error for use with [errors.Is] and [errors.As].
func (e *DoError) Unwrap() error {
	return e.err
}

var (
	errDoGetServerInfo = func(err error) *DoError {
		return &DoError{
			operation: "get server info",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateServerHostname = func(err error) *DoError {
		return &DoError{
			operation: "update server hostname",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdatePortNewAccessKeys = func(err error) *DoError {
		return &DoError{
			operation: "update port for new access keys",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateServerName = func(err error) *DoError {
		return &DoError{
			operation: "update server name",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoGetMetricsEnabled = func(err error) *DoError {
		return &DoError{
			operation: "get metrics enabled",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateMetricsEnabled = func(err error) *DoError {
		return &DoError{
			operation: "update metrics enabled",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateKeyLimitBytes = func(err error) *DoError {
		return &DoError{
			operation: "update key limit bytes",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoDeleteKeyLimitBytes = func(err error) *DoError {
		return &DoError{
			operation: "delete key limit bytes",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoCreateAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "create access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoGetAccessKeys = func(err error) *DoError {
		return &DoError{
			operation: "get access keys",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoGetAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "get access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "update access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoDeleteAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "delete access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateNameAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "update name access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoUpdateDataLimitAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "update data limit access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoDeleteDataLimitAccessKey = func(err error) *DoError {
		return &DoError{
			operation: "delete data limit access key",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoGetMetricsTransfer = func(err error) *DoError {
		return &DoError{
			operation: "get metrics transfer",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
	errDoGetExperimentalMetrics = func(err error) *DoError {
		return &DoError{
			operation: "get experimental metrics",
			message:   fmt.Sprintf("%s: %s", ClientOutlineError.Error(), DoOperationError.Error()),
			err:       errors.Join(ClientOutlineError, DoOperationError, err),
		}
	}
)

func withLastError(message string, err error) string {
	var lastErr error
	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		errs := uw.Unwrap()
		if len(errs) > 0 {
			lastErr = errs[len(errs)-1]
		}
	}
	if lastErr != nil {
		message = fmt.Sprintf("%s; reason: %s", message, lastErr.Error())
	}

	return message + "."
}

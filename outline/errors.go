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
	ClientOutlineError        = errors.New(clientOutlineErrStr)
	InvalidBaseURLError       = errors.New(invalidBaseURLErrStr)
	UnmarshalFailedError      = errors.New(unmarshalFailedErrStr)
	UnmarshalEmptyBodyError   = errors.New(unmarshalEmptyBodyErrStr)
	InvalidHostnameError      = errors.New(invalidHostnameErrStr)
	InternalHostNameError     = errors.New(internalHostNameErrStr)
	InvalidPortError          = errors.New(invalidPortErrStr)
	PortAlreadyInUseError     = errors.New(portAlreadyInUseErrStr)
	InvalidServerNameError    = errors.New(invalidServerNameErrStr)
	InvalidRequestError       = errors.New(invalidRequestErrStr)
	InvalidDataLimitError     = errors.New(invalidDataLimitErrStr)
	AccessKeyNotFoundError    = errors.New(accessKeyNotFoundErrStr)
	UnexpectedStatusCodeError = errors.New(unexpectedStatusCodeErrStr)
	DoOperationError          = errors.New(doOperationErrStr)
)

type ClientError struct {
	statusCode int
	data       []byte
	message    string
	err        error
}

func (e *ClientError) Error() string {
	e.message = fmt.Sprintf("%s; status code: %d", e.message, e.statusCode)
	if len(e.data) > 0 {
		e.message = fmt.Sprintf("%s; data: %s", e.message, e.data)
	}
	return withLastError(e.message, e.err)
}

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

type ParseURLError struct {
	baseURL string
	message string
	err     error
}

func (e *ParseURLError) Error() string {
	if e.baseURL == "" {
		e.message = fmt.Sprintf("%s; baseUrl is empty", e.message)
	} else {
		e.message = fmt.Sprintf("%s; (base url: %s)", e.message, e.baseURL)
	}
	return withLastError(e.message, e.err)
}

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

type UnmarshalError struct {
	data    []byte
	typeStr string
	message string
	err     error
}

func (e *UnmarshalError) Error() string {
	if e.typeStr != "" {
		e.message = fmt.Sprintf("%s; (type: %s)", e.message, e.typeStr)
	}
	if len(e.data) > 0 {
		e.message = fmt.Sprintf("%s; data: %s", e.message, string(e.data))
	}
	return withLastError(e.message, e.err)
}

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

type DoError struct {
	operation string
	message   string
	err       error
}

func (e *DoError) Error() string {
	e.message = fmt.Sprintf("%s; operation: %s", e.message, e.operation)
	return withLastError(e.message, e.err)
}

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

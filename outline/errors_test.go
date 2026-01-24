package outline

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInvalidHostname(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		hostnameOrIP string
		expectedMsg  string
	}{
		{
			name:         "valid hostname",
			statusCode:   400,
			hostnameOrIP: "invalid.host",
			expectedMsg:  "outline client error: (host name or ip: invalid.host); status code: 400; reason: invalid hostname or IP address.",
		},
		{
			name:         "empty hostname",
			statusCode:   404,
			hostnameOrIP: "",
			expectedMsg:  "outline client error: (host name or ip: ); status code: 404; reason: invalid hostname or IP address.",
		},
		{
			name:         "IP address",
			statusCode:   500,
			hostnameOrIP: "192.168.1.1",
			expectedMsg:  "outline client error: (host name or ip: 192.168.1.1); status code: 500; reason: invalid hostname or IP address.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errInvalidHostname(tt.statusCode, tt.hostnameOrIP)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidHostnameError)
		})
	}
}

func TestErrInternalHostname(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		hostnameOrIP string
		expectedMsg  string
	}{
		{
			name:         "valid hostname",
			statusCode:   400,
			hostnameOrIP: "invalid.host",
			expectedMsg:  "outline client error: (host name or ip: invalid.host); status code: 400; reason: internal error occurred while validating hostname or IP address.",
		},
		{
			name:         "empty hostname",
			statusCode:   404,
			hostnameOrIP: "",
			expectedMsg:  "outline client error: (host name or ip: ); status code: 404; reason: internal error occurred while validating hostname or IP address.",
		},
		{
			name:         "IP address",
			statusCode:   500,
			hostnameOrIP: "192.168.1.1",
			expectedMsg:  "outline client error: (host name or ip: 192.168.1.1); status code: 500; reason: internal error occurred while validating hostname or IP address.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errInternalHostname(tt.statusCode, tt.hostnameOrIP)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InternalHostNameError)
		})
	}
}

func TestErrInvalidPort(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		port        uint16
		expectedMsg string
	}{
		{
			name:        "port 0",
			statusCode:  400,
			port:        0,
			expectedMsg: "outline client error: (port: 0); status code: 400; reason: requested port wasn't integer from 1 through 65535, or request had no port parameter.",
		},
		{
			name:        "port 80",
			statusCode:  404,
			port:        80,
			expectedMsg: "outline client error: (port: 80); status code: 404; reason: requested port wasn't integer from 1 through 65535, or request had no port parameter.",
		},
		{
			name:        "port 65535",
			statusCode:  500,
			port:        65535,
			expectedMsg: "outline client error: (port: 65535); status code: 500; reason: requested port wasn't integer from 1 through 65535, or request had no port parameter.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errInvalidPort(tt.statusCode, tt.port)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidPortError)
		})
	}
}

func TestErrPortAlreadyInUse(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		port        uint16
		expectedMsg string
	}{
		{
			name:        "port 8080",
			statusCode:  409,
			port:        8080,
			expectedMsg: "outline client error: (port: 8080); status code: 409; reason: requested port was already in use by another service.",
		},
		{
			name:        "port 443",
			statusCode:  500,
			port:        443,
			expectedMsg: "outline client error: (port: 443); status code: 500; reason: requested port was already in use by another service.",
		},
		{
			name:        "port 22",
			statusCode:  400,
			port:        22,
			expectedMsg: "outline client error: (port: 22); status code: 400; reason: requested port was already in use by another service.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errPortAlreadyInUse(tt.statusCode, tt.port)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, PortAlreadyInUseError)
		})
	}
}

func TestErrInvalidServerName(t *testing.T) {
	tests := []struct {
		testName    string
		statusCode  int
		serverName  string
		expectedMsg string
	}{
		{
			testName:    "empty server name",
			statusCode:  400,
			serverName:  "",
			expectedMsg: "outline client error: (server name: ); status code: 400; reason: invalid server name.",
		},
		{
			testName:    "valid server name",
			statusCode:  404,
			serverName:  "MyServer",
			expectedMsg: "outline client error: (server name: MyServer); status code: 404; reason: invalid server name.",
		},
		{
			testName:    "special characters",
			statusCode:  500,
			serverName:  "Server@123",
			expectedMsg: "outline client error: (server name: Server@123); status code: 500; reason: invalid server name.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := errInvalidServerName(tt.statusCode, tt.serverName)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidServerNameError)
		})
	}
}

func TestErrInvalidRequest(t *testing.T) {
	tests := []struct {
		testName    string
		statusCode  int
		body        string
		expectedMsg string
	}{
		{
			testName:    "empty body",
			statusCode:  400,
			body:        "",
			expectedMsg: "outline client error: (response body: ); status code: 400; reason: invalid request.",
		},
		{
			testName:    "json body",
			statusCode:  422,
			body:        `{"error": "invalid input"}`,
			expectedMsg: "outline client error: (response body: {\"error\": \"invalid input\"}); status code: 422; reason: invalid request.",
		},
		{
			testName:    "text body",
			statusCode:  500,
			body:        "Internal server error",
			expectedMsg: "outline client error: (response body: Internal server error); status code: 500; reason: invalid request.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := errInvalidRequest(tt.statusCode, tt.body)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidRequestError)
		})
	}
}

func TestErrInvalidDataLimit(t *testing.T) {
	tests := []struct {
		testName    string
		statusCode  int
		bytes       uint64
		expectedMsg string
	}{
		{
			testName:    "zero bytes",
			statusCode:  400,
			bytes:       0,
			expectedMsg: "outline client error: (data limit bytes: 0); status code: 400; reason: invalid data limit.",
		},
		{
			testName:    "small limit",
			statusCode:  422,
			bytes:       1024,
			expectedMsg: "outline client error: (data limit bytes: 1024); status code: 422; reason: invalid data limit.",
		},
		{
			testName:    "large limit",
			statusCode:  500,
			bytes:       1000000000,
			expectedMsg: "outline client error: (data limit bytes: 1000000000); status code: 500; reason: invalid data limit.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := errInvalidDataLimit(tt.statusCode, tt.bytes)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidDataLimitError)
		})
	}
}

func TestErrAccessKeyNotFound(t *testing.T) {
	tests := []struct {
		testName    string
		statusCode  int
		accessKeyID string
		expectedMsg string
	}{
		{
			testName:    "empty access key ID",
			statusCode:  404,
			accessKeyID: "",
			expectedMsg: "outline client error: (access key id: ); status code: 404; reason: access key not found.",
		},
		{
			testName:    "valid access key ID",
			statusCode:  500,
			accessKeyID: "abc123",
			expectedMsg: "outline client error: (access key id: abc123); status code: 500; reason: access key not found.",
		},
		{
			testName:    "special characters",
			statusCode:  400,
			accessKeyID: "key@456",
			expectedMsg: "outline client error: (access key id: key@456); status code: 400; reason: access key not found.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := errAccessKeyNotFound(tt.statusCode, tt.accessKeyID)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, AccessKeyNotFoundError)
		})
	}
}

func TestErrUnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		testName    string
		statusCode  int
		data        []byte
		expectedMsg string
	}{
		{
			testName:    "empty data",
			statusCode:  500,
			data:        []byte{},
			expectedMsg: "outline client error: unexpected status code; status code: 500; reason: unexpected status code.",
		},
		{
			testName:    "with data",
			statusCode:  404,
			data:        []byte("Not Found"),
			expectedMsg: "outline client error: unexpected status code; status code: 404; data: Not Found; reason: unexpected status code.",
		},
		{
			testName:    "json data",
			statusCode:  422,
			data:        []byte(`{"error": "unprocessable"}`),
			expectedMsg: "outline client error: unexpected status code; status code: 422; data: {\"error\": \"unprocessable\"}; reason: unexpected status code.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := errUnexpectedStatusCode(tt.statusCode, tt.data)

			// Check type
			assert.IsType(t, &ClientError{}, err)

			// Check error can be assigned to ClientError
			var ce *ClientError
			assert.ErrorAs(t, err, &ce)

			// Check status code
			assert.Equal(t, tt.statusCode, err.statusCode)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying error
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, UnexpectedStatusCodeError)
		})
	}
}

func TestErrParseBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "valid baseURL with error",
			baseURL:     "http://example.com",
			inputErr:    errors.New("parse error"),
			expectedMsg: "outline client error: invalid baseURL; (base url: http://example.com); reason: parse error.",
		},
		{
			name:        "empty baseURL with error",
			baseURL:     "",
			inputErr:    errors.New("parse error"),
			expectedMsg: "outline client error: invalid baseURL; baseUrl is empty; reason: parse error.",
		},
		{
			name:        "valid baseURL with nil error",
			baseURL:     "https://valid.com",
			inputErr:    nil,
			expectedMsg: "outline client error: invalid baseURL; (base url: https://valid.com); reason: invalid baseURL.",
		},
		{
			name:        "empty baseURL with nil error",
			baseURL:     "",
			inputErr:    nil,
			expectedMsg: "outline client error: invalid baseURL; baseUrl is empty; reason: invalid baseURL.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errParseBaseURL(tt.baseURL, tt.inputErr)

			// Check type
			assert.IsType(t, &ParseURLError{}, err)

			// Check baseURL
			assert.Equal(t, tt.baseURL, err.baseURL)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, InvalidBaseURLError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		typeStr     string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with data, typeStr, and error",
			data:        []byte(`{"key": "value"}`),
			typeStr:     "Server",
			inputErr:    errors.New("json unmarshal error"),
			expectedMsg: "outline client error: unmarshal failed; (type: Server); data: {\"key\": \"value\"}; reason: json unmarshal error.",
		},
		{
			name:        "with data, empty typeStr, and error",
			data:        []byte("some data"),
			typeStr:     "",
			inputErr:    errors.New("parse error"),
			expectedMsg: "outline client error: unmarshal failed; data: some data; reason: parse error.",
		},
		{
			name:        "empty data, with typeStr, and error",
			data:        []byte{},
			typeStr:     "AccessKey",
			inputErr:    errors.New("invalid syntax"),
			expectedMsg: "outline client error: unmarshal failed; (type: AccessKey); reason: invalid syntax.",
		},
		{
			name:        "empty data, empty typeStr, and error",
			data:        []byte{},
			typeStr:     "",
			inputErr:    errors.New("decode error"),
			expectedMsg: "outline client error: unmarshal failed; reason: decode error.",
		},
		{
			name:        "with data, typeStr, nil error",
			data:        []byte(`[1,2,3]`),
			typeStr:     "Metrics",
			inputErr:    nil,
			expectedMsg: "outline client error: unmarshal failed; (type: Metrics); data: [1,2,3]; reason: unmarshal failed.",
		},
		{
			name:        "empty data, empty typeStr, nil error",
			data:        []byte{},
			typeStr:     "",
			inputErr:    nil,
			expectedMsg: "outline client error: unmarshal failed; reason: unmarshal failed.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errUnmarshal(tt.data, tt.typeStr, tt.inputErr)

			// Check type
			assert.IsType(t, &UnmarshalError{}, err)

			// Check data
			assert.Equal(t, tt.data, err.data)

			// Check typeStr
			assert.Equal(t, tt.typeStr, err.typeStr)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, UnmarshalFailedError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrUnmarshalEmptyBody(t *testing.T) {
	tests := []struct {
		name        string
		typeStr     string
		expectedMsg string
	}{
		{
			name:        "with typeStr",
			typeStr:     "Server",
			expectedMsg: "outline client error: unmarshal failed; (type: Server); reason: empty body.",
		},
		{
			name:        "empty typeStr",
			typeStr:     "",
			expectedMsg: "outline client error: unmarshal failed; reason: empty body.",
		},
		{
			name:        "different typeStr",
			typeStr:     "AccessKey",
			expectedMsg: "outline client error: unmarshal failed; (type: AccessKey); reason: empty body.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errUnmarshalEmptyBody(tt.typeStr)

			// Check type
			assert.IsType(t, &UnmarshalError{}, err)

			// Check data is empty
			assert.Equal(t, []byte{}, err.data)

			// Check typeStr
			assert.Equal(t, tt.typeStr, err.typeStr)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, UnmarshalFailedError)
			assert.ErrorIs(t, err.err, UnmarshalEmptyBodyError)
		})
	}
}

func TestErrDoGetServerInfo(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get server info; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get server info; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetServerInfo(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get server info", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateServerHostname(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update server hostname; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update server hostname; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateServerHostname(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update server hostname", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateServerName(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update server name; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update server name; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateServerName(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update server name", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdatePortNewAccessKeys(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update port for new access keys; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update port for new access keys; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdatePortNewAccessKeys(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update port for new access keys", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoGetMetricsEnabled(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get metrics enabled; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get metrics enabled; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetMetricsEnabled(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get metrics enabled", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateMetricsEnabled(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update metrics enabled; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update metrics enabled; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateMetricsEnabled(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update metrics enabled", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateKeyLimitBytes(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update key limit bytes; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update key limit bytes; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateKeyLimitBytes(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update key limit bytes", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoDeleteKeyLimitBytes(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: delete key limit bytes; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: delete key limit bytes; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoDeleteKeyLimitBytes(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "delete key limit bytes", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoCreateAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: create access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: create access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoCreateAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "create access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoGetAccessKeys(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get access keys; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get access keys; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetAccessKeys(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get access keys", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoGetAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoDeleteAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: delete access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: delete access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoDeleteAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "delete access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateNameAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update name access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update name access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateNameAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update name access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoUpdateDataLimitAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: update data limit access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: update data limit access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoUpdateDataLimitAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "update data limit access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoDeleteDataLimitAccessKey(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: delete data limit access key; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: delete data limit access key; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoDeleteDataLimitAccessKey(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "delete data limit access key", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoGetMetricsTransfer(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get metrics transfer; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get metrics transfer; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetMetricsTransfer(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get metrics transfer", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestErrDoGetExperimentalMetrics(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectedMsg string
	}{
		{
			name:        "with error",
			inputErr:    errors.New("network error"),
			expectedMsg: "outline client error: do operation error; operation: get experimental metrics; reason: network error.",
		},
		{
			name:        "with nil error",
			inputErr:    nil,
			expectedMsg: "outline client error: do operation error; operation: get experimental metrics; reason: do operation error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errDoGetExperimentalMetrics(tt.inputErr)

			// Check type
			assert.IsType(t, &DoError{}, err)

			// Check operation
			assert.Equal(t, "get experimental metrics", err.operation)

			// Check error message
			assert.EqualError(t, err, tt.expectedMsg)

			// Check underlying errors
			assert.ErrorIs(t, err.err, ClientOutlineError)
			assert.ErrorIs(t, err.err, DoOperationError)
			if tt.inputErr != nil {
				assert.ErrorIs(t, err.err, tt.inputErr)
			}
		})
	}
}

func TestWithLastError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		inputErr error
		expected string
	}{
		{
			name:     "nil error",
			message:  "test message",
			inputErr: nil,
			expected: "test message.",
		},
		{
			name:     "simple error without Unwrap",
			message:  "test message",
			inputErr: errors.New("simple error"),
			expected: "test message.",
		},
		{
			name:     "joined error with one error",
			message:  "test message",
			inputErr: errors.Join(errors.New("first error")),
			expected: "test message; reason: first error.",
		},
		{
			name:     "joined error with two errors",
			message:  "test message",
			inputErr: errors.Join(errors.New("first error"), errors.New("second error")),
			expected: "test message; reason: second error.",
		},
		{
			name:     "joined error with three errors",
			message:  "test message",
			inputErr: errors.Join(errors.New("first error"), errors.New("second error"), errors.New("third error")),
			expected: "test message; reason: third error.",
		},
		{
			name:     "empty message with error",
			message:  "",
			inputErr: errors.Join(errors.New("some error")),
			expected: "; reason: some error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := withLastError(tt.message, tt.inputErr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

package outline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// createTestClientForAccessKeys creates a Client with mock doer for testing access key operations.
func createTestClientForAccessKeys(doer contracts.Doer) *Client {
	baseURL, _ := url.Parse("http://localhost:8081/api/")
	return MustNewClient(baseURL.String(), "", WithClient(doer))
}

// newMockDoerAccessKey configures mock to return provided response/error and capture request.
func newMockDoerAccessKey(
	t *testing.T,
	resp *contracts.Response,
	err error,
	capture **contracts.Request,
) *MockDoer {
	m := NewMockDoer(t)
	m.On("Do", mock.Anything, mock.AnythingOfType("*contracts.Request")).
		Run(func(args mock.Arguments) {
			if capture != nil {
				if req, ok := args[1].(*contracts.Request); ok {
					*capture = req
				}
			}
		}).
		Return(resp, err)
	return m
}

// === CreateAccessKey Tests ===

func TestCreateAccessKey_Success(t *testing.T) {
	tests := []struct {
		name            string
		createAccessKey *types.CreateAccessKey
		expectedKey     types.AccessKey
	}{
		{
			name: "with all fields",
			createAccessKey: &types.CreateAccessKey{
				Method:   "aes-192-gcm",
				Name:     "Test Key",
				Password: "securepassword123",
				Port:     8388,
				Limit:    &types.Limit{Bytes: 10000},
			},
			expectedKey: types.AccessKey{
				ID:        "key-123",
				Name:      "Test Key",
				Password:  "securepassword123",
				Port:      8388,
				Method:    "aes-192-gcm",
				AccessURL: "ss://test@example.com:8388",
			},
		},
		{
			name: "with required fields only",
			createAccessKey: &types.CreateAccessKey{
				Method: "aes-256-gcm",
			},
			expectedKey: types.AccessKey{
				ID:        "key-456",
				Name:      "",
				Password:  "generated-password",
				Port:      8080,
				Method:    "aes-256-gcm",
				AccessURL: "ss://generated@example.com:8080",
			},
		},
		{
			name:            "with nil createAccessKey",
			createAccessKey: nil,
			expectedKey: types.AccessKey{
				ID:        "key-789",
				Name:      "",
				Password:  "default-password",
				Port:      8080,
				Method:    "chacha20-ietf-poly1305",
				AccessURL: "ss://default@example.com:8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			respBody, _ := json.Marshal(tt.expectedKey)
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusCreated,
				Body:       respBody,
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.CreateAccessKey(ctx, tt.createAccessKey)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedKey.ID, result.ID)
			assert.Equal(t, tt.expectedKey.Name, result.Name)
			assert.Equal(t, tt.expectedKey.Password, result.Password)
			assert.Equal(t, tt.expectedKey.Port, result.Port)
			assert.Equal(t, tt.expectedKey.Method, result.Method)
			assert.Equal(t, tt.expectedKey.AccessURL, result.AccessURL)
			assert.Equal(t, http.MethodPost, req.Method)
		})
	}
}

func TestCreateAccessKey_RequestBody(t *testing.T) {
	// Arrange
	createAccessKey := &types.CreateAccessKey{
		Method:   "aes-192-gcm",
		Name:     "My Access Key",
		Password: "mypassword",
		Port:     9000,
		Limit:    &types.Limit{Bytes: 50000},
	}

	expectedKey := types.AccessKey{
		ID:        "key-test",
		Name:      "My Access Key",
		Password:  "mypassword",
		Port:      9000,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:9000",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	_, err := client.CreateAccessKey(ctx, createAccessKey)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)

	var sentBody types.CreateAccessKey
	err = json.Unmarshal(capturedReq.Body, &sentBody)
	require.NoError(t, err)

	assert.Equal(t, createAccessKey.Method, sentBody.Method)
	assert.Equal(t, createAccessKey.Name, sentBody.Name)
	assert.Equal(t, createAccessKey.Password, sentBody.Password)
	assert.Equal(t, createAccessKey.Port, sentBody.Port)
	require.NotNil(t, sentBody.Limit)
	assert.Equal(t, createAccessKey.Limit.Bytes, sentBody.Limit.Bytes)
}

func TestCreateAccessKey_NilRequestBody(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "key-nil",
		Name:      "Default",
		Password:  "generated",
		Port:      8080,
		Method:    "chacha20-ietf-poly1305",
		AccessURL: "ss://gen@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.CreateAccessKey(ctx, nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, capturedReq)
	assert.Empty(t, capturedReq.Body)
}

func TestCreateAccessKey_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	createAccessKey := &types.CreateAccessKey{
		Method: "aes-192-gcm",
	}

	// Act
	result, err := client.CreateAccessKey(ctx, createAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestCreateAccessKey_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       []byte("invalid json"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	createAccessKey := &types.CreateAccessKey{
		Method: "aes-192-gcm",
	}

	// Act
	result, err := client.CreateAccessKey(ctx, createAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestCreateAccessKey_EmptyBody(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       []byte{},
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	createAccessKey := &types.CreateAccessKey{
		Method: "aes-192-gcm",
	}

	// Act
	result, err := client.CreateAccessKey(ctx, createAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestCreateAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			body:       []byte("Not Found"),
		},
		{
			name:       "teapot",
			statusCode: http.StatusTeapot,
			body:       []byte("I'm a teapot"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			createAccessKey := &types.CreateAccessKey{
				Method: "aes-192-gcm",
			}

			// Act
			result, err := client.CreateAccessKey(ctx, createAccessKey)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestCreateAccessKey_Headers(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "key-headers",
		Name:      "Test",
		Password:  "pass",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	createAccessKey := &types.CreateAccessKey{
		Method: "aes-192-gcm",
	}

	// Act
	_, err := client.CreateAccessKey(ctx, createAccessKey)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

// === GetAccessKeys Tests ===

func TestGetAccessKeys_Success(t *testing.T) {
	tests := []struct {
		name         string
		expectedKeys []*types.AccessKey
	}{
		{
			name: "multiple keys",
			expectedKeys: []*types.AccessKey{
				{
					ID:        "key-1",
					Name:      "First Key",
					Password:  "pass1",
					Port:      8080,
					Method:    "aes-192-gcm",
					AccessURL: "ss://first@example.com:8080",
				},
				{
					ID:        "key-2",
					Name:      "Second Key",
					Password:  "pass2",
					Port:      8081,
					Method:    "aes-256-gcm",
					AccessURL: "ss://second@example.com:8081",
				},
			},
		},
		{
			name: "single key",
			expectedKeys: []*types.AccessKey{
				{
					ID:        "key-single",
					Name:      "Only Key",
					Password:  "onlypass",
					Port:      9000,
					Method:    "chacha20-ietf-poly1305",
					AccessURL: "ss://only@example.com:9000",
				},
			},
		},
		{
			name:         "empty list",
			expectedKeys: []*types.AccessKey{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			responseBody := struct {
				AccessKeys []*types.AccessKey `json:"accessKeys"`
			}{
				AccessKeys: tt.expectedKeys,
			}
			respBody, _ := json.Marshal(responseBody)
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusOK,
				Body:       respBody,
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.GetAccessKeys(ctx)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Len(t, result, len(tt.expectedKeys))
			for i, key := range result {
				assert.Equal(t, tt.expectedKeys[i].ID, key.ID)
				assert.Equal(t, tt.expectedKeys[i].Name, key.Name)
				assert.Equal(t, tt.expectedKeys[i].Password, key.Password)
				assert.Equal(t, tt.expectedKeys[i].Port, key.Port)
				assert.Equal(t, tt.expectedKeys[i].Method, key.Method)
				assert.Equal(t, tt.expectedKeys[i].AccessURL, key.AccessURL)
			}
			assert.Equal(t, http.MethodGet, req.Method)
		})
	}
}

func TestGetAccessKeys_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKeys(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestGetAccessKeys_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("invalid json"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKeys(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestGetAccessKeys_EmptyBody(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte{},
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKeys(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestGetAccessKeys_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			body:       []byte("Not Found"),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       []byte("Unauthorized"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.GetAccessKeys(ctx)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestGetAccessKeys_Headers(t *testing.T) {
	// Arrange
	responseBody := struct {
		AccessKeys []*types.AccessKey `json:"accessKeys"`
	}{
		AccessKeys: []*types.AccessKey{},
	}
	respBody, _ := json.Marshal(responseBody)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	_, err := client.GetAccessKeys(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

func TestGetAccessKeys_NilResponseKeys(t *testing.T) {
	// Arrange
	responseBody := struct {
		AccessKeys []*types.AccessKey `json:"accessKeys"`
	}{
		AccessKeys: nil,
	}
	respBody, _ := json.Marshal(responseBody)
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKeys(ctx)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

// === GetAccessKey Tests ===

func TestGetAccessKey_Success(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
		expectedKey types.AccessKey
	}{
		{
			name:        "valid key",
			accessKeyID: "key-123",
			expectedKey: types.AccessKey{
				ID:        "key-123",
				Name:      "Test Key",
				Password:  "securepassword",
				Port:      8388,
				Method:    "aes-192-gcm",
				AccessURL: "ss://test@example.com:8388",
			},
		},
		{
			name:        "key with empty name",
			accessKeyID: "key-456",
			expectedKey: types.AccessKey{
				ID:        "key-456",
				Name:      "",
				Password:  "pass456",
				Port:      9000,
				Method:    "chacha20-ietf-poly1305",
				AccessURL: "ss://user@example.com:9000",
			},
		},
		{
			name:        "numeric id",
			accessKeyID: "12345",
			expectedKey: types.AccessKey{
				ID:        "12345",
				Name:      "Numeric Key",
				Password:  "numpass",
				Port:      8080,
				Method:    "aes-256-gcm",
				AccessURL: "ss://num@example.com:8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			respBody, _ := json.Marshal(tt.expectedKey)
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusOK,
				Body:       respBody,
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.GetAccessKey(ctx, tt.accessKeyID)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedKey.ID, result.ID)
			assert.Equal(t, tt.expectedKey.Name, result.Name)
			assert.Equal(t, tt.expectedKey.Password, result.Password)
			assert.Equal(t, tt.expectedKey.Port, result.Port)
			assert.Equal(t, tt.expectedKey.Method, result.Method)
			assert.Equal(t, tt.expectedKey.AccessURL, result.AccessURL)
			assert.Equal(t, http.MethodGet, req.Method)
		})
	}
}

func TestGetAccessKey_NotFound(t *testing.T) {
	// Arrange
	accessKeyID := "non-existent-key"
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKey(ctx, accessKeyID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
}

func TestGetAccessKey_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKey(ctx, "key-123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestGetAccessKey_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("invalid json"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKey(ctx, "key-123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestGetAccessKey_EmptyBody(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte{},
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKey(ctx, "key-123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestGetAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       []byte("Unauthorized"),
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       []byte("Forbidden"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.GetAccessKey(ctx, "key-123")

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestGetAccessKey_Headers(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "key-headers",
		Name:      "Test",
		Password:  "pass",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	_, err := client.GetAccessKey(ctx, "key-headers")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

func TestGetAccessKey_EmptyAccessKeyID(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "",
		Name:      "Empty ID Key",
		Password:  "pass",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetAccessKey(ctx, "")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.MethodGet, capturedReq.Method)
}

func TestGetAccessKey_SpecialCharactersInID(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
	}{
		{
			name:        "id with hyphen",
			accessKeyID: "key-with-hyphen",
		},
		{
			name:        "id with underscore",
			accessKeyID: "key_with_underscore",
		},
		{
			name:        "uuid format",
			accessKeyID: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			expectedKey := types.AccessKey{
				ID:        tt.accessKeyID,
				Name:      "Special Key",
				Password:  "pass",
				Port:      8080,
				Method:    "aes-192-gcm",
				AccessURL: "ss://test@example.com:8080",
			}

			respBody, _ := json.Marshal(expectedKey)
			var capturedReq *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusOK,
				Body:       respBody,
			}, nil, &capturedReq)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.GetAccessKey(ctx, tt.accessKeyID)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.accessKeyID, result.ID)
		})
	}
}

// === UpdateAccessKey Tests ===

func TestUpdateAccessKey_Success(t *testing.T) {
	tests := []struct {
		name            string
		accessKeyID     string
		updateAccessKey *types.AccessKey
		expectedKey     types.AccessKey
	}{
		{
			name:        "update all fields",
			accessKeyID: "key-123",
			updateAccessKey: &types.AccessKey{
				ID:        "key-123",
				Name:      "Updated Key",
				Password:  "newpassword",
				Port:      9000,
				Method:    "aes-256-gcm",
				AccessURL: "ss://updated@example.com:9000",
			},
			expectedKey: types.AccessKey{
				ID:        "key-123",
				Name:      "Updated Key",
				Password:  "newpassword",
				Port:      9000,
				Method:    "aes-256-gcm",
				AccessURL: "ss://updated@example.com:9000",
			},
		},
		{
			name:        "update name only",
			accessKeyID: "key-456",
			updateAccessKey: &types.AccessKey{
				ID:   "key-456",
				Name: "New Name Only",
			},
			expectedKey: types.AccessKey{
				ID:        "key-456",
				Name:      "New Name Only",
				Password:  "existing-pass",
				Port:      8080,
				Method:    "aes-192-gcm",
				AccessURL: "ss://existing@example.com:8080",
			},
		},
		{
			name:            "nil update access key",
			accessKeyID:     "key-789",
			updateAccessKey: nil,
			expectedKey: types.AccessKey{
				ID:        "key-789",
				Name:      "Default",
				Password:  "default-pass",
				Port:      8080,
				Method:    "chacha20-ietf-poly1305",
				AccessURL: "ss://default@example.com:8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			respBody, _ := json.Marshal(tt.expectedKey)
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusCreated,
				Body:       respBody,
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			result, err := client.UpdateAccessKey(ctx, tt.accessKeyID, tt.updateAccessKey)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedKey.ID, result.ID)
			assert.Equal(t, tt.expectedKey.Name, result.Name)
			assert.Equal(t, tt.expectedKey.Password, result.Password)
			assert.Equal(t, tt.expectedKey.Port, result.Port)
			assert.Equal(t, tt.expectedKey.Method, result.Method)
			assert.Equal(t, tt.expectedKey.AccessURL, result.AccessURL)
			assert.Equal(t, http.MethodPut, req.Method)
		})
	}
}

func TestUpdateAccessKey_RequestBody(t *testing.T) {
	// Arrange
	accessKeyID := "key-body-test"
	updateAccessKey := &types.AccessKey{
		ID:        "key-body-test",
		Name:      "Body Test Key",
		Password:  "testpassword",
		Port:      9500,
		Method:    "aes-192-gcm",
		AccessURL: "ss://body@example.com:9500",
	}

	respBody, _ := json.Marshal(updateAccessKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	_, err := client.UpdateAccessKey(ctx, accessKeyID, updateAccessKey)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)

	var sentBody types.AccessKey
	err = json.Unmarshal(capturedReq.Body, &sentBody)
	require.NoError(t, err)

	assert.Equal(t, updateAccessKey.ID, sentBody.ID)
	assert.Equal(t, updateAccessKey.Name, sentBody.Name)
	assert.Equal(t, updateAccessKey.Password, sentBody.Password)
	assert.Equal(t, updateAccessKey.Port, sentBody.Port)
	assert.Equal(t, updateAccessKey.Method, sentBody.Method)
}

func TestUpdateAccessKey_NilRequestBody(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "key-nil-body",
		Name:      "Nil Body Key",
		Password:  "generated",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://nil@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.UpdateAccessKey(ctx, "key-nil-body", nil)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, capturedReq)
	assert.Empty(t, capturedReq.Body)
}

func TestUpdateAccessKey_NotFound(t *testing.T) {
	// Arrange
	accessKeyID := "non-existent-key"
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		ID:   accessKeyID,
		Name: "Should Not Update",
	}

	// Act
	result, err := client.UpdateAccessKey(ctx, accessKeyID, updateAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
}

func TestUpdateAccessKey_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		ID:   "key-123",
		Name: "Update Error",
	}

	// Act
	result, err := client.UpdateAccessKey(ctx, "key-123", updateAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestUpdateAccessKey_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       []byte("invalid json"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		ID:   "key-123",
		Name: "Invalid Response",
	}

	// Act
	result, err := client.UpdateAccessKey(ctx, "key-123", updateAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestUpdateAccessKey_EmptyBody(t *testing.T) {
	// Arrange
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       []byte{},
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		ID:   "key-123",
		Name: "Empty Response",
	}

	// Act
	result, err := client.UpdateAccessKey(ctx, "key-123", updateAccessKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalEmptyBodyError)
}

func TestUpdateAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       []byte("Unauthorized"),
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       []byte("Forbidden"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			updateAccessKey := &types.AccessKey{
				ID:   "key-123",
				Name: "Unexpected Status",
			}

			// Act
			result, err := client.UpdateAccessKey(ctx, "key-123", updateAccessKey)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestUpdateAccessKey_Headers(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "key-headers",
		Name:      "Headers Test",
		Password:  "pass",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		ID:   "key-headers",
		Name: "Headers Test",
	}

	// Act
	_, err := client.UpdateAccessKey(ctx, "key-headers", updateAccessKey)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

func TestUpdateAccessKey_EmptyAccessKeyID(t *testing.T) {
	// Arrange
	expectedKey := types.AccessKey{
		ID:        "",
		Name:      "Empty ID Key",
		Password:  "pass",
		Port:      8080,
		Method:    "aes-192-gcm",
		AccessURL: "ss://test@example.com:8080",
	}

	respBody, _ := json.Marshal(expectedKey)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusCreated,
		Body:       respBody,
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	updateAccessKey := &types.AccessKey{
		Name: "Empty ID Update",
	}

	// Act
	result, err := client.UpdateAccessKey(ctx, "", updateAccessKey)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, http.MethodPut, capturedReq.Method)
}

// === DeleteAccessKey Tests ===

func TestDeleteAccessKey_Success(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
	}{
		{
			name:        "delete existing key",
			accessKeyID: "key-123",
		},
		{
			name:        "delete key with numeric id",
			accessKeyID: "12345",
		},
		{
			name:        "delete key with uuid",
			accessKeyID: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusNoContent,
				Body:       []byte{},
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.DeleteAccessKey(ctx, tt.accessKeyID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, http.MethodDelete, req.Method)
		})
	}
}

func TestDeleteAccessKey_NotFound(t *testing.T) {
	// Arrange
	accessKeyID := "non-existent-key"
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteAccessKey(ctx, accessKeyID)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
}

func TestDeleteAccessKey_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteAccessKey(ctx, "key-123")

	// Assert
	assert.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestDeleteAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       []byte("Unauthorized"),
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       []byte("Forbidden"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.DeleteAccessKey(ctx, "key-123")

			// Assert
			assert.Error(t, err)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestDeleteAccessKey_Headers(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteAccessKey(ctx, "key-headers")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

func TestDeleteAccessKey_EmptyAccessKeyID(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteAccessKey(ctx, "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, capturedReq.Method)
}

func TestDeleteAccessKey_SpecialCharactersInID(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
	}{
		{
			name:        "id with hyphen",
			accessKeyID: "key-with-hyphen",
		},
		{
			name:        "id with underscore",
			accessKeyID: "key_with_underscore",
		},
		{
			name:        "uuid format",
			accessKeyID: "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var capturedReq *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusNoContent,
				Body:       []byte{},
			}, nil, &capturedReq)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.DeleteAccessKey(ctx, tt.accessKeyID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, http.MethodDelete, capturedReq.Method)
		})
	}
}

// === UpdateNameAccessKey Tests ===

func TestUpdateNameAccessKey_Success(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
		newName     string
	}{
		{
			name:        "update with valid name",
			accessKeyID: "key-123",
			newName:     "New Key Name",
		},
		{
			name:        "update with empty name",
			accessKeyID: "key-456",
			newName:     "",
		},
		{
			name:        "update with long name",
			accessKeyID: "key-789",
			newName:     "This is a very long name for the access key that should still work",
		},
		{
			name:        "update with special characters",
			accessKeyID: "key-special",
			newName:     "Key with émojis 🔑 and спец символы",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusNoContent,
				Body:       []byte{},
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.UpdateNameAccessKey(ctx, tt.accessKeyID, tt.newName)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, http.MethodPut, req.Method)
		})
	}
}

func TestUpdateNameAccessKey_RequestBody(t *testing.T) {
	// Arrange
	accessKeyID := "key-body-test"
	newName := "Updated Name"

	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateNameAccessKey(ctx, accessKeyID, newName)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)

	var sentBody struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(capturedReq.Body, &sentBody)
	require.NoError(t, err)
	assert.Equal(t, newName, sentBody.Name)
}

func TestUpdateNameAccessKey_NotFound(t *testing.T) {
	// Arrange
	accessKeyID := "non-existent-key"
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateNameAccessKey(ctx, accessKeyID, "New Name")

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
}

func TestUpdateNameAccessKey_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, networkError, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateNameAccessKey(ctx, "key-123", "New Name")

	// Assert
	assert.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestUpdateNameAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			body:       []byte("Bad Request"),
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			body:       []byte("Internal Server Error"),
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       []byte("Unauthorized"),
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       []byte("Forbidden"),
		},
		{
			name:       "service unavailable",
			statusCode: http.StatusServiceUnavailable,
			body:       []byte("Service Unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.UpdateNameAccessKey(ctx, "key-123", "New Name")

			// Assert
			assert.Error(t, err)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
		})
	}
}

func TestUpdateNameAccessKey_Headers(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateNameAccessKey(ctx, "key-headers", "Test Name")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
}

func TestUpdateNameAccessKey_EmptyAccessKeyID(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateNameAccessKey(ctx, "", "New Name")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, capturedReq.Method)
}

// === DeleteDataLimitAccessKey Tests ===

func TestDeleteDataLimitAccessKey_Success(t *testing.T) {
	tests := []struct {
		name        string
		accessKeyID string
	}{
		{
			name:        "delete data limit for valid key",
			accessKeyID: "key-123",
		},
		{
			name:        "delete data limit for key with numbers",
			accessKeyID: "key-456",
		},
		{
			name:        "delete data limit for key with special chars",
			accessKeyID: "key-special_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: http.StatusNoContent,
				Body:       []byte{},
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.DeleteDataLimitAccessKey(ctx, tt.accessKeyID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, http.MethodDelete, req.Method)
		})
	}
}

func TestDeleteDataLimitAccessKey_NotFound(t *testing.T) {
	// Arrange
	accessKeyID := "nonexistent-key"
	var req *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte(`{"error": "access key not found"}`),
	}, nil, &req)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteDataLimitAccessKey(ctx, accessKeyID)

	// Assert
	require.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.Equal(t, http.MethodDelete, req.Method)
	assert.Contains(t, req.URL, accessKeyID)
}

func TestDeleteDataLimitAccessKey_DoerError(t *testing.T) {
	// Arrange
	accessKeyID := "key-doer-error"
	expectedErr := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, expectedErr, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteDataLimitAccessKey(ctx, accessKeyID)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
}

func TestDeleteDataLimitAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody []byte
	}{
		{
			name:         "internal server error",
			statusCode:   http.StatusInternalServerError,
			responseBody: []byte(`{"error": "internal error"}`),
		},
		{
			name:         "bad request",
			statusCode:   http.StatusBadRequest,
			responseBody: []byte(`{"error": "bad request"}`),
		},
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			responseBody: []byte(`{"error": "unauthorized"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			accessKeyID := "key-unexpected"
			var req *contracts.Request
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.responseBody,
			}, nil, &req)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.DeleteDataLimitAccessKey(ctx, accessKeyID)

			// Assert
			require.Error(t, err)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
			assert.Equal(t, http.MethodDelete, req.Method)
			assert.Contains(t, req.URL, accessKeyID)
		})
	}
}

func TestDeleteDataLimitAccessKey_Headers(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteDataLimitAccessKey(ctx, "key-headers")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, capturedReq)
	assert.Equal(t, "application/json", capturedReq.Headers["Content-Type"])
	assert.Equal(t, "application/json", capturedReq.Headers["Accept"])
	assert.Equal(t, http.MethodDelete, capturedReq.Method)
	assert.Contains(t, capturedReq.URL, "key-headers")
}

func TestDeleteDataLimitAccessKey_EmptyAccessKeyID(t *testing.T) {
	// Arrange
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteDataLimitAccessKey(ctx, "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, capturedReq.Method)
}

func TestUpdateDataLimitAccessKey_Success(t *testing.T) {
	// Arrange
	accessKeyID := "key-123"
	bytes := uint64(1000000)
	var capturedReq *contracts.Request
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &capturedReq)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateDataLimitAccessKey(ctx, accessKeyID, bytes)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, capturedReq.Method)
	assert.Contains(t, capturedReq.URL, accessKeyID)

	// Verify request body
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	err = json.Unmarshal(capturedReq.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, bytes, reqBody.Limit.Bytes)
}

func TestUpdateDataLimitAccessKey_DoerError(t *testing.T) {
	// Arrange
	accessKeyID := "key-doer-error"
	bytes := uint64(500000)
	expectedErr := errors.New("network error")
	mockDoer := newMockDoerAccessKey(t, nil, expectedErr, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateDataLimitAccessKey(ctx, accessKeyID, bytes)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
}

func TestUpdateDataLimitAccessKey_InvalidDataLimit(t *testing.T) {
	// Arrange
	accessKeyID := "key-invalid-limit"
	bytes := uint64(0) // Assuming 0 is invalid
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte(`{"error": "invalid data limit"}`),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateDataLimitAccessKey(ctx, accessKeyID, bytes)

	// Assert
	require.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidDataLimitError)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
}

func TestUpdateDataLimitAccessKey_AccessKeyNotFound(t *testing.T) {
	// Arrange
	accessKeyID := "key-not-found"
	bytes := uint64(1000000)
	mockDoer := newMockDoerAccessKey(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte(`{"error": "access key not found"}`),
	}, nil, nil)

	client := createTestClientForAccessKeys(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateDataLimitAccessKey(ctx, accessKeyID, bytes)

	// Assert
	require.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, AccessKeyNotFoundError)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
}

func TestUpdateDataLimitAccessKey_UnexpectedStatusCode(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody []byte
	}{
		{
			name:         "internal server error",
			statusCode:   http.StatusInternalServerError,
			responseBody: []byte(`{"error": "internal error"}`),
		},
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			responseBody: []byte(`{"error": "unauthorized"}`),
		},
		{
			name:         "forbidden",
			statusCode:   http.StatusForbidden,
			responseBody: []byte(`{"error": "forbidden"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			accessKeyID := "key-unexpected"
			bytes := uint64(2000000)
			mockDoer := newMockDoerAccessKey(t, &contracts.Response{
				StatusCode: tt.statusCode,
				Body:       tt.responseBody,
			}, nil, nil)

			client := createTestClientForAccessKeys(mockDoer)
			ctx := context.Background()

			// Act
			err := client.UpdateDataLimitAccessKey(ctx, accessKeyID, bytes)

			// Assert
			require.Error(t, err)
			var clientErr *ClientError
			assert.ErrorAs(t, err, &clientErr)
			assert.ErrorIs(t, err, ClientOutlineError)
			assert.ErrorIs(t, err, UnexpectedStatusCodeError)
			assert.Equal(t, tt.statusCode, clientErr.statusCode)
		})
	}
}

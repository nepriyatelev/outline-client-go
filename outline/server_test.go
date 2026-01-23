package outline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/internal/logger"
	"github.com/nepriyatelev/outline-client-go/outline/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// createTestClient creates a Client with mock doer for testing
func createTestClient(doer contracts.Doer) *Client {
	baseURL, _ := url.Parse("http://localhost:8081/api/")
	return &Client{
		secret:                             "test-secret",
		getServerInfoPath:                  urlJoin(baseURL, "server"),
		putServerHostnamePath:              urlJoin(baseURL, "server/hostname-for-access-keys"),
		putServerPortPath:                  urlJoin(baseURL, "server/port-for-new-access-keys"),
		putServerNamePath:                  urlJoin(baseURL, "server/name"),
		getMetricsEnabledPath:              urlJoin(baseURL, "metrics/enabled"),
		putMetricsEnabledPath:              urlJoin(baseURL, "metrics/enabled"),
		putServerAccessKeyDataLimitPath:    urlJoin(baseURL, "server/access-key-data-limit"),
		deleteServerAccessKeyDataLimitPath: urlJoin(baseURL, "server/access-key-data-limit"),
		doer:                               doer,
		logger:                             logger.NewNoopLogger(),
	}
}

func urlJoin(baseURL *url.URL, path string) *url.URL {
	u, _ := url.Parse(baseURL.String() + path)
	return u
}

// newMockDoer configures generated mock to return provided response/error and capture the request.
func newMockDoer(t *testing.T, resp *contracts.Response, err error, capture **contracts.Request) *MockDoer {
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

// === GetServerInfo Tests ===

func TestGetServerInfo_Success(t *testing.T) {
	// Arrange
	serverInfo := types.ServerInfoResponse{
		Name:                  "Test Server",
		ServerID:              "server-123",
		MetricsEnabled:        true,
		CreatedTimestampMs:    1234567890.0,
		Version:               "1.0.0",
		PortForNewAccessKeys:  8000,
		HostnameForAccessKeys: "example.com",
	}

	respBody, _ := json.Marshal(serverInfo)
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetServerInfo(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test Server", result.Name)
	assert.Equal(t, "server-123", result.ServerID)
	assert.True(t, result.MetricsEnabled)
	assert.Equal(t, "1.0.0", result.Version)
	assert.Equal(t, 8000, result.PortForNewAccessKeys)
	assert.Equal(t, "example.com", result.HostnameForAccessKeys)
	assert.Equal(t, http.MethodGet, req.Method)
}

func TestGetServerInfo_DoerError(t *testing.T) {
	// Arrange
	networkError := errors.New("network error")
	mockDoer := newMockDoer(t, nil, networkError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetServerInfo(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkError)
}

func TestGetServerInfo_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("invalid json"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetServerInfo(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var ue *UnmarshalError
	assert.ErrorAs(t, err, &ue)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestGetServerInfo_NotFound(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetServerInfo(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

func TestGetServerInfo_UnexpectedStatus(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("Teapot"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetServerInfo(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === UpdateServerHostname Tests ===

func TestUpdateServerHostname_Success(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	hostname := "new-hostname.com"

	// Act
	err := client.UpdateServerHostname(ctx, hostname)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)

	// Verify request body
	var reqBody struct {
		Hostname string `json:"hostname"`
	}
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, hostname, reqBody.Hostname)
}

func TestUpdateServerHostname_DoerError(t *testing.T) {
	// Arrange
	connectionFailedError := errors.New("connection failed")
	mockDoer := newMockDoer(t, nil, connectionFailedError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateServerHostname(ctx, "hostname.com")

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, connectionFailedError)
}

func TestUpdateServerHostname_InvalidHostname(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid hostname"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	hostname := "invalid@hostname"

	// Act
	err := client.UpdateServerHostname(ctx, hostname)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidHostnameError)
}

func TestUpdateServerHostname_InternalServerError(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       []byte("Server error"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	hostname := "valid-hostname.com"

	// Act
	err := client.UpdateServerHostname(ctx, hostname)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusInternalServerError, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InternalHostNameError)
}

func TestUpdateServerHostname_UnexpectedError(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("teapot"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	hostname := "valid-hostname.com"

	// Act
	err := client.UpdateServerHostname(ctx, hostname)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === UpdatePortNewAccessKeys Tests ===

func TestUpdatePortNewAccessKeys_Success(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(9000)

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)

	// Verify request body
	var reqBody struct {
		Port uint16 `json:"port"`
	}
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, port, reqBody.Port)
}

func TestUpdatePortNewAccessKeys_InvalidPort(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid port"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(0)

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidPortError)
}

func TestUpdatePortNewAccessKeys_PortInUse(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusConflict,
		Body:       []byte("Port already in use"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(8080)

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusConflict, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, PortAlreadyInUseError)
}

func TestUpdatePortNewAccessKeys_DoerError(t *testing.T) {
	// Arrange
	networkTimeoutError := errors.New("network timeout")
	mockDoer := newMockDoer(t, nil, networkTimeoutError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, 8000)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, networkTimeoutError)
}

func TestUpdatePortNewAccessKeys_MaxPort(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(65535)

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)
	var body struct {
		Port uint16 `json:"port"`
	}
	require.NoError(t, json.Unmarshal(req.Body, &body))
	assert.Equal(t, port, body.Port)
}

func TestUpdatePortNewAccessKeys_TooLargePort(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid port"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(0) // simulate >65535 by relying on server validation returning 400

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidPortError)
}

func TestUpdatePortNewAccessKeys_UnexpectedError(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("teapot"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	port := uint16(8080)

	// Act
	err := client.UpdatePortNewAccessKeys(ctx, port)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === UpdateServerName Tests ===

func TestUpdateServerName_Success(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	serverName := "My Outline Server"

	// Act
	err := client.UpdateServerName(ctx, serverName)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)

	// Verify request body
	var reqBody struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, serverName, reqBody.Name)
}

func TestUpdateServerName_InvalidName(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid name"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	serverName := "" // Empty name might be invalid

	// Act
	err := client.UpdateServerName(ctx, serverName)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidServerNameError)
}

func TestUpdateServerName_DoerError(t *testing.T) {
	// Arrange
	requestFailedError := errors.New("request failed")
	mockDoer := newMockDoer(t, nil, requestFailedError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateServerName(ctx, "Test Server")

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, requestFailedError)
}

func TestUpdateServerName_UnexpectedError(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("teapot"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	serverName := "" // Empty name might be invalid

	// Act
	err := client.UpdateServerName(ctx, serverName)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === GetMetricsEnabled Tests ===

func TestGetMetricsEnabled_Success(t *testing.T) {
	// Arrange
	metricsEnabled := types.MetricsEnabled{Enabled: true}
	respBody, _ := json.Marshal(metricsEnabled)

	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetMetricsEnabled(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
	assert.Equal(t, http.MethodGet, req.Method)
}

func TestGetMetricsEnabled_DisabledMetrics(t *testing.T) {
	// Arrange
	metricsEnabled := types.MetricsEnabled{Enabled: false}
	respBody, _ := json.Marshal(metricsEnabled)

	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       respBody,
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetMetricsEnabled(ctx)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Enabled)
}

func TestGetMetricsEnabled_InvalidJSON(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusOK,
		Body:       []byte("not json"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetMetricsEnabled(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var eu *UnmarshalError
	assert.ErrorAs(t, err, &eu)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnmarshalFailedError)
}

func TestGetMetricsEnabled_NotFound(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetMetricsEnabled(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

func TestGetMetricsEnabled_DoerError(t *testing.T) {
	// Arrange
	connectionRefusedError := errors.New("connection refused")
	mockDoer := newMockDoer(t, nil, connectionRefusedError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	result, err := client.GetMetricsEnabled(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, connectionRefusedError)
}

func TestUpdateMetricsEnabled_EnableMetrics(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateMetricsEnabled(ctx, true)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)

	// Verify request body
	var reqBody types.MetricsEnabled
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.True(t, reqBody.Enabled)
}

func TestUpdateMetricsEnabled_DoerError(t *testing.T) {
	// Arrange
	connectionLostError := errors.New("connection lost")
	mockDoer := newMockDoer(t, nil, connectionLostError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateMetricsEnabled(ctx, true)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, connectionLostError)
}

func TestUpdateMetricsEnabled_DisableMetrics(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateMetricsEnabled(ctx, false)

	// Assert
	require.NoError(t, err)

	// Verify request body
	var reqBody types.MetricsEnabled
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.False(t, reqBody.Enabled)
}

func TestUpdateMetricsEnabled_BadRequest(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid request"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateMetricsEnabled(ctx, true)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidRequestError)
}

func TestUpdateMetricsEnabled_UnexpectedStatus(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("weird"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateMetricsEnabled(ctx, true)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === UpdateKeyLimitBytes Tests ===

func TestUpdateKeyLimitBytes_Success(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	limitBytes := uint64(1000000000) // 1GB

	// Act
	err := client.UpdateKeyLimitBytes(ctx, limitBytes)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, req.Method)

	// Verify request body
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, limitBytes, reqBody.Limit.Bytes)
}

func TestUpdateKeyLimitBytes_ZeroBytes(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	limitBytes := uint64(0)

	// Act
	err := client.UpdateKeyLimitBytes(ctx, limitBytes)

	// Assert
	require.NoError(t, err)

	// Verify request body
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	err = json.Unmarshal(req.Body, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), reqBody.Limit.Bytes)
}

func TestUpdateKeyLimitBytes_InvalidLimit(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusBadRequest,
		Body:       []byte("Invalid data limit"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	limitBytes := uint64(1000000000)

	// Act
	err := client.UpdateKeyLimitBytes(ctx, limitBytes)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusBadRequest, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, InvalidDataLimitError)
}

func TestUpdateKeyLimitBytes_MaxUint64(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()
	limitBytes := uint64(^uint64(0))

	// Act
	err := client.UpdateKeyLimitBytes(ctx, limitBytes)

	// Assert
	require.NoError(t, err)
	var body struct {
		Limit types.Limit `json:"limit"`
	}
	require.NoError(t, json.Unmarshal(req.Body, &body))
	assert.Equal(t, limitBytes, body.Limit.Bytes)
}

func TestUpdateKeyLimitBytes_DoerError(t *testing.T) {
	// Arrange
	requestFailedError := errors.New("request failed")
	mockDoer := newMockDoer(t, nil, requestFailedError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateKeyLimitBytes(ctx, 1000000000)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, requestFailedError)
}

func TestUpdateKeyLimitBytes_UnexpectedError(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusTeapot,
		Body:       []byte("teapot"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.UpdateKeyLimitBytes(ctx, 1000000000)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusTeapot, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

// === DeleteKeyLimitBytes Tests ===

func TestDeleteKeyLimitBytes_Success(t *testing.T) {
	// Arrange
	var req *contracts.Request
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNoContent,
		Body:       []byte{},
	}, nil, &req)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteKeyLimitBytes(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, req.Method)
}

func TestDeleteKeyLimitBytes_NotFound(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusNotFound,
		Body:       []byte("Not Found"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteKeyLimitBytes(ctx)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusNotFound, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

func TestDeleteKeyLimitBytes_DoerError(t *testing.T) {
	// Arrange
	connectionLostError := errors.New("connection lost")
	mockDoer := newMockDoer(t, nil, connectionLostError, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteKeyLimitBytes(ctx)

	// Assert
	require.Error(t, err)
	var doErr *DoError
	assert.ErrorAs(t, err, &doErr)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, DoOperationError)
	assert.ErrorIs(t, err, connectionLostError)
}

func TestDeleteKeyLimitBytes_UnexpectedStatus(t *testing.T) {
	// Arrange
	mockDoer := newMockDoer(t, &contracts.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       []byte("boom"),
	}, nil, nil)

	client := createTestClient(mockDoer)
	ctx := context.Background()

	// Act
	err := client.DeleteKeyLimitBytes(ctx)

	// Assert
	assert.Error(t, err)
	var clientErr *ClientError
	assert.ErrorAs(t, err, &clientErr)
	assert.Equal(t, http.StatusInternalServerError, clientErr.statusCode)
	assert.ErrorIs(t, err, ClientOutlineError)
	assert.ErrorIs(t, err, UnexpectedStatusCodeError)
}

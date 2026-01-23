package outline

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// GetServerInfo retrieves and returns information about the server.
// It makes an HTTP GET request to fetch server details including name, version,
// and other metadata.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//
// Returns:
//   - *types.ServerInfoResponse - The server information data containing name,
//     version, and other metadata.
//   - error - One of the following:
//   - nil - On success (HTTP 200 OK)
//   - *ClientError - For unexpected HTTP status codes
//   - *UnmarshalError - If JSON parsing fails
//   - error - Other errors from the HTTP request execution
func (c *Client) GetServerInfo(ctx context.Context) (*types.ServerInfoResponse, error) {
	// Create an HTTP GET request to fetch server information
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getServerInfoPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "GetServerInfo", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetServerInfo(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Parse and return the JSON response on success
		return unmarshalJSONWithError[types.ServerInfoResponse](resp.Body)
	default:
		// Return an error for unexpected status codes
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateServerHostname changes the hostname or IP address for access keys.
// The hostname or IP address must be valid. If a hostname is provided,
// DNS must be set up independently of this API.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//   - hostnameOrIP: string - A valid hostname or IP address for the server.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful update (HTTP 204 No Content)
//   - *ClientError with Code 400 - If the hostname or IP address is invalid
//     (HTTP 400 Bad Request)
//   - *ClientError with Code 500 - If an internal server error occurs,
//     possibly due to network validation issues (HTTP 500 Internal Server Error)
//   - *ClientError - For other unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) UpdateServerHostname(ctx context.Context, hostnameOrIP string) error {
	// Build the request body with the hostname
	var reqBody struct {
		Hostname string `json:"hostname"`
	}

	reqBody.Hostname = hostnameOrIP
	reqBodyBytes, _ := json.Marshal(&reqBody)

	// Create an HTTP PUT request to update the hostname
	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerHostnamePath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "UpdateServerHostname", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateServerHostname(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - hostname was updated
		return nil
	case http.StatusBadRequest:
		// Invalid hostname or IP address was provided
		return errInvalidHostname(http.StatusBadRequest, hostnameOrIP)
	case http.StatusInternalServerError:
		// Internal server error, possibly due to network validation issues
		return errInternalHostname(http.StatusInternalServerError, hostnameOrIP)
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdatePortNewAccessKeys changes the default port for newly created access keys.
// The port can be one that is already in use by existing access keys.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//   - port: uint16 - The port number (1-65535) to set for new access keys.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful update (HTTP 204 No Content)
//   - *ClientError with Code 400 - If the port is invalid or outside the range
//     1-65535 (HTTP 400 Bad Request)
//   - *ClientError with Code 409 - If the port is already in use by another
//     service (HTTP 409 Conflict)
//   - *ClientError - For other unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) UpdatePortNewAccessKeys(ctx context.Context, port uint16) error {
	// Build the request body with the port number
	var reqBody struct {
		Port uint16 `json:"port"`
	}

	reqBody.Port = port
	reqBodyBytes, _ := json.Marshal(&reqBody)

	// Create an HTTP PUT request to update the port
	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerPortPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "UpdatePortNewAccessKeys", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdatePortNewAccessKeys(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - port was updated
		return nil
	case http.StatusBadRequest:
		// Port is invalid or outside the valid range
		return errInvalidPort(http.StatusBadRequest, port)
	case http.StatusConflict:
		// Port is already in use by another service
		return errPortAlreadyInUse(http.StatusConflict, port)
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateServerName renames the server.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//   - name: string - The new name for the server.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful update (HTTP 204 No Content)
//   - *ClientError with Code 400 - If the provided server name is invalid
//     (HTTP 400 Bad Request)
//   - *ClientError - For other unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) UpdateServerName(ctx context.Context, name string) error {
	// Build the request body with the new server name
	var reqBody struct {
		Name string `json:"name"`
	}

	reqBody.Name = name
	reqBodyBytes, _ := json.Marshal(&reqBody)

	// Create an HTTP PUT request to update the server name
	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerNamePath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "UpdateServerName", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateServerName(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - server name was updated
		return nil
	case http.StatusBadRequest:
		// Invalid server name was provided
		return errInvalidServerName(http.StatusBadRequest, name)
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// GetMetricsEnabled retrieves the current metrics sharing status.
// It returns whether metrics is being shared with the server.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//
// Returns:
//   - *types.MetricsEnabled - The metrics sharing status (true if enabled,
//     false if disabled).
//   - error - One of the following:
//   - nil - On success (HTTP 200 OK)
//   - *ClientError - For unexpected HTTP status codes
//   - *UnmarshalError - If JSON parsing fails
//   - error - Other errors from the HTTP request execution
func (c *Client) GetMetricsEnabled(ctx context.Context) (*types.MetricsEnabled, error) {
	// Create an HTTP GET request to fetch metrics status
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "GetMetricsEnabled", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetMetricsEnabled(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Parse and return the metrics status on success
		return unmarshalJSONWithError[types.MetricsEnabled](resp.Body)
	default:
		// Return an error for unexpected status codes
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateMetricsEnabled enables or disables sharing of metrics.
// Metrics sharing can be toggled on or off independently on the server.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//   - enabled: bool - Boolean flag to enable (true) or disable (false) metrics
//     sharing.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful update (HTTP 204 No Content)
//   - *ClientError with Code 400 - If the request body or parameters are
//     invalid (HTTP 400 Bad Request)
//   - *ClientError - For other unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) UpdateMetricsEnabled(ctx context.Context, enabled bool) error {
	// Build the request body with the metrics enabled status
	var reqBody types.MetricsEnabled
	reqBody.Enabled = enabled

	reqBodyBytes, _ := json.Marshal(&reqBody)

	// Create an HTTP PUT request to update metrics sharing status
	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "UpdateMetricsEnabled", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateMetricsEnabled(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - metrics status was updated
		return nil
	case http.StatusBadRequest:
		// Invalid request body or parameters
		return errInvalidRequest(http.StatusBadRequest, string(resp.Body))
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateKeyLimitBytes sets a server-wide data limit for newly created access keys.
// This limit applies to all new access keys that will be created after this call.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//   - bytes: uint64 - The data limit in bytes to apply to new access keys.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful update (HTTP 204 No Content)
//   - *ClientError with Code 400 - If the provided data limit value is invalid
//     (HTTP 400 Bad Request)
//   - *ClientError - For other unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) UpdateKeyLimitBytes(ctx context.Context, bytes uint64) error {
	// Build the request body with the data limit
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	reqBody.Limit.Bytes = bytes

	reqBodyBytes, _ := json.Marshal(reqBody)

	// Create an HTTP PUT request to update the data limit
	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerAccessKeyDataLimitPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "UpdateKeyLimitBytes", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateKeyLimitBytes(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - data limit was updated
		return nil
	case http.StatusBadRequest:
		// Invalid data limit value was provided
		return errInvalidDataLimit(http.StatusBadRequest, bytes)
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// DeleteKeyLimitBytes removes the server-wide data limit for access keys.
// After this call, newly created access keys will not have a data limit.
//
// Arguments:
//   - ctx: context.Context - The request context for managing request lifecycle,
//     timeouts, and cancellation.
//
// Returns:
//   - error - One of the following:
//   - nil - On successful deletion (HTTP 204 No Content)
//   - *ClientError - For unexpected HTTP status codes
//   - error - Other errors from the HTTP request execution
func (c *Client) DeleteKeyLimitBytes(ctx context.Context) error {
	// Create an HTTP DELETE request to remove the data limit
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     c.deleteServerAccessKeyDataLimitPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	// Log the request for debugging purposes
	c.logRequest(ctx, "DeleteKeyLimitBytes", req)

	// Execute the HTTP request
	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoDeleteKeyLimitBytes(err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success - data limit was deleted
		return nil
	default:
		// Unexpected status code
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

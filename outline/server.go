package outline

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// GetServerInfo retrieves information about the Outline server,
// including name, version, and other metadata.
//
// It returns [*ClientError] for unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) GetServerInfo(ctx context.Context) (*types.ServerInfoResponse, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getServerInfoPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "GetServerInfo", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetServerInfo(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONWithError[types.ServerInfoResponse](resp.Body)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateServerHostname changes the hostname or IP address for access keys.
// The provided value must be a valid hostname or IP address.
// If a hostname is provided, DNS must be configured independently.
//
// It returns [*ClientError] with code 400 if the hostname is invalid,
// [*ClientError] with code 500 for internal server errors (e.g., network validation issues),
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateServerHostname(ctx context.Context, hostnameOrIP string) error {
	var reqBody struct {
		Hostname string `json:"hostname"`
	}

	reqBody.Hostname = hostnameOrIP
	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerHostnamePath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateServerHostname", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateServerHostname(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidHostname(http.StatusBadRequest, hostnameOrIP)
	case http.StatusInternalServerError:
		return errInternalHostname(http.StatusInternalServerError, hostnameOrIP)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdatePortNewAccessKeys changes the default port for newly created access keys.
// The specified port can already be in use by existing access keys.
//
// It returns [*ClientError] with code 400 if the port is invalid,
// [*ClientError] with code 409 if the port is already in use by another service,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdatePortNewAccessKeys(ctx context.Context, port uint16) error {
	var reqBody struct {
		Port uint16 `json:"port"`
	}

	reqBody.Port = port
	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerPortPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdatePortNewAccessKeys", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdatePortNewAccessKeys(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidPort(http.StatusBadRequest, port)
	case http.StatusConflict:
		return errPortAlreadyInUse(http.StatusConflict, port)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateServerName renames the server to the specified name.
//
// It returns [*ClientError] with code 400 if the name is invalid,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateServerName(ctx context.Context, name string) error {
	var reqBody struct {
		Name string `json:"name"`
	}

	reqBody.Name = name
	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerNamePath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateServerName", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateServerName(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidServerName(http.StatusBadRequest, name)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// GetMetricsEnabled retrieves the current metrics sharing status.
//
// It returns [*ClientError] for unexpected HTTP status codes,
// [*UnmarshalError] if JSON parsing fails,
// or [*DoError] if the HTTP request fails.
func (c *Client) GetMetricsEnabled(ctx context.Context) (*types.MetricsEnabled, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "GetMetricsEnabled", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return nil, errDoGetMetricsEnabled(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return unmarshalJSONWithError[types.MetricsEnabled](resp.Body)
	default:
		return nil, errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateMetricsEnabled enables or disables sharing of metrics.
//
// It returns [*ClientError] with code 400 if the request body is invalid,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateMetricsEnabled(ctx context.Context, enabled bool) error {
	var reqBody types.MetricsEnabled
	reqBody.Enabled = enabled

	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateMetricsEnabled", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateMetricsEnabled(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidRequest(http.StatusBadRequest, string(resp.Body))
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// UpdateKeyLimitBytes sets a server-wide data limit for newly created access keys.
// This limit applies to all new access keys that will be created after this call.
//
// It returns [*ClientError] with code 400 if the data limit value is invalid,
// or [*DoError] if the HTTP request fails.
func (c *Client) UpdateKeyLimitBytes(ctx context.Context, bytes uint64) error {
	var reqBody struct {
		Limit types.Limit `json:"limit"`
	}
	reqBody.Limit.Bytes = bytes

	reqBodyBytes, _ := json.Marshal(reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerAccessKeyDataLimitPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "UpdateKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoUpdateKeyLimitBytes(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidDataLimit(http.StatusBadRequest, bytes)
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

// DeleteKeyLimitBytes removes the server-wide data limit for access keys.
// After this call, newly created access keys will not have a data limit.
//
// It returns [*ClientError] for unexpected HTTP status codes,
// or [*DoError] if the HTTP request fails.
func (c *Client) DeleteKeyLimitBytes(ctx context.Context) error {
	req := &contracts.Request{
		Method:  http.MethodDelete,
		URL:     c.deleteServerAccessKeyDataLimitPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "DeleteKeyLimitBytes", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return errDoDeleteKeyLimitBytes(err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return errUnexpectedStatusCode(resp.StatusCode, resp.Body)
	}
}

package outline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

// === Get Server Information ===

// GetServerInfo Returns information about the server.
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
		return nil, err
	}

	var serverInfo *types.ServerInfoResponse
	if err = json.Unmarshal(resp.Body, serverInfo); err != nil {
		return nil, err
	}

	return serverInfo, nil
}

// === Server Configuration ===

// UpdateServerHostname Changes the hostname for access keys. Must be a valid hostname or IP address.
// If it's a hostname, DNS must be set up independently of this API.
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
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return &ClientError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf("An invalid hostname or IP address was provided: %s.",
				hostnameOrIP),
		}
	case http.StatusInternalServerError:
		return &ClientError{
			Code: http.StatusInternalServerError,
			Message: fmt.Sprintf("An internal error occurred for host or IP: %s. "+
				"This could be thrown if there were network errors "+
				"while validating the hostname.",
				hostnameOrIP),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

// UpdatePortNewAccessKeys Changes the default port for newly created access keys.
// This can be a port already used for access keys.
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
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return &ClientError{
			Code: http.StatusBadRequest,
			Message: fmt.Sprintf(
				"The requested port wasn't an integer from 1 through 65535, "+
					"or the request had no port parameter. Provided: %d.", port),
		}
	case http.StatusConflict:
		return &ClientError{
			Code: http.StatusConflict,
			Message: fmt.Sprintf(
				"The requested port was already in use by another service: %d.",
				port),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

// UpdateServerName Renames the server.
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

	c.logRequest(ctx, "SetServerName", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return &ClientError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("An invalid server name was provided: %s.", name),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

// GetMetricsEnabled Returns whether metrics is being shared.
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
		return nil, err
	}

	metricsEnabled, err := unmarshalJSONWithError[types.MetricsEnabled](resp.Body)
	if err != nil {
		return nil, err
	}

	return metricsEnabled, nil
}

// Enables or disables sharing of metrics.
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
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return &ClientError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid request: %s.", string(reqBodyBytes)),
		}
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

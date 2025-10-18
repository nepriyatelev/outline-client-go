package outline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nepriyatelev/outline-client-go/internal/contracts"
	"github.com/nepriyatelev/outline-client-go/outline/types"
)

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

// PutServerHostname Changes the hostname for access keys. Must be a valid hostname or IP address.
// If it's a hostname, DNS must be set up independently of this API.
func (c *Client) PutServerHostname(ctx context.Context, hostname string) error {
	var reqBody struct {
		Hostname string `json:"hostname"`
	}

	reqBody.Hostname = hostname
	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putServerHostnamePath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "SetHostname", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	var (
		errInvalidHostnameOrIP = func(hostOrIP string) *ClientError {
			return &ClientError{
				Code: 400,
				Message: fmt.Sprintf("An invalid hostname or IP address was provided: %s.",
					hostOrIP),
			}
		}
		errInternal = func(hostOrIP string) *ClientError {
			return &ClientError{
				Code: 500,
				Message: fmt.Sprintf("An internal error occurred for host or IP: %s. "+
					"This could be thrown if there were network errors "+
					"while validating the hostname.",
					hostOrIP),
			}
		}
	)

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return errInvalidHostnameOrIP(hostname)
	case http.StatusInternalServerError:
		return errInternal(hostname)
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

// PutServerName Renames the server.
func (c *Client) PutServerName(ctx context.Context, name string) error {
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

	errInvalidServerName := func(name string) *ClientError {
		return &ClientError{
			Code:    400,
			Message: fmt.Sprintf("An invalid server name was provided: %s.", name),
		}
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidServerName(name)
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}

// GetMetricsEnabled Returns whether metrics is being shared.
func (c *Client) GetMetricsEnabled(ctx context.Context) (bool, error) {
	req := &contracts.Request{
		Method:  http.MethodGet,
		URL:     c.getMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    nil,
	}

	c.logRequest(ctx, "GetMetricsEnabled", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return false, err
	}

	var metricsResp struct {
		Enabled bool `json:"enabled"`
	}
	if err = json.Unmarshal(resp.Body, &metricsResp); err != nil {
		return false, err
	}

	return metricsResp.Enabled, nil
}

// Enables or disables sharing of metrics.
func (c *Client) PutMetricsEnabled(ctx context.Context, enabled bool) error {
	var reqBody struct {
		Enabled bool `json:"enabled"`
	}
	reqBody.Enabled = enabled

	reqBodyBytes, _ := json.Marshal(&reqBody)

	req := &contracts.Request{
		Method:  http.MethodPut,
		URL:     c.putMetricsEnabledPath.String(),
		Headers: DefaultHeaders(),
		Body:    reqBodyBytes,
	}

	c.logRequest(ctx, "PutMetricsEnabled", req)

	resp, err := c.doer.Do(ctx, req)
	if err != nil {
		return err
	}

	errInvalidRequest := func() *ClientError {
		return &ClientError{
			Code:    400,
			Message: fmt.Sprintf("Invalid request: %s.", string(reqBodyBytes)),
		}
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errInvalidRequest()
	default:
		return errUnexpected(resp.StatusCode, resp.Body)
	}
}
